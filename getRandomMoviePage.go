package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
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

var findRedlink = regexp.MustCompile("redlink=1")

// Returns a random WikiQuotes page about a movie.
func getRandomMoviePage() (moviePage string) {
	rand.Seed(time.Now().UnixNano())
	startURL := fmt.Sprintf("https://en.wikiquote.org/wiki/List_of_films_(%s)", listParts[rand.Intn(8)])

	document := stringToDom(fetchPage(startURL))

	movieLinks := make([]string, 0)
	for _, node := range querySelectorAll(document, "i") {
		link, err := extractLink(node)
			if err == nil {
				movieLinks = append(movieLinks, link)
			}
	}

	randomLink := fmt.Sprintf("%s", movieLinks[rand.Intn(len(movieLinks))])
	moviePage = fetchPage("https://en.wikiquote.org" + randomLink)
	return
}



func extractLink(n *html.Node) (string, error) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				if findRedlink.MatchString(a.Val) {
					return "", errors.New("Found link was a red link")
				}
				return a.Val, nil
			}
		}
	}
	if n.FirstChild == nil {
		return "", errors.New("No link found")
	}
	return extractLink(n.FirstChild)
}

func fetchPage(URL string) (page string) {
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	page = string(body)
	return
}