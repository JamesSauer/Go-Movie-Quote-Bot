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
	"strings"
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

// Returns a random WikiQuotes page about a movie.
func getRandomMoviePage() (moviePage string) {
	rand.Seed(time.Now().UnixNano())
	startURL := fmt.Sprintf("https://en.wikiquote.org/wiki/List_of_films_(%s)", listParts[rand.Intn(8)])

	listPage := fetchPage(startURL)
	doc, err := html.Parse(strings.NewReader(listPage))
	if err != nil {
		log.Fatal(err)
	}

	movieLinks := make([]string, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "i" {
			link, err := extractLink(n)
			if err == nil {
				movieLinks = append(movieLinks, link)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	randomLink := fmt.Sprintf("%s", movieLinks[rand.Intn(len(movieLinks))])
	moviePage = fetchPage("https://en.wikiquote.org" + randomLink)
	return
}

// Takes a wiki quotes page as string and returns a map of slices of quotes on the page.
// The keys of the map are the characters the quotes are attributed to.
func getQuotes(moviePage string) (quotes map[string][]string) {
	// TODO: Implement function.
	quotes = make(map[string][]string)
	return
}

func main() {
	// TODO: Select random quote.
	fmt.Println(`My momma always said, "Life was like a box of chocolates. You never know what you're gonna get."`)
	ioutil.WriteFile("test.html", []byte(getRandomMoviePage()), 0600)
}