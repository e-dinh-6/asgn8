package main

import (
	"fmt"
)

// ExprC ----------------
type ExprC interface{}
type NumC struct {
	n float64
}
type AppC struct {
	fun  ExprC
	args []ExprC
}
type StrC struct {
	str string
}
type IfC struct {
	cond  ExprC
	True  ExprC
	False ExprC
}
type IdC struct {
	name string
}
type LamC struct {
	args []string
	body ExprC
}

// Values ---------------
type Value interface{}

type NumV struct {
	n float64
}

type StrV struct {
	str string
}

type BoolV struct {
	b bool
}
type CloV struct {
	params  []string
	body    ExprC
	clo_env Env
}
type PrimV struct {
	operand string
}

// Environment -----------
type Env []Binding
type Binding struct {
	name string
	val  Value
}

func main() {
	fmt.Println("Hi Chenyi, Katy, and Sharon")
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
	}
	fmt.Println(topEnv)
}
