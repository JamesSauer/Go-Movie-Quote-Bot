package main

import(
	"errors"
	"fmt"
	"math"
	"regexp"
	"time"
)

// Takes a wikiquote URL to a movie page and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapePage(movieURL string) (page *Page, err error) {
	document, err := fetch(movieURL)
	if err != nil {
		return
	}
	movie := &Movie{
		wikiquoteURL: regexp.MustCompile(`\.org(.+)`).FindStringSubmatch(movieURL)[1],
		title: extractText(querySelectorAll(document, "#firstHeading")[0]),
	}

	characters := make([]*Character, 0)
	quotes := make([]*Quote, 0)
	
	headings := querySelectorAll(document, ".mw-headline")
	
	for _, heading := range headings {
		char := extractText(heading)
		if !isCharacter(char) {
			continue
		}
		characters = append(characters, &Character{name: char})

		ul := getNextElementSibling(heading.Parent)
		items := querySelectorAll(ul, "li")

		for _, item := range items {
			q := &Quote{body: extractText(item), movie: movie, author: characters[len(characters)-1]}
			quotes = append(quotes, q)
		}
	}
	page = &Page{movie: movie, characters: characters, quotes: quotes}
	return
}

func scrapeAll(params ...float64) (numScraped int, took time.Duration, err error) {
	var reqsPerSec float64
	switch len(params) {
	case 0:
		reqsPerSec = 5.0
	case 1:
		if params[0] > 100.0 {
			reqsPerSec = 100.0
		} else if params[0] < 0.1 {
			reqsPerSec = 0.1
		} else {
			reqsPerSec = params[0]
		}
	default:
		err = errors.New("scrapeAll expects one argument at most that specifies the requests per second as float64")
	}
	timeout := int(math.Ceil(1000 / reqsPerSec))

	start := time.Now()
	urls, err := getAllURLs()
	if err != nil {
		return
	}

	numScrapeErrors := 0
	numDBErrors := 0
	for i, url := range urls {
		fmt.Printf("\rCurrently scraping page %d.", i+1)
		page, err := scrapePage("https://en.wikiquote.org" + url)
		if err != nil {
			numScrapeErrors++
		}

		err = page.save()
		if err != nil {
			numDBErrors++
		}

		time.Sleep(time.Duration(timeout) * time.Millisecond)
	}

	numScraped = len(urls)
	took = time.Since(start)

	if numScrapeErrors == 0 && numDBErrors == 0 {
		return
	}

	errStr := "Encountered %d errors while scraping and %d errors while writing to the database"
	err = fmt.Errorf(errStr, numScrapeErrors, numDBErrors)
	return
}

var nonCharacterHeadings = [...]string{
	"About",
	"Cast",
	"Contents",
	"Dialogue",
	"External links",
	"Navigation menu",
	"See also",
	"Taglines",
	"Footnote",
}

func isCharacter(title string) bool {
	for _, v := range nonCharacterHeadings {
		if v == title {
			return false
		}
	}
	return true
}