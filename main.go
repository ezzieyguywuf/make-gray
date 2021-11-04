package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
)

type makeGray struct {
	domain string
}

func (mk makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	target := url.URL{
		Host: mk.domain,
		Path: request.URL.Path,
	}

	fetch := http.Request{
		Method: "GET",
		URL:    &target,
	}

	log.Println("fetching ", target)
	client := http.Client{Timeout: 5}

	response, err := client.Do(&fetch)
	if err != nil {
		// do something
	}

	// defer response.Body.Close()
	writer.Write([]byte("<h1>Hello, response</h1>"))

	log.Println("respons: ", response)
}

func main() {
	const addr = "127.0.0.1:8080"

	domain := flag.String("domain", "https://maps.wikimedia.org", "The origin server against which incoming requests will be proxied")
	flag.Parse()

	log.Printf("Starting server on %s, proxying %s", addr, *domain)
	http.ListenAndServe(addr, &makeGray{*domain})
}
