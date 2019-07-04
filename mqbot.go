package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

/* TODO:
- Fix TODOs tied to specific code snippets before tackling the ones below.
- Use gofmt.
- Split off querySelectors and other DOM helpers into a separate package.
- Make the scraping more robust by considering edge cases.
  (E.g. quotes following the "Others" heading, where the name of the character appears before the quote.)
- Write some tests.

BONUS:
- Post the quotes somewhere instead of just printing them to the console. (Twitter? Separate web page?)
- Persist the quotes in a database.
- Make it an actual bot instead of a command line tool by letting it scrape continously.
- Make the selection of quotes more interesting than just selecting them at random. (Maybe find quotes that fit a theme?)
- Maybe fetch some meta data from themoviedb.org and incorporate it into the tweets or web page.
- Look into making the code more idiomatic. Start here: https://golang.org/doc/effective_go.html
*/

func main() {
	if len(os.Args) <= 1 {
		printRandom(getRandomQuote)
	} else {
		switch os.Args[1] {
		case "test":
			fmt.Println("Nope! Chuck Testa!")
		case "testdb":
			db, err := connectPostgres()
			if err != nil {
				log.Fatalln(errors.New("Could not connect to DB. Please ensure the MQBOT_POSTGRES environment variable is set correctly"))
			}
			fmt.Println("Successfully connected to DB!")
			db.Close()
		// Force scraping a fresh quote:
		case "--fresh":
			printRandom(getRandomQuoteFresh)
		case "-f":
			printRandom(getRandomQuoteFresh)
		// Force retrieving a quote from the database:
		case "--database":
			printRandom(getRandomQuoteDB)
		case "-db":
			printRandom(getRandomQuoteDB)
		default:
			fmt.Println("Movie quote bot doesn't have that command, but here's a random quote instead:")
			printRandom(getRandomQuote)
		}
	}
}

// Attempts to scrape a random movie page from Wikiquote.
func getRandomPage() (page *Page, err error) {
	rand.Seed(time.Now().UnixNano())
	url, err := getRandomURL()
	if err != nil {
		return
	}

	page, err = scrapePage(url)
	if err != nil {
		return
	}
	for len(page.quotes) == 0 {
		page, err = scrapePage(url)
		if err != nil {
			return
		}
	}
	return
}

// Attempts to retrieve a random quote from the database.
// If it fails, fetch a fresh one from Wikiquote instead.
// This is the default behaviour.
func getRandomQuote() (quote *Quote, err error) {
	quote, err = getRandomQuoteDB()
	if err != nil {
		quote, err = getRandomQuoteFresh()
		if err != nil {
			return nil, errors.New("Couldn't retrieve a quote from either the database or Wikiquote")
		}
		return
	}
	return
}

// Attempts to retrieve a random quote from the database.
func getRandomQuoteDB() (quote *Quote, err error) {
	db, err := connectPostgres()
	if err != nil {
		return
	}
	defer db.Close()

	var (
		body string
		author string
		title string
		wikiquote_url string
	)
	row := db.QueryRow(sqlStatements["select_random_quote"])
	err = row.Scan(&body, &author, &title)
	if err != nil {
		return
	}
	movie := &Movie{
		title: title,
		wikiquoteURL: wikiquote_url,
	}
	char := &Character{
		name: author,
	}
	quote = &Quote{
		movie: movie,
		author: char,
		body: body,
	}
	return
}

// Attempts to scrape Wikiquote for a random quote.
func getRandomQuoteFresh() (quote *Quote, err error) {
	page, err := getRandomPage()
	if err != nil {
		return
	}
	rand.Seed(time.Now().UnixNano())
	quote = page.quotes[rand.Intn(len(page.quotes))]
	return
}

// Wrapper for the getRandomQuote...() functions to handle errors.
func printRandom(fn func() (*Quote, error)) {
	q, err := fn()
	if err != nil {
		log.Fatalln(err)
	}
	q.print()
}