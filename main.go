package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
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

	data, err := io.ReadAll(response.Body)

	if err != nil {
		log.Println(err)
		return
	}

	writer.Write([]byte("<h1>Hello, response</h1>"))

	log.Println("read ", len(data), "bytes of data")
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
