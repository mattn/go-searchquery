package meilisearch_test

import (
	"testing"

	"github.com/mattn/go-searchquery/dialect/meilisearch"
)

func TestToQuery(t *testing.T) {
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
			want:  "hello world",
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "three terms",
			query: "cat dog bird",
			want:  "cat dog bird",
		},
		{
			name:  "empty query",
			query: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := meilisearch.ToQuery(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToQuery(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestToFilter(t *testing.T) {
	tests := []struct {
		name  string
		field string
		query string
		want  string
	}{
		{
			name:  "single term",
			field: "title",
			query: "hello",
			want:  `title = "hello"`,
		},
		{
			name:  "implicit AND",
			field: "title",
			query: "hello world",
			want:  `(title = "hello" AND title = "world")`,
		},
		{
			name:  "phrase search",
			field: "title",
			query: `"hello world"`,
			want:  `title = "hello world"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := meilisearch.ToFilter(tt.field, tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToFilter(%q, %q) = %q, want %q", tt.field, tt.query, got, tt.want)
			}
		})
	}
}
