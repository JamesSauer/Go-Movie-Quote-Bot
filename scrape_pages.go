package main

import(
	"regexp"
)

// Takes a wikiquote URL to a movie page and returns a slice of quotes,
// a slice of characters the quotes are attributed to as well as the title of the movie.
func scrapePage(movieURL string) (page *Page, err error) {
	document, err := fetch(movieURL)
	if err != nil {
		return
	}
	movie := &Movie{
		wikiquoteURL: regexp.MustCompile(`\.org(.+)`).FindStringSubmatch(movieURL)[1],
		title: extractText(querySelectorAll(document, "#firstHeading")[0]),
	}

	characters := make([]*Character, 0)
	quotes := make([]*Quote, 0)
	
	headings := querySelectorAll(document, ".mw-headline")
	
	for _, heading := range headings {
		char := extractText(heading)
		if !isCharacter(char) {
			continue
		}
		characters = append(characters, &Character{name: char})

		ul := getNextElementSibling(heading.Parent)
		items := querySelectorAll(ul, "li")

		for _, item := range items {
			q := &Quote{body: extractText(item), movie: movie, author: characters[len(characters)-1]}
			quotes = append(quotes, q)
		}
	}
	page = &Page{movie: movie, characters: characters, quotes: quotes}
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