package main

import (
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected Input
		errMsg   string
	}{
		{
			name:  "case1",
			input: []byte("1\n*\n3\n3"),
			expected: Input{
				Id:   "1",
				Op:   "*",
				Val1: 3,
				Val2: 3,
			},
			errMsg: "",
		},
		{
			name:   "too few lines",
			input:  []byte("1\n+\n3"),
			errMsg: "invalid input format",
		},
		{
			name:  "case2",
			input: []byte("2\n/\n10\n5"),
			expected: Input{
				Id:   "2",
				Op:   "/",
				Val1: 10,
				Val2: 5,
			},
			errMsg: "",
		},
		{
			name:     "case3",
			input:    []byte("2\n/\n\n5"),
			expected: Input{},
			errMsg:   "invalid syntax",
		},
		{
			name:     "case4",
			input:    []byte("2\n/\n10\n"),
			expected: Input{},
			errMsg:   "invalid syntax",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if !strings.Contains(errMsg, tt.errMsg) {
				t.Errorf("Expected err message `%s`, got `%s`", tt.errMsg, errMsg)
			}

		})
	}
}
func FuzzParse(f *testing.F) {
	testcases := [][]byte{
		[]byte("1\n*\n3\n3"),
		[]byte("0\n"),
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, a []byte) {
		_, err := parser(a)
		if reason := shouldSkip(err); reason != "" {
			t.Skip(reason)
		}

	})

}

func shouldSkip(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "invalid input format"):
		return "input too short"
	case strings.Contains(msg, "invalid Id"):
		return "invalid Id"
	case strings.Contains(msg, "invalid op"):
		return "invalid operator"
	case strings.Contains(msg, "invalid syntax"): // Atoi error
		return "not a number"
	default:
		// можно вернуть пустую строку, если хочешь упасть на неожиданных ошибках
		return ""
	}
}
