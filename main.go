package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/timocheu/shortify/utils"
)

func main() {
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/shorten", shortenPage)

	// Listener for incoming request
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexPage(w http.ResponseWriter, r *http.Request) {

}

func shortenPage(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	fmt.Println("Payload: ", url)

	shortURL := utils.GetShortCode()
	// fullShortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortURL)

	// Print generated url
	fmt.Println("Generated URL: ", shortURL)
}
