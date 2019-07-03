package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"regexp"
)

var (
	connectStr = "postgres://postgres:abc123@localhost:5432/mqbot_dev?sslmode=disable"
	sqlDirName = "sql"
	sqlStatements map[string]string // Gets populated after connection to Postgres is established.
)

type movie struct{
	title, wikiquoteURL string
}

type quote struct{
	text string
	character int
}

// type quote2 struct{
// 	movie, author, body string
// }

func connectPostgres() *sql.DB {
	defer loadSQL()

	db, err := sql.Open("postgres", connectStr)
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

// Loads all .sql files from the sql directory and puts them in a global map.
func loadSQL() {
	sqlStatements = make(map[string]string)

	sqlFiles, err := ioutil.ReadDir(sqlDirName)
	if err != nil {
		log.Fatal(err)
	}

	stripExtension := regexp.MustCompile(`(.*)\.([a-zA-Z0-9]+)$`)

	for _, f := range sqlFiles {

		submatches := stripExtension.FindStringSubmatch(f.Name())
		name := submatches[1]
		extension := submatches[2]

		if extension == "sql" {
			var sqlFileContent, err = ioutil.ReadFile(sqlDirName + "/" + f.Name())
			if err != nil {
				log.Fatal(err)
			}

			sqlStatements[name] = string(sqlFileContent)
		}
	}
	return
}

func executeSchema() {
	schema := sqlStatements["db_schema"]
	statements := regexp.MustCompile(`;\s*`).Split(schema, -1)

	db := connectPostgres()

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			panic(err)
		}
	}
}

func (movie *movie) save() {
	_, err := db.Exec(sqlStatements["insert_movie"], movie.title, movie.wikiquoteURL)
	if err != nil {
		log.Fatal(err)
	}
}