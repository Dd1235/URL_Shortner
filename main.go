package main

import (
	"context"
	"fmt"
	"net/http"
	"url_shortner/store"
	"url_shortner/utils"
)

var urlStore store.URLStore

func main() {
	ctx := context.Background()
	urlStore = store.NewRedisStore()
	defer urlStore.Close()

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(writer, "Welcome to URL Shortner!")
	})
	http.HandleFunc("/shorten", func(writer http.ResponseWriter, req *http.Request) {

		longUrl := req.FormValue("url")
		gen := &utils.TimestampGenerator{}
		shortCode := gen.Generate(longUrl)

		if err := urlStore.Set(ctx, shortCode, longUrl); err != nil {
			http.Error(writer, "Failed to store URL", http.StatusInternalServerError)
			return
		}

		shortUrl := fmt.Sprintf("http://localhost:8080/r/%s", shortCode)
		fmt.Fprintln(writer, shortUrl)
		fmt.Printf("The short url creaetd for %s is %s", longUrl, shortUrl)
	})
	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Path[len("/r/"):]
		longURL, err := urlStore.Get(ctx, code)
		if err != nil {
			http.NotFound(w, r) // reply with 404 not found error
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	})
	fmt.Println("Starting a server on localhost:8080!")
	http.ListenAndServe(":8080", nil)
}
