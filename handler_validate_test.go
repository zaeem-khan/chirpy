package main

import "testing"

func TestProfanityCleanup(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"hello world", "hello world"},
		{"KERFUFFLE is bad", "**** is bad"},
		{"sharBERT and forNAx are bad", "**** and **** are bad"},
	}

	for _, tc := range cases {
		result := profanityCleanup(tc.input)
		if result != tc.expected {
			t.Errorf("profanityCleanup(%q) = %q; expected %q", tc.input, result, tc.expected)
		}
	}
}
