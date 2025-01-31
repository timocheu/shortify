package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/timocheu/shortify/utils"
)

var dbClient = utils.NewRedisClient()
var ctx = context.Background()

func main() {
	if dbClient == nil {
		fmt.Println("Failed to connect to redis")
		return
	}

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/shorten", shortenPage)
	http.HandleFunc("/r/{code}", redirect)

	// Listener for incoming request
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))

	tmpl.Execute(w, nil)
}

func shortenPage(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	fmt.Println("Payload: ", url)

	shortURL := utils.GetShortCode()
	fullShortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortURL)

	// Print generated url
	fmt.Println("Generated URL: ", shortURL)

	utils.SetKey(&ctx, dbClient, shortURL, url, 0)

	// Write the response into the responese write
	fmt.Fprintf(w, `<p class="mt-4 text-blue-400">Shortened URL: <a href="/r/%s" class="underline">%s</a></p>`, shortURL, fullShortURL)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("code")

	// Check if the key exist or empty
	if key == "" {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	longURL, err := utils.GetLongURL(&ctx, dbClient, key)
	// Check if URL exist in redis server
	if err != nil {
		http.Error(w, "Unable to find the Shortened URL", http.StatusBadRequest)
		return
	}

	// Do the redirect after all error check
	http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
}
