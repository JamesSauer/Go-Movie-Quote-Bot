package main

import (
	"fmt"
	"log"
)

// Movie ...
type Movie struct{
	title, wikiquoteURL string
}

func (movie *Movie) save() {
	_, err := db.Exec(sqlStatements["insert_movie"], movie.title, movie.wikiquoteURL)
	if err != nil {
		log.Fatal(err)
	}
}

// Character ...
type Character struct{
	name string
}

func (char *Character) save() {
	_, err := db.Exec(sqlStatements["insert_character"], char.name)
	if err != nil {
		log.Fatal(err)
	}
}

// Quote ...
type Quote struct{
	movie *Movie
	author *Character
	body string
}

func (q *Quote) save() {
	_, err := db.Exec(sqlStatements["insert_quote"], q.body, q.author.name, q.movie.wikiquoteURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (q *Quote) saveFull() {
	q.author.save()
	q.movie.save()
	q.save()
}

func (q *Quote) print() {
	fmt.Printf("%s\n    - %s, %s", q.body, q.author.name, q.movie.title)
}