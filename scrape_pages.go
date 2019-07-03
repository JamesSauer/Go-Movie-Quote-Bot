package main

import(
	"regexp"
)

// Takes a wikiquote URL to a movie page and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapeQuotes(movieURL string) (quotes []*quote, characters []*character, film *movie) {
	document := fetch(movieURL)
	film = &movie{
		wikiquoteURL: regexp.MustCompile(`\.org(.+)`).FindStringSubmatch(movieURL)[1],
		title: extractText(querySelectorAll(document, "#firstHeading")[0]),
	}

	quotes = make([]*quote, 0)
	characters = make([]*character, 0)
	
	headings := querySelectorAll(document, ".mw-headline")
	
	for _, heading := range headings {
		char := extractText(heading)
		if !isCharacter(char) {
			continue
		}
		characters = append(characters, &character{name: char})

		ul := getNextElementSibling(heading.Parent)
		items := querySelectorAll(ul, "li")

		for _, item := range items {
			q := &quote{body: extractText(item), movie: film, author: characters[len(characters)-1]}
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