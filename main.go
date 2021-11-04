package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type makeGray struct {
	host url.URL
}

// The returned string will represent the format of the image
func fetchImage(target url.URL) (image.Image, string, error) {
	fetch := http.Request{
		Method: "GET",
		URL:    &target,
	}

	client := http.Client{Timeout: 5 * time.Second}

	response, err := client.Do(&fetch)
	if err != nil {
		log.Println(err)
		return image.NewGray16(image.Rect(0, 0, 0, 0)), "", fmt.Errorf("make-grey: could not fetch image from %s", target)
	}

	defer response.Body.Close()

	return image.Decode(response.Body)
}

func transformImage(colorImage image.Image) (greyImage image.Image) {
	bounds := colorImage.Bounds()
	min, max := bounds.Min, bounds.Max
	grayImage := image.NewGray16(bounds)

	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			grayImage.Set(x, y, colorImage.At(x, y))
		}
	}

	return grayImage
}

func (mk makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target := mk.host
	target.Path = request.URL.Path

	imageData, format, err := fetchImage(target)
	if err != nil {
		log.Println(err)
		return
	}

	grayImg := transformImage(imageData)
	buffer := new(bytes.Buffer)

	if format == "png" {
		writer.Header().Set("Content-Type", "image/png")
		if err := png.Encode(buffer, grayImg); err != nil {
			log.Println("Could not encode image back to png")
			return
		}
	} else if format == "jpeg" {
		writer.Header().Set("Content-Type", "image/jpeg")
		if err := jpeg.Encode(buffer, grayImg, nil); err != nil {
			log.Println("Could not encode image back to jpeg")
			return
		}
	}

	writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		log.Println("Unable to write image to HTTP response")
		return
	}
}

func main() {
	host := flag.String("host", "https://maps.wikimedia.org", "The origin server against which incoming requests will be proxied")
	server := flag.String("server", "127.0.0.1", "The server to which requests can be sent")
	port := flag.String("port", "8080", "The port on the server on which to listen")
	flag.Parse()

	addr := *server + ":" + *port

	url, err := url.Parse(*host)
	if err != nil {
		log.Printf("Invalid target host, %s. Please provide a valid URL", *host)
	}

	log.Printf("Starting server on %s, proxying %s", addr, url)
	http.ListenAndServe(addr, &makeGray{*url})
}
