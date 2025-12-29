package searchquery_test

import (
	"strings"
	"testing"

	"github.com/mattn/go-searchquery"
)

func TestWithLookup(t *testing.T) {
	tests := []struct {
		name   string
		query  string
		lookup searchquery.LookupFunc
		want   string
	}{
		{
			name:  "uppercase transformation",
			query: "hello world",
			lookup: func(token string) string {
				return strings.ToUpper(token)
			},
			want: "(HELLO AND WORLD)",
		},
		{
			name:  "filter out specific token",
			query: "good bad ugly",
			lookup: func(token string) string {
				if token == "bad" {
					return ""
				}
				return token
			},
			want: "(good AND ugly)",
		},
		{
			name:  "alias replacement",
			query: "x OR fb",
			lookup: func(token string) string {
				aliases := map[string]string{
					"x":  "twitter",
					"fb": "facebook",
				}
				if alias, ok := aliases[token]; ok {
					return alias
				}
				return token
			},
			want: "(twitter OR facebook)",
		},
		{
			name:  "stopword removal",
			query: "the cat and dog",
			lookup: func(token string) string {
				stopwords := map[string]bool{"the": true, "a": true, "an": true}
				if stopwords[token] {
					return ""
				}
				return token
			},
			want: "((cat AND and) AND dog)",
		},
		{
			name:  "phrase with lookup",
			query: `"hello world"`,
			lookup: func(token string) string {
				return strings.ToUpper(token)
			},
			want: "HELLO WORLD",
		},
		{
			name:  "filter all tokens",
			query: "bad bad bad",
			lookup: func(token string) string {
				return "" // Remove all
			},
			want: "", // Should result in error, but we'll handle gracefully
		},
		{
			name:   "no lookup function",
			query:  "hello world",
			lookup: nil,
			want:   "(hello AND world)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var parser *searchquery.Parser
			if tt.lookup != nil {
				parser = searchquery.NewParser(tt.query, searchquery.WithLookup(tt.lookup))
			} else {
				parser = searchquery.NewParser(tt.query)
			}

			ast, err := parser.Parse()
			if tt.want == "" {
				// Expect error for empty result
				if err == nil {
					t.Errorf("Expected error for query %q, but got none", tt.query)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.query, err)
			}

			got := nodeToString(ast)
			if got != tt.want {
				t.Errorf("Parse(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestLookupWithMatch(t *testing.T) {
	// Test that lookup works with Match function
	content := "I love TWITTER and FACEBOOK"

	// Define aliases
	aliases := map[string]string{
		"x":  "twitter",
		"fb": "facebook",
	}

	// Create parser with lookup
	parser := searchquery.NewParser("x OR fb", searchquery.WithLookup(func(token string) string {
		if alias, ok := aliases[token]; ok {
			return alias
		}
		return token
	}))

	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	// Manually evaluate the AST
	result := evalNode(ast, strings.ToLower(content))
	if !result {
		t.Errorf("Expected match for 'x OR fb' against %q", content)
	}
}

// Helper functions for testing
func nodeToString(node searchquery.Node) string {
	switch n := node.(type) {
	case *searchquery.TermNode:
		return n.Phrase
	case *searchquery.AndNode:
		return "(" + nodeToString(n.Left) + " AND " + nodeToString(n.Right) + ")"
	case *searchquery.OrNode:
		return "(" + nodeToString(n.Left) + " OR " + nodeToString(n.Right) + ")"
	default:
		return ""
	}
}

func evalNode(node searchquery.Node, content string) bool {
	switch n := node.(type) {
	case *searchquery.TermNode:
		return strings.Contains(content, strings.ToLower(n.Phrase))
	case *searchquery.AndNode:
		return evalNode(n.Left, content) && evalNode(n.Right, content)
	case *searchquery.OrNode:
		return evalNode(n.Left, content) || evalNode(n.Right, content)
	default:
		return false
	}
}
