package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	dialect "github.com/mattn/go-searchquery/dialect/mysql"
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

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if initFlag {
		_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS example (
            id INT AUTO_INCREMENT PRIMARY KEY,
            data TEXT NOT NULL,
            FULLTEXT KEY ft_data (data)
        ) ENGINE=InnoDB`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("TRUNCATE TABLE example")
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare("INSERT INTO example(data) VALUES(?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		// Note: MySQL InnoDB fulltext search has a minimum token size (default 3 characters).
		// Words shorter than this (like "go") won't be indexed.
		// Also, common stopwords are ignored by default.
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
		rows, err = db.Query(
			"SELECT e.id, e.data FROM example e " +
				"ORDER BY e.id")
	} else {
		query, err := dialect.ToBoolean(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = db.Query(
			"SELECT id, data FROM example "+
				"WHERE MATCH(data) AGAINST(? IN BOOLEAN MODE) "+
				"ORDER BY id", query)
		if err != nil {
			log.Fatal(err)
		}
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
