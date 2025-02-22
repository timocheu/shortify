package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/timocheu/shortify/utils"
)

var dbClient = utils.NewLocalRedisClient()
var ctx = context.Background()
var counter int64 = 1000000000000

// html templates
type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("template/*.html")),
	}
}

// Sample counter
type Counter struct {
	Count int
}

// Test count
var count = Counter{Count: 0}

func main() {
	if dbClient == nil {
		fmt.Println("Failed to connect to redis")
		return
	}

	// Intitalize echo web
	e := echo.New()
	e.Use(middleware.Logger())
	// Render the page
	e.Renderer = newTemplate()

	e.GET("/", index)
	http.HandleFunc("/shorten", shortenPage)
	http.HandleFunc("/r/{code}", redirect)

	e.Logger.Fatal(e.Start(":8080"))
}

func index(c echo.Context) error {
	count.Count++
	return c.Render(200, "index", count)
}

func shortenPage(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	fmt.Println("Payload: ", url)

	shortURL := utils.GetShortCode(counter)
	counter++
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

	fmt.Printf("Key: %s\n", key)
	longURL, err := utils.GetLongURL(&ctx, dbClient, key)
	// Check if URL exist in redis server
	if err != nil {
		http.Error(w, "Unable to find the Shortened URL", http.StatusBadRequest)
		return
	}

	// Do the redirect after all error check
	http.Redirect(w, r, longURL, http.StatusPermanentRedirect)
}
