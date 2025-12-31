package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
	dialect "github.com/mattn/go-searchquery/dialect/bleve"
)

func main() {
	var initFlag bool
	var listFlag bool
	flag.BoolVar(&initFlag, "init", false, "Initialize the index")
	flag.BoolVar(&listFlag, "list", false, "List all items")
	flag.Parse()

	if flag.NArg() < 1 && !initFlag && !listFlag {
		log.Fatal("Please provide a search query or use -init to initialize the index, or -list to list all items")
	}

	indexPath := "example.bleve"

	if initFlag {
		// Remove existing index if it exists
		os.RemoveAll(indexPath)

		// Create a new index
		mapping := bleve.NewIndexMapping()
		index, err := bleve.New(indexPath, mapping)
		if err != nil {
			log.Fatal(err)
		}

		// Index some sample documents
		documents := []struct {
			ID   string
			Text string
		}{
			{"1", "Hello World"},
			{"2", "Great World"},
			{"3", "Hello Go"},
			{"4", "Golang is awesome"},
			{"5", "Bleve search engine"},
		}

		for _, doc := range documents {
			err = index.Index(doc.ID, map[string]interface{}{
				"text": doc.Text,
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		index.Close()
		fmt.Println("Index initialized successfully")
		return
	}

	// Open the existing index
	index, err := bleve.Open(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	defer index.Close()

	var sq query.Query
	if listFlag {
		sq = bleve.NewMatchAllQuery()
	} else {
		// Convert the search query to Bleve format
		queryString, err := dialect.ToQueryString(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// Create a query string query
		sq = bleve.NewQueryStringQuery(queryString)
	}

	searchRequest := bleve.NewSearchRequest(sq)
	searchRequest.Fields = []string{"text"}

	// Execute the search
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Display results
	for i, hit := range searchResults.Hits {
		fmt.Printf("%d. Document ID: %s, Score: %f\n", i+1, hit.ID, hit.Score)
		if text, ok := hit.Fields["text"].(string); ok {
			fmt.Printf("   Text: %s\n", text)
		}
	}
}
