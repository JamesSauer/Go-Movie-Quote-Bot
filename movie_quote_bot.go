package main

import (
	"fmt"
	"math/rand"
)

func main() {
	quotes, characters, title := scrapeQuotes(getRandomMoviePage())
	for len(quotes) == 0 {
		quotes, characters, title = scrapeQuotes(getRandomMoviePage())
	}
	q := quotes[rand.Intn(len(quotes))]
	fmt.Printf("%s\n    - %s, %s", q.text, characters[q.character], title)
}