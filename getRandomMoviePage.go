package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"math/rand"
	"net/http"
	"regexp"
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

	document := fetchPage(startURL)

	movieLinks = make([]string, 0)
	for _, node := range querySelectorAll(document, "i") {
		link, err := extractLink(node)
			if err == nil {
				movieLinks = append(movieLinks, link)
			}
	}

	return
}

func getAllURLs() (movieLinks []string) {
	movieLinks = make([]string, 0)
	for i := range listParts {
		movieLinks = append(movieLinks, getURLListFromPart(i)...)
	}
	return
}

// Returns a random page from a slice of relative WikiQuote.org URLs as *html.Node.
func getRandomMoviePage() (moviePage *html.Node) {
	rand.Seed(time.Now().UnixNano())
	URLList := getURLListFromPart(rand.Intn(len(listParts)))

	randomLink := fmt.Sprintf("%s", URLList[rand.Intn(len(URLList))])
	moviePage = fetchPage("https://en.wikiquote.org" + randomLink)
	return
}

// Extracts the "href" attribute of the first a-tag child it finds within a given container node.
// This function assumes every node it checks to have only one child.
// Siblings do not get checked.
func extractLink(node *html.Node) (string, error) {
	// TODO: Make use of querySelector for this, as it does most of what this function does as well.
	if node.Type == html.ElementNode && node.Data == "a" {
		for _, a := range node.Attr {
			if a.Key == "href" {
				findRedlink := regexp.MustCompile("redlink=1")
				if findRedlink.MatchString(a.Val) {
					return "", errors.New("Found link was a red link")
				}
				return a.Val, nil
			}
		}
	}
	if node.FirstChild == nil {
		return "", errors.New("No link found")
	}
	return extractLink(node.FirstChild)
}

// Takes a URL and returns the fetched page as *html.Node.
func fetchPage(URL string) (page *html.Node) {
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