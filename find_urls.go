package main

import (
	"fmt"
	"math/rand"
	"time"
)

// The list of movies on wikiquote is split up into these parts:
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

// Returns a slice of relative WikiQuote.org URLs to pages about movies.
func getURLListFromPart(i int) (movieLinks []string, err error) {
	// Make sure i is in range:
	if i < 0 {
		i = i * -1
	}
	i = i % len(listParts)

	startURL := fmt.Sprintf("https://en.wikiquote.org/wiki/List_of_films_(%s)", listParts[i])

	document, err := fetch(startURL)
	if err != nil {
		return
	}

	movieLinks = make([]string, 0)
	for _, node := range querySelectorAll(document, "i") {
		link, err := extractLink(node)
		if err == nil {
			movieLinks = append(movieLinks, link)
		}
	}

	return
}

// Gets all relative movie page URLs on wikiquote and returns them as slice of strings.
func getAllURLs() (movieLinks []string, err error) {
	movieLinks = make([]string, 0)
	for i := range listParts {
		urlList, err := getURLListFromPart(i)
		if err != nil {
			return nil, err
		}
		movieLinks = append(movieLinks, urlList...)
	}
	return
}

// Returns a random absolute URL to a wikiquote page about a movie.
func getRandomURL() (movieURL string, err error) {
	rand.Seed(time.Now().UnixNano())
	urlList, err := getURLListFromPart(rand.Intn(len(listParts)))
	if err != nil {
		return
	}
	randomLink := fmt.Sprintf("%s", urlList[rand.Intn(len(urlList))])
	movieURL = "https://en.wikiquote.org" + randomLink
	return
}
