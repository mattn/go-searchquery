# MongoDB Dialect

Converts search queries to MongoDB `$text` search or `$regex` query format.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/mongodb
```

## Usage

### Text Search (requires text index)

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/mattn/go-searchquery/dialect/mongodb"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

func main() {
    userQuery := "hello world"
    textSearch, err := mongodb.ToTextSearch(userQuery)
    if err != nil {
        panic(err)
    }
    // textSearch = {"$text":{"$search":"\"hello\" \"world\""}}
    
    // Parse JSON and use in MongoDB query
    var filter bson.M
    json.Unmarshal([]byte(textSearch), &filter)
    
    // collection.Find(context.Background(), filter)
}
```

### Regex Query (no index required)

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/mattn/go-searchquery/dialect/mongodb"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    userQuery := "hello world"
    regexQuery, err := mongodb.ToRegexQuery(userQuery, "content")
    if err != nil {
        panic(err)
    }
    // regexQuery = {"$and":[{"content":{"$regex":"hello","$options":"i"}},{"content":{"$regex":"world","$options":"i"}}]}
    
    // Parse and use
    var filter bson.M
    json.Unmarshal([]byte(regexQuery), &filter)
}
```

## Query Conversion Examples

| Input Query | Text Search | Regex Query |
|-------------|-------------|-------------|
| `hello world` | `{"$text":{"$search":"\"hello\" \"world\""}}` | `{"$and":[{"content":{"$regex":"hello"}},{"content":{"$regex":"world"}}]}` |
| `"hello world"` | `{"$text":{"$search":"\"hello world\""}}` | `{"content":{"$regex":"hello world"}}` |

## Features

- **Two Query Types**: `$text` (fast, requires index) and `$regex` (flexible, slower)
- **JSON Output**: Ready to use with MongoDB drivers
- **Case Insensitive**: Regex queries include case-insensitive option
- **AND/OR Support**: Properly structured MongoDB queries

## Notes

- `$text` search requires a text index: `db.collection.createIndex({content: "text"})`
- `$text` is faster but requires index setup
- `$regex` queries are more flexible but slower (full collection scan)
- Both methods support phrase search and boolean operations
