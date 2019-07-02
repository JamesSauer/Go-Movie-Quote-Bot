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
func getURLListFromPart(i int) (movieLinks []string) {
	// Make sure i is in range:
	if i < 0 {
		i = i * -1
	}
	i = i % len(listParts)

	startURL := fmt.Sprintf("https://en.wikiquote.org/wiki/List_of_films_(%s)", listParts[i])

	document := fetch(startURL)

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
func getAllURLs() (movieLinks []string) {
	movieLinks = make([]string, 0)
	for i := range listParts {
		movieLinks = append(movieLinks, getURLListFromPart(i)...)
	}
	return
}

// Returns a random absolute URL to a wikiquote page about a movie.
func getRandomURL() (movieURL string) {
	rand.Seed(time.Now().UnixNano())
	URLList := getURLListFromPart(rand.Intn(len(listParts)))

	randomLink := fmt.Sprintf("%s", URLList[rand.Intn(len(URLList))])
	movieURL = "https://en.wikiquote.org" + randomLink
	return
}