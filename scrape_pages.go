package main

type quote struct{
	text string
	character int
}

// Takes a wikiquote URL to a movie page and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapeQuotes(movieURL string) (quotes []quote, characters []string, title string) {
	document := fetch(movieURL)

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