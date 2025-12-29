package bleve

import (
	"testing"
)

func TestToQueryString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single term",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "implicit AND",
			input:    "hello world",
			expected: "+hello +world",
		},
		{
			name:     "explicit OR",
			input:    "hello OR world",
			expected: "hello world",
		},
		{
			name:     "phrase search",
			input:    `"hello world"`,
			expected: `"hello world"`,
		},
		{
			name:     "complex AND OR",
			input:    "cat AND (dog OR mouse)",
			expected: "+cat +(dog mouse)",
		},
		{
			name:     "multiple terms with AND",
			input:    "apple AND banana AND cherry",
			expected: "+apple +banana +cherry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToQueryString(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestToQueryStringErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "mismatched parentheses",
			input: "hello (world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToQueryString(tt.input)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
