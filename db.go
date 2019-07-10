package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
			var sqlFileContent, err = ioutil.ReadFile(filepath.Join(sqlDirName, f.Name()))
			if err != nil {
				log.Fatalln(err)
			}

			sqlStatements[name] = string(sqlFileContent)
		}
	}
	return
}

// Set up database schema.
func executeSchema() (err error) {
	if db == nil {
		db, err = connectPostgres()
		defer db.Close()
	}
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.Exec(sqlStatements["db_schema"])
	return
}
