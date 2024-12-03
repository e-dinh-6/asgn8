package main

import "testing"

func TestNumCtoV(t *testing.T) {
	topEnv := []Binding{
		{"+", PrimV{"+"}},
		{"-", PrimV{"-"}},
		{"*", PrimV{"*"}},
		{"/", PrimV{"/"}},
		{"<=", PrimV{"<="}},
		{"error", PrimV{"error"}},
		{"equal?", PrimV{"equal?"}},
		{"true", BoolV{true}},
		{"false", BoolV{false}},
		{"println", PrimV{"println"}},
		{"read-num", PrimV{"read-num"}},
		{"read-str", PrimV{"read-str"}},
		{"++", PrimV{"++"}},
	}
	got1, got2 := interp(NumC{12}, topEnv)
	want1 := NumV{12}

	if got1 != want1 {
		t.Errorf("got %v, wanted %v", got1, want1)
	}

	if got2 != nil {
		t.Errorf("got %v, wanted %v", got2, nil)
	}
}

func TestInterpNumC(t *testing.T) {
	expr := NumC{42}
	expectedValue := NumV{42}
	env := Env{}
	result, err := interp(expr, env)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestInterpStrC(t *testing.T) {
	expr := StrC{str: "go!"}
	expectedValue := StrV{str: "go!"}
	env := Env{}
	result, err := interp(expr, env)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestInterpIdC(t *testing.T) {
	expr := IdC{name: "+"}
	expectedValue := PrimV{operand: "+"}
	topEnv := []Binding{
		{"+", PrimV{"+"}},
		{"-", PrimV{"-"}},
		{"*", PrimV{"*"}},
		{"/", PrimV{"/"}},
		{"<=", PrimV{"<="}},
		{"error", PrimV{"error"}},
		{"equal?", PrimV{"equal?"}},
		{"true", BoolV{true}},
		{"false", BoolV{false}},
		{"println", PrimV{"println"}},
		{"read-num", PrimV{"read-num"}},
		{"read-str", PrimV{"read-str"}},
		{"++", PrimV{"++"}},
	}
	result, err := interp(expr, topEnv)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}

func TestApplyPrim(t *testing.T) {
	value := []Value{NumV{n: 1}, NumV{n: 1}}
    expectedValue := BoolV{b: true}
	result, err := applyPrim("equal?", value)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if result != expectedValue {
		t.Errorf("Expected %v, got %v", expectedValue, result)
	}
}