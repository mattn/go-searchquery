# SQLite Dialect

Converts search queries to SQLite FTS5 MATCH format for full-text search.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/sqlite
```

## Usage

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery/dialect/sqlite"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "test.db")
    
    // Create FTS5 table
    db.Exec(`CREATE VIRTUAL TABLE articles USING fts5(title, content)`)
    
    // Convert user query to FTS5 format
    userQuery := "hello world"
    ftsQuery, err := sqlite.ToFTS5Query(userQuery)
    if err != nil {
        panic(err)
    }
    // ftsQuery = "hello AND world"
    
    // Use in SQLite FTS5 query
    rows, err := db.Query(`
        SELECT title, content 
        FROM articles 
        WHERE articles MATCH ?
    `, ftsQuery)
    
    // Phrase search example
    phraseQuery := `"hello world"`
    ftsQuery, _ = sqlite.ToFTS5Query(phraseQuery)
    // ftsQuery = "\"hello world\""  (quoted phrase)
}
```

## Query Conversion Examples

| Input Query | SQLite FTS5 Query |
|-------------|-------------------|
| `hello world` | `hello AND world` |
| `"hello world"` | `"hello world"` |
| `cat dog bird` | `cat AND dog AND bird` |

## Features

- **Implicit AND**: Multiple terms use `AND` keyword
- **Phrase Search**: Quoted phrases remain quoted
- **OR Support**: Explicit OR operators are preserved
- **Simple Syntax**: Clean, readable query format

## Notes

- Requires SQLite FTS5 extension
- Use `CREATE VIRTUAL TABLE ... USING fts5(...)` to create searchable tables
- FTS5 is available in SQLite 3.9.0 and later
