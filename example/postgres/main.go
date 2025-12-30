package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	dialect "github.com/mattn/go-searchquery/dialect/postgres"
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

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if initFlag {
		_, err = db.Exec("CREATE TABLE IF NOT EXISTS example (id SERIAL PRIMARY KEY, data TEXT)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE TABLE IF NOT EXISTS example_tsvector (id INTEGER PRIMARY KEY, data_tsv tsvector)")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(
			"CREATE OR REPLACE FUNCTION example_ai() RETURNS TRIGGER AS $$ BEGIN " +
				"INSERT INTO example_tsvector(id, data_tsv) VALUES (NEW.id, to_tsvector('simple', NEW.data)) " +
				"ON CONFLICT (id) DO UPDATE SET data_tsv = to_tsvector('simple', NEW.data); RETURN NEW; END; $$ LANGUAGE plpgsql")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("DROP TRIGGER IF EXISTS example_ai_trigger ON example")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("CREATE TRIGGER example_ai_trigger AFTER INSERT ON example FOR EACH ROW EXECUTE FUNCTION example_ai()")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("TRUNCATE TABLE example RESTART IDENTITY")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("TRUNCATE TABLE example_tsvector")
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare("INSERT INTO example(data) VALUES($1)")
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
	}

	var rows *sql.Rows

	if listFlag {
		rows, err = db.Query(
			"SELECT e.id, e.data FROM example e " +
				"ORDER BY e.id")
	} else {
		query, err := dialect.ToTsQuery(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		rows, err = db.Query(
			"SELECT e.id, e.data FROM example e JOIN example_tsvector f ON e.id = f.id "+
				"WHERE f.data_tsv @@ to_tsquery('simple', $1) "+
				"ORDER BY e.id", query)
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
