package main

import (
	"fmt"
)

// TODO: Make the save() methods in this file return errors instead of terminating the program.

// Page ...
type Page struct {
	movie      *Movie
	characters []*Character
	quotes     []*Quote
}

func (page *Page) save() (err error) {
	err = page.movie.save()

	for _, char := range page.characters {
		err = char.save()
	}

	for _, quote := range page.quotes {
		err = quote.save()
	}

	// Returning just the last error that occured is sufficient for now.
	return
}

// Movie ...
type Movie struct {
	title, wikiquoteURL string
}

func (movie *Movie) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_movie"], movie.title, movie.wikiquoteURL)
	return
}

// Character ...
type Character struct {
	name string
}

func (char *Character) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_character"], char.name)
	return
}

// Quote ...
type Quote struct {
	movie  *Movie
	author *Character
	body   string
}

func (q *Quote) save() (err error) {
	_, err = db.Exec(sqlStatements["insert_quote"], q.body, q.author.name, q.movie.wikiquoteURL)
	return
}

func (q *Quote) saveFull() (err error) {
	err = q.author.save()
	err = q.movie.save()
	err = q.save()
	// Returning just the last error that occured is sufficient for now.
	return
}

func (q *Quote) print() {
	fmt.Printf("%s\n	- %s, %s", q.body, q.author.name, q.movie.title)
}

func getStats() (stats map[string]int, err error) {
	db, err = connectPostgres()
	if err != nil {
		return
	}
	defer db.Close()

	result, err := db.Query(sqlStatements["select_stats"])
	if err != nil {
		return
	}

	stats = make(map[string]int)
	for result.Next() {
		var (
			stat string
			value int
		)
		result.Scan(&stat, &value)
		stats[stat] = value
	}
	err = result.Err()
	return
}