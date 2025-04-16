package exercise03

import (
	"reflect"
	"testing"
)

func TestWordCount(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]int
	}{
		{
			name:     "empty string",
			input:    "",
			expected: map[string]int{},
		},
		{
			name:     "single word",
			input:    "go",
			expected: map[string]int{"go": 1},
		},
		{
			name:     "multiple words",
			input:    "hello world",
			expected: map[string]int{"hello": 1, "world": 1},
		},
		{
			name:     "repeated words",
			input:    "the quick brown fox jumps over the lazy fox",
			expected: map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 2, "jumps": 1, "over": 1, "lazy": 1},
		},
		{
			name:     "words with punctuation",
			input:    "hello, world!",
			expected: map[string]int{"hello,": 1, "world!": 1},
		},
		{
			name:     "words with different casing",
			input:    "Go go GO",
			expected: map[string]int{"Go": 1, "go": 1, "GO": 1},
		},
		{
			name:     "leading and trailing spaces",
			input:    "  leading and trailing  ",
			expected: map[string]int{"leading": 1, "and": 1, "trailing": 1},
		},
		{
			name:     "multiple spaces between words",
			input:    "multiple   spaces",
			expected: map[string]int{"multiple": 1, "spaces": 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := WordCount(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("WordCount(%q) = %v, expected %v", tc.input, actual, tc.expected)
			}
		})
	}
}
