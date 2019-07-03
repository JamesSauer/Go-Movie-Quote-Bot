package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"regexp"
)

func connectPostgres() *sql.DB {
	connStr := "postgres://postgres:abc123@localhost:5432/mqbot_dev?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
	  panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func executeSchema() {
	file, err := ioutil.ReadFile("db_schema.sql")
	if err != nil {
		panic(err)
	}
	schema := string(file)
	statements := regexp.MustCompile(`;\s*`).Split(schema, -1)

	db := connectPostgres()

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}