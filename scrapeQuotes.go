package main

import(
	"golang.org/x/net/html"
)

var nonCharacterHeadings = [7]string{
	"Contents",
	"Dialogue",
	"Cast",
	"External links",
	"Navigation menu",
	"Taglines",
	"See also",
}

type quote struct{
	text string
	character int
}

// Takes a wiki quotes page as string and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapeQuotes(moviePage string) (quotes []quote, characters []string, title string) {
	document := stringToDom(moviePage)

	quotes = make([]quote, 0)
	characters = make([]string, 0)
	title = querySelectorAll(document, "#firstHeading")[0].FirstChild.FirstChild.Data

	headings := querySelectorAll(document, ".mw-headline")
	
	for _, heading := range headings {
		character := heading.FirstChild.Data
		if !isCharacterHeading(character) {
			continue
		}
		characters = append(characters, character)
		ul := getNextElementSibling(heading.Parent)

		

		items := querySelectorAll(ul, "li")

		for _, item := range items {
			q := quote{text: extractText(item), character: len(characters) - 1}
			quotes = append(quotes, q)
		}
	}
	return
}

func extractText(root *html.Node) (text string) {
	var walker func(*html.Node)
	walker = func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			walker(c)
		}
	}
	walker(root)
	return
}

func getNextElementSibling(node *html.Node) (sibling *html.Node) {
	for sibling = node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type == html.ElementNode {
			return
		}
	}
	return nil
}

func isCharacterHeading(title string) bool {
	for _, v := range nonCharacterHeadings {
		if v == title {
			return false
		}
	}
	return true
}