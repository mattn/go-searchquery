package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	dialect "github.com/mattn/go-searchquery/dialect/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var initFlag bool
	var listFlag bool
	flag.BoolVar(&initFlag, "init", false, "Initialize the database")
	flag.BoolVar(&listFlag, "list", false, "List all items")
	flag.Parse()

	if flag.NArg() < 1 && !initFlag && !listFlag {
		log.Fatal("Please provide a search query or use -init to initialize the index, or -list to list all items")
	}

	db, err := sql.Open("sqlite3", "database.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if initFlag {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS example (id INTEGER PRIMARY KEY, data TEXT)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE VIRTUAL TABLE IF NOT EXISTS example_fts USING fts5(data, content='example', content_rowid='id')")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE TRIGGER IF NOT EXISTS example_ai AFTER INSERT ON example BEGIN INSERT INTO example_fts(rowid, data) VALUES (new.id, new.data); END;")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("DELETE FROM example")
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare("INSERT INTO example(data) VALUES(?)")
		if err != nil {
			log.Fatal(err)
		}

		data := []string{"Hello World", "Great World", "Hello Go", "Golang programming", "Rust language"}
		for _, d := range data {
			_, err = stmt.Exec(d)
			if err != nil {
				log.Fatal(err)
			}
		}

		return
	}

	var rows *sql.Rows

	if listFlag {
		rows, err = db.Query("SELECT rowid, data FROM example_fts ORDER BY rowid")
	} else {
		query, err := dialect.ToFTS5Query(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		rows, err = db.Query("SELECT rowid, data FROM example_fts WHERE data MATCH ? ORDER BY rowid", query)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var text string
		err = rows.Scan(&id, &text)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Data: %s\n", id, text)
	}
}
