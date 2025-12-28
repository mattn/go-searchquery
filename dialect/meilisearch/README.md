# Meilisearch Dialect

Converts search queries to Meilisearch query or filter format.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/meilisearch
```

## Usage

### Query Format (for search)

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery/dialect/meilisearch"
)

func main() {
    userQuery := "hello world"
    query, err := meilisearch.ToQuery(userQuery)
    if err != nil {
        panic(err)
    }
    // query = "hello world"
    
    // Use in Meilisearch search
    // POST /indexes/articles/search
    // {
    //   "q": "hello world"
    // }
    
    // Phrase search
    phraseQuery := `"hello world"`
    query, _ = meilisearch.ToQuery(phraseQuery)
    // query = "\"hello world\""
}
```

### Filter Format (for exact matching)

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery/dialect/meilisearch"
)

func main() {
    userQuery := "hello world"
    filter, err := meilisearch.ToFilter("title", userQuery)
    if err != nil {
        panic(err)
    }
    // filter = "(title = \"hello\" AND title = \"world\")"
    
    // Use in Meilisearch filter
    // POST /indexes/articles/search
    // {
    //   "filter": "(title = \"hello\" AND title = \"world\")"
    // }
}
```

## Query Conversion Examples

| Input Query | Query Format | Filter Format |
|-------------|--------------|---------------|
| `hello world` | `hello world` | `(title = "hello" AND title = "world")` |
| `"hello world"` | `"hello world"` | `title = "hello world"` |
| `cat dog bird` | `cat dog bird` | Nested AND filters |

## Features

- **Simple Query Syntax**: Natural, easy-to-read format
- **Two Modes**: Query for fuzzy search, Filter for exact matching
- **Phrase Support**: Quoted phrases preserved
- **OR Support**: Explicit OR operators work as expected

## Notes

- Meilisearch automatically handles typos and relevance ranking
- Query format uses Meilisearch's natural search capabilities
- Filter format is for precise filtering (no typo tolerance)
- Very fast with minimal configuration
