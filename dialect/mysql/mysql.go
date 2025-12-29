package mysql

import (
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToBoolean converts a search query to MySQL MATCH AGAINST boolean mode format
// Example: "hello world" -> "+hello +world"
// Example: "\"hello world\"" -> "\"hello world\""
func ToBoolean(query string, opts ...searchquery.ParserOption) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query, opts...)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToBoolean(ast), nil
}

// nodeToBoolean recursively converts AST nodes to MySQL boolean mode format
func nodeToBoolean(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term - prefix with + for required
		return "+" + escapeTerm(n.Phrase)
	case *searchquery.AndNode:
		// In MySQL boolean mode, + prefix means required (AND)
		left := nodeToBoolean(n.Left)
		right := nodeToBoolean(n.Right)
		return left + " " + right
	case *searchquery.OrNode:
		// In MySQL boolean mode, no prefix or parentheses for OR
		left := nodeToBoolean(n.Left)
		right := nodeToBoolean(n.Right)
		// Remove + prefix for OR terms
		left = strings.TrimPrefix(left, "+")
		right = strings.TrimPrefix(right, "+")
		return left + " " + right
	default:
		return ""
	}
}

// escapeTerm escapes special characters in MySQL boolean mode terms
func escapeTerm(term string) string {
	// Remove quotes if present
	term = strings.Trim(term, "\"")

	// MySQL boolean mode special characters: + - > < ( ) ~ * " @
	// For simple terms, we just need to handle quotes
	if strings.ContainsAny(term, "+-><()~*\"@") {
		// Quote the term if it contains special characters
		term = strings.ReplaceAll(term, "\"", "\\\"")
		return "\"" + term + "\""
	}
	return term
}
