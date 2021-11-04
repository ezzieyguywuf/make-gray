package main

import (
	"log"
	"net/http"
)

type makeGray struct{}

func (makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("%s %s %s", request.RemoteAddr, request.Method, request.URL)
}

func main() {
	const addr = "127.0.0.1:8080"
	log.Printf("Starting server on %s", addr)

	http.ListenAndServe(addr, &makeGray{})
}
