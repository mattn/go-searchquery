# PostgreSQL Dialect

Converts search queries to PostgreSQL `tsquery` format for full-text search.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/postgres
```

## Usage

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery/dialect/postgres"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "...")
    
    // Convert user query to tsquery
    userQuery := "hello world"
    tsquery, err := postgres.ToTsQuery(userQuery)
    if err != nil {
        panic(err)
    }
    // tsquery = "(hello & world)"
    
    // Use in PostgreSQL query
    rows, err := db.Query(`
        SELECT title, content 
        FROM articles 
        WHERE to_tsvector('english', content) @@ to_tsquery('english', $1)
    `, tsquery)
    
    // Phrase search example
    phraseQuery := `"hello world"`
    tsquery, _ = postgres.ToTsQuery(phraseQuery)
    // tsquery = "hello <-> world"  (followed-by operator)
}
```

## Query Conversion Examples

| Input Query | PostgreSQL tsquery |
|-------------|-------------------|
| `hello world` | `(hello & world)` |
| `"hello world"` | `hello <-> world` |
| `cat dog bird` | `((cat & dog) & bird)` |
| `"quick brown fox"` | `quick <-> brown <-> fox` |

## Features

- **Implicit AND**: Multiple terms are combined with `&` operator
- **Phrase Search**: Quoted phrases use `<->` (followed-by) operator
- **Special Character Escaping**: Handles special characters properly
- **Case Handling**: Terms are converted according to PostgreSQL tsquery rules

## Notes

- Requires a full-text search index on your PostgreSQL table
- Use `to_tsvector()` and `to_tsquery()` functions in your SQL queries
- Consider using the appropriate language configuration (e.g., 'english', 'japanese')
