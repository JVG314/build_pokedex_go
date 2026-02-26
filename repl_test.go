package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		}, {
			input:    " HELLO WORLD ",
			expected: []string{"hello", "world"},
		}, {
			input:    " hello   world ",
			expected: []string{"hello", "world"},
		}, {
			input:    " helloworld ",
			expected: []string{"helloworld"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("fail: %v, length: %v, expected: %v", c.input, len(actual), len(c.expected))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("fail: %v, word: %s, expected: %s", c.input, word, expectedWord)
			}
		}
	}
}
