package main

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
)

// Takes a URL and returns the fetched page as *html.Node.
func fetch(URL string) (page *html.Node) {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	page, err = html.Parse(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return
}