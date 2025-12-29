package elasticsearch

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mattn/go-searchquery"
)

// ToQueryString converts a search query to Elasticsearch Query String format
// Example: "hello world" -> "hello AND world"
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

// nodeToQueryString recursively converts AST nodes to query string format
func nodeToQueryString(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase: keep quotes
			return fmt.Sprintf("\"%s\"", strings.Trim(n.Phrase, "\""))
		}
		// Single term - escape if needed
		return escapeQueryString(n.Phrase)
	case *searchquery.AndNode:
		left := nodeToQueryString(n.Left)
		right := nodeToQueryString(n.Right)
		return fmt.Sprintf("(%s AND %s)", left, right)
	case *searchquery.OrNode:
		left := nodeToQueryString(n.Left)
		right := nodeToQueryString(n.Right)
		return fmt.Sprintf("(%s OR %s)", left, right)
	default:
		return ""
	}
}

// escapeQueryString escapes special characters in query string
func escapeQueryString(term string) string {
	// Remove quotes if present
	term = strings.Trim(term, "\"")

	// Elasticsearch special characters: + - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /
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
		return "\"" + term + "\""
	}
	return term
}

// MatchQuery represents an Elasticsearch match query
type MatchQuery struct {
	Match map[string]interface{} `json:"match"`
}

// BoolQuery represents an Elasticsearch bool query
type BoolQuery struct {
	Bool BoolClause `json:"bool"`
}

// BoolClause contains must, should, etc.
type BoolClause struct {
	Must   []interface{} `json:"must,omitempty"`
	Should []interface{} `json:"should,omitempty"`
}

// ToMatchQuery converts a search query to Elasticsearch Match Query DSL
// Returns JSON string for the query
func ToMatchQuery(query, field string, opts ...searchquery.ParserOption) (string, error) {
	if query == "" {
		return "{}", nil
	}

	parser := searchquery.NewParser(query, opts...)
	ast, err := parser.Parse()
	if err != nil {
		return "", err
	}

	queryMap := nodeToMatchQuery(ast, field)
	jsonBytes, err := json.Marshal(queryMap)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// nodeToMatchQuery recursively converts AST nodes to match query DSL
func nodeToMatchQuery(node searchquery.Node, field string) interface{} {
	switch n := node.(type) {
	case *searchquery.TermNode:
		// Check if it's a phrase (contains spaces)
		if strings.Contains(n.Phrase, " ") {
			// Phrase query
			return map[string]interface{}{
				"match_phrase": map[string]interface{}{
					field: strings.Trim(n.Phrase, "\""),
				},
			}
		}
		// Single term match
		return map[string]interface{}{
			"match": map[string]interface{}{
				field: n.Phrase,
			},
		}
	case *searchquery.AndNode:
		left := nodeToMatchQuery(n.Left, field)
		right := nodeToMatchQuery(n.Right, field)
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{left, right},
			},
		}
	case *searchquery.OrNode:
		left := nodeToMatchQuery(n.Left, field)
		right := nodeToMatchQuery(n.Right, field)
		return map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []interface{}{left, right},
			},
		}
	default:
		return map[string]interface{}{}
	}
}
