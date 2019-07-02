package main

import(
	"golang.org/x/net/html"
)

type quote struct{
	text string
	character int
}

// Takes a wikiquote page as *html.Node and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapeQuotes(document *html.Node) (quotes []quote, characters []string, title string) {
	quotes = make([]quote, 0)
	characters = make([]string, 0)
	title = extractText(querySelectorAll(document, "#firstHeading")[0])

	headings := querySelectorAll(document, ".mw-headline")
	
	for _, heading := range headings {
		character := extractText(heading)
		if !isCharacter(character) {
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

// Strips the tags off the text within a given node.
// Without this, inline tags like <b>...</b> would screw up the quotes.
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

// Because html.Node.NextSibling doesn't differentiate between text and element nodes.
func getNextElementSibling(node *html.Node) (sibling *html.Node) {
	for sibling = node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type == html.ElementNode {
			return
		}
	}
	return nil
}

func getFirstElementChild(node *html.Node) (firstChild *html.Node) {
	// TODO: This function does almost the exact same thing as the one above. Merge them!
	for firstChild = node.FirstChild; firstChild != nil; firstChild = firstChild.NextSibling {
		if firstChild.Type == html.ElementNode {
			return
		}
	}
	return nil
}

var nonCharacterHeadings = [8]string{
	"About",
	"Cast",
	"Contents",
	"Dialogue",
	"External links",
	"Navigation menu",
	"See also",
	"Taglines",
}

func isCharacter(title string) bool {
	for _, v := range nonCharacterHeadings {
		if v == title {
			return false
		}
	}
	return true
}