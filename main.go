package main

import (
	"flag"
	"log"
	"net/http"
)

type makeGray struct {
	domain string
}

func (makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("%s %s %s", request.RemoteAddr, request.Method, request.URL)
}

func main() {
	const addr = "127.0.0.1:8080"

	domain := flag.String("domain", "https://maps.wikimedia.org", "The origin server against which incoming requests will be proxied")
	flag.Parse()

	log.Printf("Starting server on %s, proxying %s", addr, *domain)
	http.ListenAndServe(addr, &makeGray{*domain})
}
