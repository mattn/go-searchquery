package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	dialect "github.com/mattn/go-searchquery/dialect/meilisearch"
	"github.com/meilisearch/meilisearch-go"
)

func main() {
	var apikey string
	var endpoint string
	var catelog string
	var filter string
	flag.StringVar(&apikey, "apikey", "foobar", "API key for Meilisearch")
	flag.StringVar(&endpoint, "endpoint", "http://localhost:7700", "Meilisearch endpoint")
	flag.StringVar(&catelog, "index", "movies", "Index name")
	flag.StringVar(&filter, "filter", "", "Filter expression")
	flag.Parse()
	client := meilisearch.New(endpoint, meilisearch.WithAPIKey(apikey))

	var query string
	var err error
	if filter != "" {
		query, err = dialect.ToFilter(filter, flag.Arg(0))
	} else {
		query, err = dialect.ToQuery(flag.Arg(0))
	}
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Index(catelog).Search(query, &meilisearch.SearchRequest{
		Limit:            10,
		MatchingStrategy: meilisearch.All,
	})
	if err != nil {
		log.Fatal(err)
	}
	var items []map[string]any
	err = resp.Hits.DecodeInto(&items)
	if err != nil {
		log.Fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	for _, item := range items {
		enc.Encode(item)
	}
}
