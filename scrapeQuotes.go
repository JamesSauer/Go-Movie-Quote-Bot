package main

import(
	"github.com/PuerkitoBio/goquery"
)

type quote struct{
	text string
	character int
}

// Takes a wikiquote page as *html.Node and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapeQuotes(movieURL string) (quotes []quote, characters []string, title string) {
	document := fetchPage(movieURL)

	quotes = make([]quote, 0)
	characters = make([]string, 0)
	title = extractText(document.Find("#firstHeading"))

	document.Find(".mw-headline").Each(func(i int, heading *goquery.Selection) {
		character := extractText(heading)
		if isCharacter(character) {
			characters = append(characters, character)
			ul := heading.Parent().Next()
			items := ul.ChildrenFiltered("li")
			items.Each(func(i int, item *goquery.Selection) {
				text := extractText(item)
				q := quote{text: text, character: len(characters) - 1}
				quotes = append(quotes, q)
			})
		}
	})

	return
}

// Strips the tags off the text within a given node.
// Without this, inline tags like <b>...</b> would screw up the quotes.
func extractText(node *goquery.Selection) (text string) {
	node.Children().Each(func(i int, child *goquery.Selection) {
		nodeName := goquery.NodeName(child)
		if nodeName == "#text" {
			innerHTML, _ := child.Html()
			text += innerHTML
		}
	})
	return
}


var nonCharacterHeadings = [7]string{
	"Contents",
	"Dialogue",
	"Cast",
	"External links",
	"Navigation menu",
	"Taglines",
	"See also",
}

func isCharacter(title string) bool {
	for _, v := range nonCharacterHeadings {
		if v == title {
			return false
		}
	}
	return true
}