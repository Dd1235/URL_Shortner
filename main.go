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
		if code == "" {
			http.Error(w, "Short URL code is required", http.StatusBadRequest)
			return
		}
		longURL, err := urlStore.Get(ctx, code)
		if err != nil {
			http.NotFound(w, r) // reply with 404 not found error
		}
		http.Redirect(w, r, longURL, http.StatusFound)
	})
	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		code := r.URL.Path[len("/delete/"):]
		if code == "" {
			http.Error(w, "Short URL code is required", http.StatusBadRequest)
			return
		}
		err := urlStore.Delete(ctx, code)
		if err != nil {
			http.Error(w, "Failed to delete short URL: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, "Deleted short URL with code: %s\n", code)
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := urlStore.All(ctx)
		if err != nil {
			http.Error(w, "Failed to fetch data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for short, long := range data {
			fmt.Fprintf(w, `
				<li class="py-1 flex items-start justify-between">
				<span class="truncate max-w-[6rem] mr-2">%s →</span>

				<div class="flex-1 min-w-0">
					<a href="%s" target="_blank"
					class="text-blue-600 underline break-words">%s</a>
				</div>

				<button type="button"                 
						hx-delete="/delete/%s"
						hx-target="closest li"
						hx-swap="outerHTML"
						class="ml-3 text-xs text-red-500 hover:text-red-700">
					Delete
				</button>
				</li>`, short, long, long, short)

		}

	})

	fmt.Println("Starting a server on localhost:8080!")
	http.ListenAndServe(":8080", nil)

}
