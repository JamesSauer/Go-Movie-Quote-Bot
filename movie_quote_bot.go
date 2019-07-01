package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"strings"
)

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
	page := getRandomMoviePage()
	ioutil.WriteFile("test.html", []byte(page), 0600)

	doc, _ := html.Parse(strings.NewReader(page))
	nodeList := querySelectorAll(doc, "li")

	fmt.Println(len(nodeList))
	for _, node := range nodeList {
		if node.Type == html.ElementNode {
			fmt.Println(node.FirstChild.Data + "\n")
		}
	}
}