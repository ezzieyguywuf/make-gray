package main

import (
	"bytes"
	"flag"
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

func (mk makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target := mk.host
	target.Path = request.URL.Path

	fetch := http.Request{
		Method: "GET",
		URL:    &target,
	}

	client := http.Client{Timeout: 5 * time.Second}

	response, err := client.Do(&fetch)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Body.Close()

	imageData, format, err := image.Decode(response.Body)

	if err != nil {
		log.Println(err)
		return
	}

	bounds := imageData.Bounds()

	log.Println("Image bounds = ", bounds, ", format = ", format)

	min, max := bounds.Min, bounds.Max
	grayImg := image.NewGray16(bounds)

	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			grayImg.Set(x, y, imageData.At(x, y))
		}
	}

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
	const addr = "127.0.0.1:8080"

	host := flag.String("host", "https://maps.wikimedia.org", "The origin server against which incoming requests will be proxied")
	flag.Parse()

	url, err := url.Parse(*host)
	if err != nil {
		log.Printf("Invalid target host, %s. Please provide a valid URL", *host)
	}

	log.Printf("Starting server on %s, proxying %s", addr, url)
	http.ListenAndServe(addr, &makeGray{*url})
}
