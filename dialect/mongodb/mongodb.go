package mongodb

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToTextSearch converts a search query to MongoDB $text search format
// Example: "hello world" -> {"$text": {"$search": "hello world"}}
// Note: MongoDB $text with multiple terms uses implicit OR by default
// For AND behavior, use quotes: "\"hello\" \"world\""
func ToTextSearch(query string) (string, error) {
	if query == "" {
		return "{}", nil
	}

	parser := searchquery.NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	searchString := nodeToTextSearch(ast)
	result := map[string]interface{}{
		"$text": map[string]interface{}{
			"$search": searchString,
		},
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// nodeToTextSearch recursively converts AST nodes to MongoDB text search string
func nodeToTextSearch(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes for exact match
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term - wrap in quotes for required term (AND behavior)
		return fmt.Sprintf("\"%s\"", n.Phrase)
	case *searchquery.AndNode:
		// MongoDB: quoted terms are required (AND)
		left := nodeToTextSearch(n.Left)
		right := nodeToTextSearch(n.Right)
		return left + " " + right
	case *searchquery.OrNode:
		// MongoDB: unquoted terms for OR
		left := removeQuotes(nodeToTextSearch(n.Left))
		right := removeQuotes(nodeToTextSearch(n.Right))
		return left + " " + right
	default:
		return ""
	}
}

// removeQuotes removes surrounding quotes from a string
func removeQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// ToRegexQuery converts a search query to MongoDB $regex query
// This provides more control but is slower than $text search
// Example: "hello world" -> {"$and": [{"field": {"$regex": "hello", "$options": "i"}}, {"field": {"$regex": "world", "$options": "i"}}]}
func ToRegexQuery(query, field string) (string, error) {
	if query == "" {
		return "{}", nil
	}

	parser := searchquery.NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	queryMap := nodeToRegexQuery(ast, field)
	jsonBytes, err := json.Marshal(queryMap)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// nodeToRegexQuery recursively converts AST nodes to MongoDB regex query
func nodeToRegexQuery(node searchquery.Node, field string) interface{} {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: exact match with spaces
			phrase := strings.Trim(n.Phrase, "\"")
			// Escape regex special characters
			phrase = escapeRegex(phrase)
			return map[string]interface{}{
				field: map[string]interface{}{
					"$regex":   phrase,
					"$options": "i", // case insensitive
				},
			}
		}
		// Single term
		term := escapeRegex(n.Phrase)
		return map[string]interface{}{
			field: map[string]interface{}{
				"$regex":   term,
				"$options": "i",
			},
		}
	case *searchquery.AndNode:
		left := nodeToRegexQuery(n.Left, field)
		right := nodeToRegexQuery(n.Right, field)
		return map[string]interface{}{
			"$and": []interface{}{left, right},
		}
	case *searchquery.OrNode:
		left := nodeToRegexQuery(n.Left, field)
		right := nodeToRegexQuery(n.Right, field)
		return map[string]interface{}{
			"$or": []interface{}{left, right},
		}
	default:
		return map[string]interface{}{}
	}
}

// escapeRegex escapes regex special characters
func escapeRegex(s string) string {
	s = strings.Trim(s, "\"")
	// Escape regex special chars: . * + ? ^ $ { } ( ) | [ ] \
	specialChars := []string{"\\", ".", "*", "+", "?", "^", "$", "{", "}", "(", ")", "|", "[", "]"}
	for _, char := range specialChars {
		s = strings.ReplaceAll(s, char, "\\"+char)
	}
	return s
}
