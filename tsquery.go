package searchquery

import (
	"fmt"
	"strings"
)

// ToTsQuery converts a search query to PostgreSQL tsquery format
// Example: "hello world" -> "hello & world"
// Example: "\"hello world\"" -> "hello <-> world"
func ToTsQuery(query string) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToTsQuery(ast), nil
}

// nodeToTsQuery recursively converts AST nodes to tsquery format
func nodeToTsQuery(node Node) string {
	switch n := node.(type) {
	case TermNode:
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
	case AndNode:
		left := nodeToTsQuery(n.Left)
		right := nodeToTsQuery(n.Right)
		return fmt.Sprintf("(%s & %s)", left, right)
	case OrNode:
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

// ToFTS5Query converts a search query to SQLite FTS5 MATCH format
// Example: "hello world" -> "hello AND world"
// Example: "\"hello world\"" -> "\"hello world\""
func ToFTS5Query(query string) (string, error) {
	if query == "" {
		return "", nil
	}

	parser := NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	return nodeToFTS5Query(ast), nil
}

// nodeToFTS5Query recursively converts AST nodes to FTS5 format
func nodeToFTS5Query(node Node) string {
	switch n := node.(type) {
	case TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term - escape if needed
		return escapeFTS5Term(n.Phrase)
	case AndNode:
		left := nodeToFTS5Query(n.Left)
		right := nodeToFTS5Query(n.Right)
		return fmt.Sprintf("%s AND %s", left, right)
	case OrNode:
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
