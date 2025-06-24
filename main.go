package main

import (
	"context"
	"fmt"
	"net/http"
	"text/template"
	"url_shortner/store"
	"url_shortner/utils"

	"github.com/joho/godotenv"
)

var urlStore store.URLStore

func main() {
	errLoad := godotenv.Load()
	if errLoad != nil {
		fmt.Println("No environment variable found!")
	}
	ctx := context.Background()
	urlStore = store.NewRedisStore()
	defer urlStore.Close()

	http.HandleFunc("/", func(writer http.ResponseWriter, req *http.Request) {
		tmpl := template.Must(template.ParseFiles("templates/index.html"))
		tmpl.Execute(writer, nil)
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
		fmt.Fprintf(writer, `
			<div class="mt-4 text-center">
				<p class="text-lg text-green-700 dark:text-green-400">
				Shortened URL: 
				<a href="/r/%s" class="underline hover:text-blue-600 dark:hover:text-blue-400">
					%s
				</a>
				</p>
				<button 
				onclick="navigator.clipboard.writeText('%s')" 
				class="mt-2 text-sm text-gray-600 dark:text-gray-400 underline hover:text-blue-500"
				>
				Copy to Clipboard
				</button>
			</div>
			`, shortCode, shortUrl, shortUrl)

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
