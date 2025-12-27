package searchquery_test

import (
	"testing"

	"github.com/mattn/go-searchquery"
)

func TestToTsQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "single term",
			query: "hello",
			want:  "hello",
		},
		{
			name:  "implicit AND",
			query: "hello world",
			want:  "(hello & world)",
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  "hello <-> world",
		},
		{
			name:  "three terms",
			query: "cat dog bird",
			want:  "((cat & dog) & bird)",
		},
		{
			name:  "phrase with multiple words",
			query: `"quick brown fox"`,
			want:  "quick <-> brown <-> fox",
		},
		{
			name:  "term with special characters",
			query: "hello@world",
			want:  "'hello@world'",
		},
		{
			name:  "empty query",
			query: "",
			want:  "",
		},
		{
			name:  "multiple phrases",
			query: `"hello world" "test case"`,
			want:  "(hello <-> world & test <-> case)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.ToTsQuery(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToTsQuery(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestToTsQueryErrors(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "mismatched parentheses",
			query: "hello (world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := searchquery.ToTsQuery(tt.query)
			if err == nil {
				t.Errorf("ToTsQuery(%q) expected error but got nil", tt.query)
			}
		})
	}
}

func TestToFTS5Query(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "single term",
			query: "hello",
			want:  "hello",
		},
		{
			name:  "implicit AND",
			query: "hello world",
			want:  "hello AND world",
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "three terms",
			query: "cat dog bird",
			want:  "cat AND dog AND bird",
		},
		{
			name:  "phrase with multiple words",
			query: `"quick brown fox"`,
			want:  `"quick brown fox"`,
		},
		{
			name:  "empty query",
			query: "",
			want:  "",
		},
		{
			name:  "multiple phrases",
			query: `"hello world" "test case"`,
			want:  `"hello world" AND "test case"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.ToFTS5Query(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToFTS5Query(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestToMySQLBoolean(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "single term",
			query: "hello",
			want:  "+hello",
		},
		{
			name:  "implicit AND",
			query: "hello world",
			want:  "+hello +world",
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "three terms",
			query: "cat dog bird",
			want:  "+cat +dog +bird",
		},
		{
			name:  "phrase with multiple words",
			query: `"quick brown fox"`,
			want:  `"quick brown fox"`,
		},
		{
			name:  "empty query",
			query: "",
			want:  "",
		},
		{
			name:  "multiple phrases",
			query: `"hello world" "test case"`,
			want:  `"hello world" "test case"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchquery.ToMySQLBoolean(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToMySQLBoolean(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}
