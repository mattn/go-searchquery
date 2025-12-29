# Bleve Dialect

This package converts search queries to [Bleve](https://github.com/blevesearch/bleve) query string format.

## Overview

Bleve uses a query string syntax similar to Lucene. This dialect converts the parsed AST into Bleve-compatible query strings.

## Functions

### `ToQueryString(query string, opts ...searchquery.ParserOption) (string, error)`

Converts a search query to Bleve query string format.

**Features:**
- Single terms are passed through as-is
- Implicit AND is converted to `+term1 +term2` (required terms)
- Explicit OR is converted to `term1 term2` (optional terms)
- Phrase searches are preserved with quotes: `"hello world"`
- Special characters are escaped when necessary

## Examples

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery/dialect/bleve"
)

func main() {
    // Single term
    query, _ := bleve.ToQueryString("hello")
    fmt.Println(query) // Output: hello
    
    // Implicit AND (multiple required terms)
    query, _ = bleve.ToQueryString("hello world")
    fmt.Println(query) // Output: +hello +world
    
    // Explicit OR (optional terms)
    query, _ = bleve.ToQueryString("hello OR world")
    fmt.Println(query) // Output: hello world
    
    // Phrase search
    query, _ = bleve.ToQueryString(`"hello world"`)
    fmt.Println(query) // Output: "hello world"
    
    // Complex query
    query, _ = bleve.ToQueryString("cat AND (dog OR mouse)")
    fmt.Println(query) // Output: +cat +(dog mouse)
}
```

### Using with Bleve

```go
package main

import (
    "log"
    
    "github.com/blevesearch/bleve/v2"
    "github.com/mattn/go-searchquery/dialect/bleve"
)

func main() {
    // Open or create a Bleve index
    index, err := bleve.Open("example.bleve")
    if err != nil {
        log.Fatal(err)
    }
    defer index.Close()
    
    // Convert search query to Bleve format
    userQuery := "golang AND search"
    bleveQuery, err := bleve.ToQueryString(userQuery)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create a Bleve query
    query := bleve.NewQueryStringQuery(bleveQuery)
    searchRequest := bleve.NewSearchRequest(query)
    
    // Execute search
    searchResults, err := index.Search(searchRequest)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for _, hit := range searchResults.Hits {
        fmt.Printf("Document: %s, Score: %f\n", hit.ID, hit.Score)
    }
}
```

### With Token Lookup

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery"
    "github.com/mattn/go-searchquery/dialect/bleve"
)

func main() {
    // Define aliases
    aliases := map[string]string{
        "go":   "golang",
        "js":   "javascript",
    }
    
    // Convert with alias replacement
    query := "go OR js"
    result, _ := bleve.ToQueryString(query, searchquery.WithLookup(func(token string) string {
        if alias, ok := aliases[token]; ok {
            return alias
        }
        return token
    }))
    
    fmt.Println(result) // Output: golang javascript
}
```

## Query Syntax Mapping

| Input | Output | Description |
|-------|--------|-------------|
| `hello` | `hello` | Single term |
| `hello world` | `+hello +world` | Implicit AND (required terms) |
| `hello OR world` | `hello world` | Explicit OR (optional terms) |
| `"hello world"` | `"hello world"` | Phrase search |
| `cat AND dog` | `+cat +dog` | Explicit AND |
| `(cat OR dog)` | `(cat dog)` | Grouped OR |

## Notes

- In Bleve/Lucene syntax, the `+` prefix indicates a required term (AND logic)
- Terms without `+` or `-` are optional (OR logic by default)
- Phrase queries use double quotes to match exact sequences
- Special characters in terms are escaped automatically
