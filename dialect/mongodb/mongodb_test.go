package mongodb_test

import (
	"encoding/json"
	"testing"

	"github.com/mattn/go-searchquery/dialect/mongodb"
)

func TestToTextSearch(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{
			name:  "single term",
			query: "hello",
			want:  `{"$text":{"$search":"\"hello\""}}`,
		},
		{
			name:  "implicit AND",
			query: "hello world",
			want:  `{"$text":{"$search":"\"hello\" \"world\""}}`,
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			want:  `{"$text":{"$search":"\"hello world\""}}`,
		},
		{
			name:  "empty query",
			query: "",
			want:  "{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mongodb.ToTextSearch(tt.query)
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
				t.Errorf("ToTextSearch(%q) = %s, want %s", tt.query, gotStr, wantStr)
			}
		})
	}
}

func TestToRegexQuery(t *testing.T) {
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
			want:  `{"content":{"$options":"i","$regex":"hello"}}`,
		},
		{
			name:  "implicit AND",
			query: "hello world",
			field: "content",
			want:  `{"$and":[{"content":{"$options":"i","$regex":"hello"}},{"content":{"$options":"i","$regex":"world"}}]}`,
		},
		{
			name:  "phrase search",
			query: `"hello world"`,
			field: "content",
			want:  `{"content":{"$options":"i","$regex":"hello world"}}`,
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
			got, err := mongodb.ToRegexQuery(tt.query, tt.field)
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
				t.Errorf("ToRegexQuery(%q, %q) = %s, want %s", tt.query, tt.field, gotStr, wantStr)
			}
		})
	}
}
