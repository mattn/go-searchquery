# Elasticsearch Dialect

Converts search queries to Elasticsearch Query String or Match Query DSL format.

## Installation

```bash
go get github.com/mattn/go-searchquery/dialect/elasticsearch
```

## Usage

### Query String Format

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery/dialect/elasticsearch"
)

func main() {
    userQuery := "hello world"
    queryString, err := elasticsearch.ToQueryString(userQuery)
    if err != nil {
        panic(err)
    }
    // queryString = "(hello AND world)"
    
    // Use in Elasticsearch query_string query
    // GET /articles/_search
    // {
    //   "query": {
    //     "query_string": {
    //       "query": "(hello AND world)",
    //       "default_field": "content"
    //     }
    //   }
    // }
}
```

### Match Query DSL Format

```go
package main

import (
    "fmt"
    "github.com/mattn/go-searchquery/dialect/elasticsearch"
)

func main() {
    userQuery := "hello world"
    matchQuery, err := elasticsearch.ToMatchQuery(userQuery, "content")
    if err != nil {
        panic(err)
    }
    // matchQuery = {"bool":{"must":[{"match":{"content":"hello"}},{"match":{"content":"world"}}]}}
    
    // Phrase search
    phraseQuery := `"hello world"`
    matchQuery, _ = elasticsearch.ToMatchQuery(phraseQuery, "content")
    // matchQuery = {"match_phrase":{"content":"hello world"}}
}
```

## Query Conversion Examples

| Input Query | Query String | Match Query DSL |
|-------------|--------------|-----------------|
| `hello world` | `(hello AND world)` | `{"bool":{"must":[{"match":{"content":"hello"}},{"match":{"content":"world"}}]}}` |
| `"hello world"` | `"hello world"` | `{"match_phrase":{"content":"hello world"}}` |
| `cat dog bird` | `((cat AND dog) AND bird)` | Nested bool query |

## Features

- **Two Output Formats**: Query String (simple) or Match Query DSL (powerful)
- **Phrase Search**: Converted to `match_phrase` queries
- **Bool Queries**: AND/OR operations use bool query structure
- **Special Character Escaping**: Handles Elasticsearch special characters

## Notes

- Query String format is simpler but less flexible
- Match Query DSL provides more control over scoring and analysis
- Works with both Elasticsearch and OpenSearch
- Consider using appropriate analyzers for your use case
