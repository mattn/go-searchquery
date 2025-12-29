package bleve

import (
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToQueryString converts a search query to Bleve query string format
// Bleve uses a query string syntax similar to Lucene
// Example: "hello world" -> "+hello +world"
// Example: "hello OR world" -> "hello world"
// Example: "\"hello world\"" -> "\"hello world\""
func ToQueryString(query string, opts ...searchquery.ParserOption) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query, opts...)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToQueryString(ast), nil
}

// nodeToQueryString recursively converts AST nodes to Bleve query string format
func nodeToQueryString(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term
		return escapeTerm(n.Phrase)
	case *searchquery.AndNode:
		left := nodeToQueryString(n.Left)
		right := nodeToQueryString(n.Right)
		// Use + prefix for required terms
		leftTerm := addPlusIfNeeded(left)
		rightTerm := addPlusIfNeeded(right)
		return fmt.Sprintf("%s %s", leftTerm, rightTerm)
	case *searchquery.OrNode:
		left := nodeToQueryString(n.Left)
		right := nodeToQueryString(n.Right)
		// For simple OR, don't wrap in parentheses
		return fmt.Sprintf("%s %s", left, right)
	default:
		return ""
	}
}

// addPlusIfNeeded adds + prefix if not already present
func addPlusIfNeeded(term string) string {
	term = strings.TrimSpace(term)
	if strings.HasPrefix(term, "+") {
		return term
	}
	// For OR groups or phrases, wrap with +
	if strings.Contains(term, " ") && !strings.HasPrefix(term, "\"") {
		return "+(" + term + ")"
	}
	if strings.HasPrefix(term, "\"") {
		return "+" + term
	}
	return "+" + term
}

// escapeTerm escapes special characters in a term
func escapeTerm(term string) string {
	// Remove quotes if present
	term = strings.Trim(term, "\"")
	
	// Bleve/Lucene special characters: + - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /
	specialChars := []string{"+", "-", "=", "&&", "||", ">", "<", "!", "(", ")", "{", "}", "[", "]", "^", "\"", "~", "*", "?", ":", "\\", "/"}
	needsEscape := false
	for _, char := range specialChars {
		if strings.Contains(term, char) {
			needsEscape = true
			break
		}
	}

	if needsEscape {
		// Escape backslashes first
		term = strings.ReplaceAll(term, "\\", "\\\\")
		// Escape quotes
		term = strings.ReplaceAll(term, "\"", "\\\"")
	}
	return term
}
