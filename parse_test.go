package main

import "testing"

// Tests for isValidId
func TestIsValidId(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{Symbol("if"), false},
		{Symbol("bind"), false},
		{Symbol("validId"), true},
		{Symbol("=>"), false},
		{42, false},
		{Symbol("="), false},
	}

	for _, test := range tests {
		result := isValidId(test.input)
		if result != test.expected {
			t.Errorf("isValidId(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}
