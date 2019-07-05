package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	// "strings"
)

var (
	db            *sql.DB
	connectStr    = os.Getenv("MQBOT_POSTGRES")
	fallbackStr   = "postgres://postgres:abc123@localhost:5432/mqbot_dev?sslmode=disable"
	sqlDirName    = "sql"
	sqlStatements map[string]string // Gets populated after connection to Postgres is established.
)

func connectPostgres() (connection *sql.DB, err error) {
	defer loadSQL()

	if connectStr == "" {
		connectStr = fallbackStr
	}

	connection, err = sql.Open("postgres", connectStr)
	if err != nil {
		return
	}

	err = connection.Ping()
	if err != nil {
		return
	}
	return
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
				log.Fatalln(err)
			}

			sqlStatements[name] = string(sqlFileContent)
		}
	}
	return
}

func executeSchema() {
	var err error
	if db == nil {
		db, err = connectPostgres()
		defer db.Close()
	}
	if err != nil {
		log.Fatalln(err)
	}

	schema := sqlStatements["db_schema"]
	statements := regexp.MustCompile(`;\s*`).Split(schema, -1)

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

// TODO: Make it work with atomic queries first. Ignore the below and worry about optimization later.

// func insertMovieAndCharacters(film *movie, characters []*character) {
// 	names := make([]string, 0)
// 	for _, char := range characters {
// 		names = append(names, char.name)
// 	}
// 	rows, err := db.Query(sqlStatements["insert_movie_and_characters"], film.wikiquoteURL, film.title, pgArrayLiteral(names))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// }

// func pgArrayLiteral(slice []string) string {
// 	return "{'" + strings.Join(slice, "','") + "'}"
// }

// func insertQuoteBatch(batch []*quote) {
// 	txn, err := db.Begin()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	stmt, _ := txn.Prepare(pq.CopyIn("quotes", "movie", "author", "body"))

// 	for _, quote := range batch {
// 		_, err := stmt.Exec(quote.movie, quote.author, quote.body)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	_, err = stmt.Exec()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = stmt.Close()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = txn.Commit()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
