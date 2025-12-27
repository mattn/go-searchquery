# go-searchquery

A general-purpose Go library and command-line tool for parsing and matching search queries, similar to X (Twitter) search syntax.

## Features

- **Intuitive Search Syntax**: Works like familiar search engines (X/Twitter, Google, etc.)
- **Implicit AND**: Multiple terms are automatically combined with AND logic
- **Phrase Search**: Support for quoted phrases with exact matching
- **Case-Insensitive**: All matching is case-insensitive by default
- **Substring Matching**: Terms match anywhere within the content
- **Command-Line Tool**: Grep-style utility for filtering text
- **NIP-50 Compatible**: Implements the minimum required behavior for Nostr search queries

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
- **Nostr/NIP-50**: Implements the minimum required behavior for Nostr relay search
- **Any Application**: That needs X/Twitter-style search functionality

### NIP-50 Compliance

For Nostr applications, this library implements the **minimum required behavior (MUST)** from NIP-50:

- Implicit AND for multiple terms  
- Case-insensitive matching  
- Substring matching against content  
- Quoted phrase support  

**Note**: Explicit boolean operators (`AND`, `OR`) and parentheses are currently not supported due to parser implementation issues. These are marked as **SHOULD** (recommended) in NIP-50, not **MUST** (required).

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

## Author

Yasuhiro Matsumoto (a.k.a. mattn)

## License

MIT
