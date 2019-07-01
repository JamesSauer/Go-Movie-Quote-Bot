package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

var listParts = [8]string{
	// Note how these are not regular dashes ("-"), but longer ones ("–").
	"A–C",
	"D–F",
	"G–I",
	"J–L",
	"M–O",
	"P–S",
	"T–V",
	"W–Z",
}

// Returns the URL of a random WikiQuotes.org page about a movie as a string.
func getRandomMoviePage() (movieURL string) {
	rand.Seed(time.Now().UnixNano())
	startURL := fmt.Sprintf("https://en.wikiquote.org/wiki/List_of_films_(%s)", listParts[rand.Intn(8)])

	document := fetchPage(startURL)

	movieLinks := make([]string, 0)
	document.Find("i > a").Each(func(i int, element *goquery.Selection) {
		href, _ := element.Attr("href")
		findRedlink := regexp.MustCompile("redlink=1")
		if !findRedlink.MatchString(href) {
			movieLinks = append(movieLinks, href)
		}
	})

	randomLink := fmt.Sprintf("%s", movieLinks[rand.Intn(len(movieLinks))])
	movieURL = "https://en.wikiquote.org" + randomLink
	return
}

// Takes a URL and returns the fetched page as *html.Node.
func fetchPage(URL string) (document *goquery.Document) {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	document, err = goquery.NewDocumentFromReader(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return
}