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

	// Index
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	// Shorten
	e.POST("/shorten", shortenPage)
	e.POST("/r{code}", redirect)

	e.Logger.Fatal(e.Start(":8080"))
}

func shortenPage(c echo.Context) error {
	url := c.FormValue("url")
	fmt.Println("Payload: ", url)

	shortURL := utils.GetShortCode(counter)
	counter++
	fullShortURL := fmt.Sprintf("http://localhost:8080/r/%s", shortURL)
	// Print generated url
	fmt.Println("Generated URL: ", shortURL)

	utils.SetKey(&ctx, dbClient, shortURL, url, 0)

	return c.Render(200, "url", fullShortURL)
}

func redirect(c echo.Context) error {
	key := c.Request().PathValue("code")

	// Check if the key exist or empty
	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid URL")
	}

	fmt.Printf("Key: %s\n", key)
	longURL, err := utils.GetLongURL(&ctx, dbClient, key)
	// Check if URL exist in redis server
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to find the Shortened URL")
	}

	// Do the redirect after all error check
	return c.Redirect(http.StatusPermanentRedirect, longURL)
}
