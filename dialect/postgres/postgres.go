package postgres

import (
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToTsQuery converts a search query to PostgreSQL tsquery format
// Example: "hello world" -> "(hello & world)"
// Example: "\"hello world\"" -> "hello <-> world"
func ToTsQuery(query string, opts ...searchquery.ParserOption) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := searchquery.NewParser(query, opts...)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToTsQuery(ast), nil
}

// nodeToTsQuery recursively converts AST nodes to tsquery format
func nodeToTsQuery(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: convert to <-> (followed by) operator
			words := strings.Fields(n.Phrase)
			escapedWords := make([]string, len(words))
			for i, word := range words {
				escapedWords[i] = escapeTsQueryTerm(word)
			}
			return strings.Join(escapedWords, " <-> ")
		}
		// Single term
		return escapeTsQueryTerm(n.Phrase)
	case *searchquery.AndNode:
		left := nodeToTsQuery(n.Left)
		right := nodeToTsQuery(n.Right)
		return fmt.Sprintf("(%s & %s)", left, right)
	case *searchquery.OrNode:
		left := nodeToTsQuery(n.Left)
		right := nodeToTsQuery(n.Right)
		return fmt.Sprintf("(%s | %s)", left, right)
	default:
		return ""
	}
}

// escapeTsQueryTerm escapes special characters in tsquery terms
func escapeTsQueryTerm(term string) string {
	// Remove quotes if present
	term = strings.Trim(term, "\"")

	// Escape single quotes by doubling them
	term = strings.ReplaceAll(term, "'", "''")

	// Check if term needs quoting (contains special chars or spaces)
	needsQuoting := false
	for _, r := range term {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			needsQuoting = true
			break
		}
	}

	if needsQuoting {
		return "'" + term + "'"
	}
	return term
}
