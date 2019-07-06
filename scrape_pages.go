package main

import (
	"errors"
	"fmt"
	"log"
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
		title:        extractText(querySelectorAll(document, "#firstHeading")[0]),
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
	warning := "This command will attempt to scrape the entirety of wikiquote.org's movie quotes.\n" +
		"This might take more than 10 minutes. Do you want to proceed?"

	if confirm(warning) {
		if db == nil {
			db, err = connectPostgres()
			defer db.Close()
		}
		if err != nil {
			log.Fatalln(err)
		}

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
			err = errors.New("scrapeAll expects one argument at most that specifies the number of requests per second as float64")
		}
		timeout := int(math.Ceil(1000 / reqsPerSec))

		// Start scraping
		start := time.Now()
		var urls []string
		urls, err = getAllURLs()
		if err != nil {
			return
		}
		numScrapeErrors := 0
		numDBErrors := 0
		runIndicator := [...]string{".....", " ....", ". ...", ".. ..", "... .", ".... "}
		for i, url := range urls {
			fmt.Printf("\rScraping page %d of %d %s", i+1, len(urls), runIndicator[i%len(runIndicator)])
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
	}
	return
}

var nonCharacterHeadings = [...]string{
	"About",
	"Cast",
	"Contents",
	"Dialogue",
	"External links",
	"Footnote",
	"Navigation menu",
	"See also",
	"Taglines",
	"Voice cast",
}

func isCharacter(title string) bool {
	for _, v := range nonCharacterHeadings {
		if v == title {
			return false
		}
	}
	return true
}
