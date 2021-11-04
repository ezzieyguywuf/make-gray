package main

import (
	"fmt"
	"net/http"
)

type makeGray struct{}

func (makeGray) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.RemoteAddr, " ", request.Method, " ", request.URL)
}

func main() {
	fmt.Println("Starting server on 127.0.0.1:8080")

	http.ListenAndServe("127.0.0.1:8080", &makeGray{})
}
