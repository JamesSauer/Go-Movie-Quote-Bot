package main

import (
	"fmt"
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
		getRandomQuote().print()
	} else {
		switch os.Args[1] {
		case "test":
			fmt.Println("Nope! Chuck Testa!")
		case "testdb":
			connectPostgres()
			defer db.Close()
		case "save1":
			connectPostgres()
			defer db.Close()

			q := getRandomQuote()
			q.saveFull()
		default:
			fmt.Println("Movie quote bot doesn't have that command, but here's a random quote instead:")
			getRandomQuote().print()
		}
	}
}

func getRandomQuote() (q *Quote) {
	rand.Seed(time.Now().UnixNano())
	quotes, _, _ := scrapeQuotes(getRandomURL())
	for len(quotes) == 0 {
		quotes, _, _ = scrapeQuotes(getRandomURL())
	}
	q = quotes[rand.Intn(len(quotes))]
	return
}