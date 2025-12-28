package elasticsearch_test

import (
	"encoding/json"
	"testing"

	"github.com/mattn/go-searchquery/dialect/elasticsearch"
)

func TestToQueryString(t *testing.T) {
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
			want:  "(hello AND world)",
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  `"hello world"`,
		},
		{
			name:  "three terms",
			query: "cat dog bird",
			want:  "((cat AND dog) AND bird)",
		},
		{
			name:  "empty query",
			query: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := elasticsearch.ToQueryString(tt.query)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("ToQueryString(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestToMatchQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		field string
		want  string
	}{
		{
			name:  "single term",
			query: "hello",
			field: "content",
			want:  `{"match":{"content":"hello"}}`,
		},
		{
			name:  "implicit AND",
			query: "hello world",
			field: "content",
			want:  `{"bool":{"must":[{"match":{"content":"hello"}},{"match":{"content":"world"}}]}}`,
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			field: "content",
			want:  `{"match_phrase":{"content":"hello world"}}`,
		},
		{
			name:  "empty query",
			query: "",
			field: "content",
			want:  "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := elasticsearch.ToMatchQuery(tt.query, tt.field)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			
			// Compare as JSON to ignore formatting differences
			var gotJSON, wantJSON interface{}
			if err := json.Unmarshal([]byte(got), &gotJSON); err != nil {
				t.Fatalf("failed to parse got JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.want), &wantJSON); err != nil {
				t.Fatalf("failed to parse want JSON: %v", err)
			}
			
			gotStr, _ := json.Marshal(gotJSON)
			wantStr, _ := json.Marshal(wantJSON)
			
			if string(gotStr) != string(wantStr) {
				t.Errorf("ToMatchQuery(%q, %q) = %s, want %s", tt.query, tt.field, gotStr, wantStr)
			}
		})
	}
}
