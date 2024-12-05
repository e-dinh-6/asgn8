package main

import (
	"fmt"
	"testing"
)

// Tests for isValidId
func TestIsValidId(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected bool
	}{
		{"if", false},
		{"bind", false},
		{"validId", true},
		{"=>", false},
		{42, false},
		{"=", false},
	}

	for _, test := range tests {
		result := isValidId(test.input)
		if result != test.expected {
			t.Errorf("isValidId(%v) = %v; expected %v", test.input, result, test.expected)
		}
	}
}

func TestParseNum(t *testing.T) {
	expected := NumC{n: 42.0}
	result, err := parse(42.0)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if num, ok := result.(NumC); ok {
		if num.n != expected.n {
			t.Errorf("Expected %v, but got %v", expected.n, num.n)
		}
	} else {
		t.Errorf("Expected NumC, but got %T", result)
	}
}

func TestParseStr(t *testing.T) {
	expected := StrC{str: "hi"}
	result, err := parse("hi")
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	if num, ok := result.(StrC); ok {
		if num.str != expected.str {
			t.Errorf("Expected %v, but got %v", expected.str, num.str)
		}
	} else {
		t.Errorf("Expected NumC, but got %T", result)
	}
}

func TestParseIf(t *testing.T) {
	expected := IfC{AppC{IdC{"equal?"}, []ExprC{NumC{3}, NumC{3}}}, 
					BoolV{true},
					BoolV{false}}
	expr := []interface{}{
		[]interface{}{"if", []interface{}{"equal?",3.0,3.0},
		true,
		false},
	}
	result, err := parse(expr)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	fmt.Println(expected)
	fmt.Println(result)
}

// func TestParseSymbol(t *testing.T) {
// 	expected := IdC{}
// 	result, err := parse("hi")
// 	if err != nil {
// 		t.Fatalf("Expected no error, but got: %v", err)
// 	}

// 	if num, ok := result.(StrC); ok {
// 		if num.str != expected {
// 			t.Errorf("Expected %v, but got %v", expected, num)
// 		}
// 	} else {
// 		t.Errorf("Expected NumC, but got %T", result)
// 	}
// }