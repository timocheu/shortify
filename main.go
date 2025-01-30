package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	utils.SetKey(&ctx, dbClient, shortURL, url, 0)
}
