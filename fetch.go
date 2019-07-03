package main

import (
	"golang.org/x/net/html"
	"log"
	"net/http"
)

// Takes a URL and returns the fetched page as *html.Node.
func fetch(URL string) (page *html.Node) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "MovieQuoteBot https://github.com/JamesSauer/Go-Movie-Quote-Bot")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	page, err = html.Parse(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return
}