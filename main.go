package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "Welcome to URL Shortner!")
	})
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(writer, "I will take in a long url and shorten it!")
	})
	fmt.Println("Starting a server on localhost:8080!")
	http.ListenAndServe(":8080", nil)
}
