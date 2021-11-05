package main

import (
	"bytes"
	"errors"
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

// fetchImage will retrieve an image from the URL. The returned string will
// contain the format of the image
func fetchImage(target url.URL) (image.Image, string, error) {
	fetch := http.Request{
		Method: "GET",
		URL:    &target,
	}

	client := http.Client{Timeout: 5 * time.Second}

	response, err := client.Do(&fetch)
	if err != nil {
		log.Printf("make-grey: could not fetch image from %v", target)
		return image.NewGray16(image.Rect(0, 0, 0, 0)), "", err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Printf("make-grey: statusCode %d in fetchImage", response.StatusCode)
		return image.NewGray16(image.Rect(0, 0, 0, 0)), "", errors.New("error in HTTP client")
	}

	defer response.Body.Close()

	img, status, err := image.Decode(response.Body)

	if err != nil {
		log.Print("make-grey: unable to decode image")
	}

	return img, status, err
}

// transformImage will change the incoming image to greyscale and return a
// buffer that is suitable for transport via HTTP
func transformImage(colorImage image.Image, format string) (*bytes.Buffer, error) {
	bounds := colorImage.Bounds()
	min, max := bounds.Min, bounds.Max
	grayImage := image.NewGray16(bounds)

	for x := min.X; x < max.X; x++ {
		for y := min.Y; y < max.Y; y++ {
			grayImage.Set(x, y, colorImage.At(x, y))
		}
	}

	buffer := new(bytes.Buffer)

	if format == "png" {
		err := png.Encode(buffer, grayImage)
		return buffer, err
	} else if format == "jpeg" {
		err := jpeg.Encode(buffer, grayImage, nil)
		return buffer, err
	}

	return nil, errors.New("make-gray: I only know how to transform png and jpeg images")
}

func (mk makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target := mk.host
	target.Path = request.URL.Path

	imageData, format, err := fetchImage(target)
	if err != nil {
		http.Error(writer, "Unable to fetch image", http.StatusBadRequest)
		return
	}

	buffer, err := transformImage(imageData, format)

	if err != nil {
		http.Error(writer, "Unable to transform image", http.StatusUnprocessableEntity)
		return
	}

	writer.Header().Set("Content-Type", "image/"+format)
	writer.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := writer.Write(buffer.Bytes()); err != nil {
		log.Println("Unable to write image to HTTP response")
		http.Error(writer, "Unable to transform image", http.StatusInternalServerError)
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
		return
	}

	log.Printf("Starting server on %s, proxying %s", addr, url)
	if err := http.ListenAndServe(addr, &makeGray{*url}); err != nil {
		log.Printf("Unable to start server. Error: %v", err)
		return
	}
}
