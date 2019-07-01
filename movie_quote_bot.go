package main

import (
	"fmt"
	"math/rand"
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
	quotes, characters, title := scrapeQuotes(getRandomMoviePage())
	for len(quotes) == 0 {
		quotes, characters, title = scrapeQuotes(getRandomMoviePage())
	}
	q := quotes[rand.Intn(len(quotes))]
	fmt.Printf("%s\n    - %s, %s", q.text, characters[q.character], title)
}