package meilisearch

import (
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToQuery converts a search query to Meilisearch query format
// Meilisearch uses a simple query syntax similar to standard search engines
// Example: "hello world" -> "hello world"
// Example: "\"hello world\"" -> "\"hello world\""
func ToQuery(query string) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToQuery(ast), nil
}

// nodeToQuery recursively converts AST nodes to Meilisearch query format
func nodeToQuery(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term
		return n.Phrase
	case *searchquery.AndNode:
		// Meilisearch treats space-separated terms as AND by default
		left := nodeToQuery(n.Left)
		right := nodeToQuery(n.Right)
		return left + " " + right
	case *searchquery.OrNode:
		// Meilisearch supports OR operator
		left := nodeToQuery(n.Left)
		right := nodeToQuery(n.Right)
		return fmt.Sprintf("(%s OR %s)", left, right)
	default:
		return ""
	}
}

// ToFilter converts a search query to Meilisearch filter format
// This is useful for exact matching on specific fields
// Example: field="title", query="hello world" -> "title = \"hello world\""
func ToFilter(field, query string) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToFilter(ast, field), nil
}

// nodeToFilter recursively converts AST nodes to Meilisearch filter format
func nodeToFilter(node searchquery.Node, field string) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Simple equality filter
		phrase := strings.Trim(n.Phrase, "\"")
		return fmt.Sprintf("%s = \"%s\"", field, escapeFilterValue(phrase))
	case *searchquery.AndNode:
		left := nodeToFilter(n.Left, field)
		right := nodeToFilter(n.Right, field)
		return fmt.Sprintf("(%s AND %s)", left, right)
	case *searchquery.OrNode:
		left := nodeToFilter(n.Left, field)
		right := nodeToFilter(n.Right, field)
		return fmt.Sprintf("(%s OR %s)", left, right)
	default:
		return ""
	}
}

// escapeFilterValue escapes special characters in filter values
func escapeFilterValue(value string) string {
	// Escape quotes by doubling them
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return value
}
