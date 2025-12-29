package searchquery_test

import (
	"testing"

	"github.com/mattn/go-searchquery"
)

func TestMatchSimpleTerms(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		content string
		want    bool
	}{
		{
			name:    "single term match",
			query:   "hello",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "single term no match",
			query:   "goodbye",
			content: "Hello World",
			want:    false,
		},
		{
			name:    "case insensitive match",
			query:   "HELLO",
			content: "hello world",
			want:    true,
		},
		{
			name:    "partial word match",
			query:   "wor",
			content: "Hello World",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchImplicitAND(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		content string
		want    bool
	}{
		{
			name:    "two terms both present",
			query:   "hello world",
			content: "Hello there World",
			want:    true,
		},
		{
			name:    "two terms one missing",
			query:   "hello mars",
			content: "Hello there World",
			want:    false,
		},
		{
			name:    "multiple terms all present",
			query:   "nostr apps",
			content: "Best Nostr Apps 2025",
			want:    true,
		},
		{
			name:    "three terms all present",
			query:   "cat dog bird",
			content: "I have a cat, a dog, and a bird",
			want:    true,
		},
		{
			name:    "three terms one missing",
			query:   "cat dog fish",
			content: "I have a cat and a dog",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchPhrases(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		content string
		want    bool
	}{
		{
			name:    "exact phrase match",
			query:   `"hello world"`,
			content: "Say hello world today",
			want:    true,
		},
		{
			name:    "phrase not contiguous",
			query:   `"hello world"`,
			content: "world hello",
			want:    false,
		},
		{
			name:    "phrase with case insensitive",
			query:   `"Hello World"`,
			content: "say hello world today",
			want:    true,
		},
		{
			name:    "phrase in middle of content",
			query:   `"quick brown"`,
			content: "The quick brown fox jumps",
			want:    true,
		},
		{
			name:    "phrase words present but not together",
			query:   `"brown fox"`,
			content: "The brown and red fox",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchExplicitANDOperator(t *testing.T) {
	// NIP-50 SHOULD support explicit AND operator
	tests := []struct {
		name       string
		query      string
		content    string
		want       bool
		skip       bool
		skipReason string
	}{
		{
			name:    "explicit AND both present",
			query:   "hello AND world",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "explicit AND first missing",
			query:   "goodbye AND world",
			content: "Hello World",
			want:    false,
		},
		{
			name:    "explicit AND second missing",
			query:   "hello AND mars",
			content: "Hello World",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchExplicitOROperator(t *testing.T) {
	// NIP-50 SHOULD support explicit OR operator
	tests := []struct {
		name       string
		query      string
		content    string
		want       bool
		skip       bool
		skipReason string
	}{
		{
			name:    "explicit OR both present",
			query:   "hello OR world",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "explicit OR first present",
			query:   "hello OR mars",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "explicit OR second present",
			query:   "goodbye OR world",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "explicit OR both missing",
			query:   "goodbye OR mars",
			content: "Hello World",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchImplicitANDOperator(t *testing.T) {
	// NIP-50 MUST support implicit AND (already tested in TestMatchImplicitAND)
	// This test focuses on the operator behavior specifically
	tests := []struct {
		name    string
		query   string
		content string
		want    bool
	}{
		{
			name:    "implicit AND match",
			query:   "cat dog",
			content: "I have a cat and a dog",
			want:    true,
		},
		{
			name:    "implicit AND no match",
			query:   "cat fish",
			content: "I have a cat and a dog",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchParentheses(t *testing.T) {
	// NIP-50 SHOULD support parentheses for grouping
	tests := []struct {
		name       string
		query      string
		content    string
		want       bool
		skip       bool
		skipReason string
	}{
		{
			name:    "simple grouping with implicit AND",
			query:   "(cat dog)",
			content: "I have a cat and a dog",
			want:    true,
		},
		{
			name:    "AND with OR in parentheses - match",
			query:   "cat AND (dog OR bird)",
			content: "I have a cat and a bird",
			want:    true,
		},
		{
			name:    "AND with OR in parentheses - no match",
			query:   "cat AND (dog OR bird)",
			content: "I have a cat",
			want:    false,
		},
		{
			name:    "OR with AND in parentheses",
			query:   "(cat AND dog) OR bird",
			content: "I have a bird",
			want:    true,
		},
		{
			name:    "precedence: AND higher than OR",
			query:   "cat AND dog OR bird AND fish",
			content: "I have a bird and a fish",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		content string
		want    bool
	}{
		{
			name:    "empty query",
			query:   "",
			content: "Hello World",
			want:    true,
		},
		{
			name:    "empty content",
			query:   "hello",
			content: "",
			want:    false,
		},
		{
			name:    "both empty",
			query:   "",
			content: "",
			want:    true,
		},
		{
			name:    "special characters in term",
			query:   "hello@world",
			content: "Contact hello@world.com",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}

func TestMatchErrors(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		content   string
		wantError bool
	}{
		{
			name:      "mismatched opening parenthesis",
			query:     "cat AND (dog OR bird",
			content:   "I have a cat and a dog",
			wantError: true,
		},
		{
			name:      "mismatched closing parenthesis",
			query:     "cat AND dog OR bird)",
			content:   "I have a cat and a dog",
			wantError: true,
		},
		{
			name:      "multiple mismatched parentheses",
			query:     "((cat AND dog)",
			content:   "I have a cat and a dog",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := searchquery.Match(tt.content, tt.query)
			if tt.wantError && err == nil {
				t.Errorf("Match(%q, %q) expected error but got nil", tt.content, tt.query)
			}
			if !tt.wantError && err != nil {
				t.Errorf("Match(%q, %q) unexpected error: %v", tt.content, tt.query, err)
			}
		})
	}
}

func TestMatchComplexQueries(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		content    string
		want       bool
		skip       bool
		skipReason string
	}{
		{
			name:    "phrase with implicit AND",
			query:   `"hello world" test`,
			content: "This is a hello world test",
			want:    true,
		},
		{
			name:    "multiple phrases with implicit AND",
			query:   `"hello world" "test case"`,
			content: "This hello world is a test case",
			want:    true,
		},
		{
			name:    "phrase with explicit AND",
			query:   `"hello world" AND test`,
			content: "This is a hello world test",
			want:    true,
		},
		{
			name:    "phrase with explicit OR",
			query:   `"hello world" OR "goodbye world"`,
			content: "Say goodbye world",
			want:    true,
		},
		{
			name:    "NIP-50 example: nostr apps",
			query:   "nostr apps",
			content: "Best Nostr Apps 2025",
			want:    true,
		},
		{
			name:    "real-world query with operators",
			query:   `(golang OR go) AND (tutorial OR guide)`,
			content: "A beginner's guide to golang programming",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skip {
				t.Skip(tt.skipReason)
			}
			got, err := searchquery.Match(tt.content, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("Match(%q, %q) = %v, want %v", tt.content, tt.query, got, tt.want)
			}
		})
	}
}
