# MySQL Dialect

Converts search queries to MySQL MATCH AGAINST boolean mode format.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/mysql
```

## Usage

```go
package main

import (
    "database/sql"
    "fmt"
    "github.com/mattn/go-searchquery/dialect/mysql"
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
    mysqlQuery, err := mysql.ToBoolean(userQuery)
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
    mysqlQuery, _ = mysql.ToBoolean(phraseQuery)
    // mysqlQuery = "\"hello world\""
}
```

## Query Conversion Examples

| Input Query | MySQL Boolean Mode |
|-------------|-------------------|
| `hello world` | `+hello +world` |
| `"hello world"` | `"hello world"` |
| `cat dog bird` | `+cat +dog +bird` |

## Features

- **Required Terms**: The `+` prefix marks terms as required (AND logic)
- **Phrase Search**: Quoted phrases are preserved for exact matching
- **OR Support**: Terms without `+` prefix (OR behavior)
- **Special Character Handling**: Properly escapes MySQL boolean mode special characters

## Notes

- Requires a FULLTEXT index on the searched column(s)
- Use `IN BOOLEAN MODE` in your MATCH AGAINST queries
- Available in MyISAM tables (all versions) and InnoDB tables (MySQL 5.6+)
