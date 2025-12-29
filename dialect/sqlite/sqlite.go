package sqlite

import (
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToFTS5Query converts a search query to SQLite FTS5 MATCH format
// Example: "hello world" -> "hello AND world"
// Example: "\"hello world\"" -> "\"hello world\""
func ToFTS5Query(query string, opts ...searchquery.ParserOption) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query, opts...)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToFTS5Query(ast), nil
}

// nodeToFTS5Query recursively converts AST nodes to FTS5 format
func nodeToFTS5Query(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term - escape if needed
		return escapeFTS5Term(n.Phrase)
	case *searchquery.AndNode:
		left := nodeToFTS5Query(n.Left)
		right := nodeToFTS5Query(n.Right)
		return fmt.Sprintf("%s AND %s", left, right)
	case *searchquery.OrNode:
		left := nodeToFTS5Query(n.Left)
		right := nodeToFTS5Query(n.Right)
		return fmt.Sprintf("%s OR %s", left, right)
	default:
		return ""
	}
}

// escapeFTS5Term escapes special characters in FTS5 terms
func escapeFTS5Term(term string) string {
	// Remove quotes if present
	term = strings.Trim(term, "\"")

	// FTS5 special characters that need quoting: " *
	needsQuoting := strings.ContainsAny(term, "\"*")

	if needsQuoting {
		// Escape quotes by doubling them
		term = strings.ReplaceAll(term, "\"", "\"\"")
		return "\"" + term + "\""
	}
	return term
}
