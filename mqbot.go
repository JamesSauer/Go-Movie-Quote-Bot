package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"os"
	"strings"
	"time"
)

/* TODO:
- Fix TODOs tied to specific code snippets before tackling the ones below.
- Use gofmt.
- Replace the functions in the dom.go file with goquery. https://github.com/PuerkitoBio/goquery
- Make the scraping more robust by considering edge cases.
  (E.g. quotes following the "Others" heading, where the name of the character appears before the quote.)
- Write some tests.

BONUS:
- Post the quotes somewhere instead of just printing them to the console. (Twitter? Separate web page?)
- Make the selection of quotes more interesting than just selecting them at random. (Maybe find quotes that fit a theme?)
- Maybe fetch some meta data from themoviedb.org and incorporate it into the tweets or web page.
- Look into making the code more idiomatic. Start here: https://golang.org/doc/effective_go.html
*/

func main() {
	// TODO: Add a subcommand for scraping the entirety of the quotes.
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
		case "--fresh", "-f":
			printRandom(getRandomQuoteFresh)
		// Force retrieving a quote from the database:
		case "--database", "-db":
			printRandom(getRandomQuoteDB)
		// Scrape one page at random:
		case "scrape1page":
			page, err := scrapeRandomPage()
			if err != nil {
				log.Fatalln(err)
			}

			db, err = connectPostgres()
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()

			err = page.save()
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Successfully scraped and saved the entry for the movie \"%s\"!\n", page.movie.title)
		// Scrape ALL the pages:
		case "scrapeall":
			warning := "This command will attempt to scrape the entirety of wikiquote.org's movie quotes.\n"+
				"This will take more than 10 minutes."

			if confirm(warning) {
				var err error
				db, err = connectPostgres()
				if err != nil {
					log.Fatalln(err)
				}
				numPages, elapsedTime, err := scrapeAll()
				fmt.Printf("\rScraped %d pages in %s!\n", numPages, elapsedTime)
				if err != nil {
					log.Fatal(err)
				}
				return
			}
			return
		// Using unknown flags or subcommands defaults to behaviour without flags or subcommands:
		default:
			fmt.Println("Movie quote bot doesn't have that command, but here's a random quote instead:")
			printRandom(getRandomQuote)
		}
	}
}

func confirm(warning string) (confirmed bool) {
	fmt.Println(warning + "\n\nDo you want to proceed? (yes/y/no/n)")

	findWords := regexp.MustCompile(`^((?i)yes|y|no|n)\s`)
	keepAsking := true

	for keepAsking {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		match := findWords.FindStringSubmatch(input)
		if len(match) == 2 {
			input = strings.ToLower(match[1])
		}

		switch input {
		case "yes", "y":
			confirmed  = true
			keepAsking = false
			return
		case "no", "n":
			confirmed  = false
			keepAsking = false
			return
		default:
			fmt.Println("You have to type yes, y, no or n.")
			continue
		}
	}
	return 
}

// Attempts to scrape a random movie page from Wikiquote.
func scrapeRandomPage() (page *Page, err error) {
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
	err = row.Scan(&body, &author, &title, &wikiquote_url)
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
	page, err := scrapeRandomPage()
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