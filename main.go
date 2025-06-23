package main

import (
	"fmt"
	"net/http"
	"url_shortner/utils"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "Welcome to URL Shortner!")
	})
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "I will take in a long url and shorten it!")
		url := req.FormValue("url")
		fmt.Println("Payload: ", url)
		gen := &utils.TimestampGenerator{}
		shortUrl := gen.Generate(url)
		fullShortUrl := fmt.Sprintf("http://localhost:8080/r/%s", shortUrl)
		fmt.Printf("Generated short url: %s", fullShortUrl)

	})
	fmt.Println("Starting a server on localhost:8080!")
	http.ListenAndServe(":8080", nil)
}
