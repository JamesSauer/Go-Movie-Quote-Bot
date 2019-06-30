package main

import (
	"fmt"
	"io/ioutil"
)

// Returns a random movie title.
func getRandomMovie() (movieTitle string) {
	// TODO: Where to look for movies? themoviedb.org API?
	// TODO: Implement function.
	movieTitle = "Forrest Gump"
	return
}

// Takes a movie title, searches WikiQuotes for it and returns its entry.
// If no entry can be found, returns an error instead.
func getMoviePage(movieTitle string) (moviePage string, err error) {
	// TODO: Implement function.
	fileContent, err := ioutil.ReadFile("test_page.html")
	moviePage = string(fileContent)
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
}