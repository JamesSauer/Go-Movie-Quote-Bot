package main

import (
	"fmt"
)

// TODO: Make the save() methods in this file return errors instead of terminating the program.

// Page ...
type Page struct{
	movie *Movie
	characters []*Character
	quotes []*Quote
}

func (page *Page) save() {
	page.movie.save()

	for _, char := range page.characters {
		char.save()
	}

	for _, quote := range page.quotes {
		quote.save()
	}
}

// Movie ...
type Movie struct{
	title, wikiquoteURL string
}

func (movie *Movie) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_movie"], movie.title, movie.wikiquoteURL)
	return
}

// Character ...
type Character struct{
	name string
}

func (char *Character) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_character"], char.name)
	return
}

// Quote ...
type Quote struct{
	movie *Movie
	author *Character
	body string
}

func (q *Quote) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_quote"], q.body, q.author.name, q.movie.wikiquoteURL)
	return
}

func (q *Quote) saveFull() {
	q.author.save()
	q.movie.save()
	q.save()
}

func (q *Quote) print() {
	fmt.Printf("%s\n    - %s, %s", q.body, q.author.name, q.movie.title)
}