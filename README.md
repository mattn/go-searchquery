# go-searchquery

A general-purpose Go library and command-line tool for parsing and matching search queries, similar to X (Twitter) search syntax.

## Features

- **Intuitive Search Syntax**: Works like familiar search engines (X/Twitter, Google, etc.)
- **Implicit AND**: Multiple terms are automatically combined with AND logic
- **Phrase Search**: Support for quoted phrases with exact matching
- **Case-Insensitive**: All matching is case-insensitive by default
- **Substring Matching**: Terms match anywhere within the content
- **Command-Line Tool**: Grep-style utility for filtering text

## Installation

```bash
go get github.com/mattn/go-searchquery
```

To install the command-line tool:

```bash
go install github.com/mattn/go-searchquery/cmd/query@latest
```

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery"
)

func main() {
    // Simple term search
    match, err := searchquery.Match("Hello World", "hello")
    fmt.Println(match) // true

    // Implicit AND (multiple terms)
    match, err = searchquery.Match("Best Nostr Apps 2025", "nostr apps")
    fmt.Println(match) // true

    // Phrase search
    match, err = searchquery.Match("Say hello world today", `"hello world"`)
    fmt.Println(match) // true

    // Phrase must be contiguous
    match, err = searchquery.Match("world hello", `"hello world"`)
    fmt.Println(match) // false
}
```

## Database Integration

### PostgreSQL Full-Text Search

Convert search queries to PostgreSQL `tsquery` format for full-text search:

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery"
    _ "github.com/lib/pq"
)

func main() {
    db, _ := sql.Open("postgres", "...")
    
    // Convert user query to tsquery
    userQuery := "hello world"
    tsquery, err := searchquery.ToTsQuery(userQuery)
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
    tsquery, _ = searchquery.ToTsQuery(phraseQuery)
    // tsquery = "hello <-> world"  (followed-by operator)
}
```

**Query Conversion Examples:**

| Input Query | PostgreSQL tsquery |
|-------------|-------------------|
| `hello world` | `(hello & world)` |
| `"hello world"` | `hello <-> world` |
| `cat dog bird` | `((cat & dog) & bird)` |
| `"quick brown fox"` | `quick <-> brown <-> fox` |

### SQLite3 Full-Text Search (FTS5)

For SQLite3, you can use the FTS5 extension with MATCH queries:

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, _ := sql.Open("sqlite3", "test.db")
    
    // Create FTS5 table
    db.Exec(`CREATE VIRTUAL TABLE articles USING fts5(title, content)`)
    
    // Convert user query to FTS5 format
    userQuery := "hello world"
    ftsQuery, err := searchquery.ToFTS5Query(userQuery)
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
    ftsQuery, _ = searchquery.ToFTS5Query(phraseQuery)
    // ftsQuery = "\"hello world\""  (quoted phrase)
}
```

**Query Conversion Examples:**

| Input Query | SQLite FTS5 Query |
|-------------|-------------------|
| `hello world` | `hello AND world` |
| `"hello world"` | `"hello world"` |
| `cat dog bird` | `cat AND dog AND bird` |

**Note**: Both PostgreSQL and SQLite3 examples assume you have appropriate full-text search indexes set up on your tables.

### MySQL Full-Text Search

For MySQL, use the MATCH AGAINST syntax with boolean mode:

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, _ := sql.Open("mysql", "user:password@/dbname")
    
    // Create table with FULLTEXT index
    db.Exec(`CREATE TABLE articles (
        id INT PRIMARY KEY,
        title VARCHAR(200),
        content TEXT,
        FULLTEXT(content)
    )`)
    
    // Convert user query to MySQL boolean mode
    userQuery := "hello world"
    mysqlQuery, err := searchquery.ToMySQLBoolean(userQuery)
    if err != nil {
        panic(err)
    }
    // mysqlQuery = "+hello +world"
    
    // Use in MySQL MATCH AGAINST query
    rows, err := db.Query(`
        SELECT title, content 
        FROM articles 
        WHERE MATCH(content) AGAINST(? IN BOOLEAN MODE)
    `, mysqlQuery)
    
    // Phrase search example
    phraseQuery := `"hello world"`
    mysqlQuery, _ = searchquery.ToMySQLBoolean(phraseQuery)
    // mysqlQuery = "\"hello world\""
}
```

**Query Conversion Examples:**

| Input Query | MySQL Boolean Mode |
|-------------|-------------------|
| `hello world` | `+hello +world` |
| `"hello world"` | `"hello world"` |
| `cat dog bird` | `+cat +dog +bird` |

**Note**: The `+` prefix in MySQL boolean mode means the term is required (AND logic). Phrases are kept in quotes for exact matching.


## Command-Line Tool

The `query` command works like grep but uses NIP-50 search query syntax.

### Basic Usage

```bash
# Search for lines containing both "hello" and "world"
query "hello world" file.txt

# Search for exact phrase
query '"hello world"' file.txt

# Read from stdin
cat file.txt | query "search term"

# Search multiple files
query "pattern" file1.txt file2.txt
```

### Options

- `-n`: Print line numbers with output lines
- `-v`: Invert match (select non-matching lines)
- `-c`: Only print count of matching lines
- `-q`: Quiet mode (suppress normal output, exit code indicates match)

### Examples

```bash
# Show line numbers
echo -e "Hello World\nGoodbye Mars\nHello there World" | query -n "hello world"
# Output:
# 1:Hello World
# 3:Hello there World

# Count matching lines
query -c "error" logfile.txt

# Invert match (show lines that DON'T match)
query -v "debug" logfile.txt

# Search for phrase in multiple files
query '"error occurred"' *.log
```

## Query Syntax

This library implements a search syntax similar to popular search engines like X (Twitter):

### Core Features

- **Implicit AND**: `hello world` matches content containing both "hello" AND "world"
- **Case-Insensitive**: `HELLO` matches "hello", "Hello", "HELLO", etc.
- **Substring Matching**: `wor` matches "world", "work", "sword", etc.
- **Phrase Search**: `"hello world"` matches the exact phrase (contiguous)

### Examples

```go
// Implicit AND - both terms must be present
Match("I have a cat and a dog", "cat dog")  // true
Match("I have a cat", "cat dog")            // false

// Phrase search - terms must be contiguous
Match("Say hello world today", `"hello world"`)  // true
Match("world hello", `"hello world"`)            // false

// Case insensitive
Match("Hello World", "HELLO")  // true

// Substring matching
Match("Hello World", "wor")    // true
```

## Use Cases

This library is suitable for:

- **Text Search Applications**: Add search functionality to your applications
- **Log Filtering**: Search through log files with intuitive query syntax
- **Content Management**: Filter and search user-generated content
- **Any Application**: That needs X/Twitter-style search functionality

## Testing

```bash
go test -v
```

The test suite includes comprehensive coverage of:
- Simple term matching
- Implicit AND behavior
- Phrase search
- Edge cases
- Error handling

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
