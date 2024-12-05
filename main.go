package main

import (
	"errors"
	"fmt"
	"strconv"
)

// custom type because Go doesn't have Symbols
type Symbol string

func (s Symbol) Equals(other Symbol) bool {
	return s == other
}

// ExprC ----------------
type ExprC interface{}
type NumC struct {
	n float64
}
type StrC struct {
	str string
}
type IdC struct {
	name Symbol
}
type IfC struct {
	cond  ExprC
	True  ExprC
	False ExprC
}
type AppC struct {
	fun  ExprC
	args []ExprC
}
type LamC struct {
	args []Symbol
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
	params  []Symbol
	body    ExprC
	clo_env Env
}
type PrimV struct {
	operand Symbol
}

// Environment -----------
type Env []Binding
type Binding struct {
	name Symbol
	val  Value
}

var topEnv = []Binding{
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
}

func lookup(name Symbol, env Env) (Value, error) {
	for _, binding := range env {
		if binding.name == name {
			return binding.val, nil
		}
	}
	return nil, fmt.Errorf("AAQZ: name not found %s", name)
}

func interp(expr ExprC, env Env) (Value, error) {
	switch e := expr.(type) {
	case NumC:
		return NumV{n: e.n}, nil
	case StrC:
		return StrV{str: e.str}, nil
	case IdC:
		return lookup(e.name, env)
	case IfC:
		condVal, err := interp(e.cond, env)
		if err != nil {
			return nil, fmt.Errorf("AAQZ: error evaluating condition in IfC: %w", err)
		}
		switch boolean := condVal.(type) {
		case BoolV:
			if boolean.b == true {
				return interp(e.True, env)
			}
			return interp(e.False, env)
		default:
			return nil, fmt.Errorf("AAQZ: not passing in conditional, got: %v", condVal)
		}
	case LamC:
		return CloV{params: e.args, body: e.body, clo_env: env}, nil
	case AppC:
		funVal, err := interp(e.fun, env)
		if err != nil {
			return nil, fmt.Errorf("AAQZ: invalid function, got: %w", err)
		}
		// eager evaluation on args
		var interped_args []Value
		for _, arg := range e.args {
			interped_arg, err := interp(arg, env)
			if err != nil {
				return nil, fmt.Errorf("AAQZ: error evaluating argument in AppC: %w", err)
			}
			interped_args = append(interped_args, interped_arg)
		}
		// checking type of interped fun
		switch f := funVal.(type) {
		case PrimV:
			result, err := applyPrim(f.operand, interped_args)
			if err != nil {
				return nil, fmt.Errorf("AAQZ: error applying primV %s: %w", f.operand, err)
			}
			return result, nil
		case CloV:
			// arity check
			if len(f.params) != len(interped_args) {
				return nil, fmt.Errorf("AAQZ: parame-arg count mismatch in AppC; expected length %d, got %d",
					len(f.params), len(interped_args))
			}
			// extend env
			newEnv := extendEnv(f.params, interped_args, f.clo_env)

			return interp(f.body, newEnv)
		default:
			return nil, fmt.Errorf("AAQZ: invalid function value in AppC: %v", funVal)
		}
	}
	return nil, fmt.Errorf("AAQZ: expression is not an ExprC %v", expr)
}

func extendEnv(params []Symbol, args []Value, env Env) Env {
	extendedEnv := append([]Binding{}, env...)
	for i := 0; i < len(params); i++ {
		extendedEnv = append(extendedEnv, Binding{name: params[i], val: args[i]})
	}
	return extendedEnv
}

func applyPrim(op Symbol, values []Value) (Value, error) {
	switch op {
	case "println":
		if len(values) != 1 {
			return nil, fmt.Errorf("AAQZ: `println` expects 1 argument, got: %v", values)
		}
		val := values[0]
		switch val := val.(type) {
		case StrV:
			fmt.Println(val.str)
			return BoolV{b: true}, nil
		default:
			return nil, fmt.Errorf("AAQZ `println` expects a string argument, got %v", val)
		}
	case "read-num":
		if len(values) != 0 {
			return nil, fmt.Errorf("AAQZ `read-num` expects 0 arguments, got: %v", values)
		}
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		num, other := strconv.ParseFloat(input, 64)
		if other != nil {
			return nil, fmt.Errorf("AAQZ Invalid input, not a real number, got: %v", other)
		}
		return NumV{n: num}, nil
	case "read-str":
		if len(values) != 0 {
			return nil, fmt.Errorf("AAQZ `read-str` expects 0 arguments, got: %v", values)
		}
		fmt.Print("> ")
		var input string
		fmt.Scanln(&input)
		return StrV{str: input}, nil
	case "seq":
		if len(values) == 0 {
			return nil, fmt.Errorf("AAQZ `seq` expects at least 1 argument, got: %v", values)
		}
		return values[len(values)-1], nil
	case "++":
		var result string
		for _, v := range values {
			switch v := v.(type) {
			case NumV:
				result += fmt.Sprintf("%f", v.n)
			case StrV:
				result += v.str
			case BoolV:
				result += fmt.Sprintf("%t", v.b)
			default:
				return nil, fmt.Errorf("AAQZ unsupported value type for ++: %v", v)
			}
		}
		return StrV{str: result}, nil
	case "equal?":
		if len(values) != 2 {
			return nil, fmt.Errorf("AAQZ `equal` expects 2 arguments, got: %v", values)
		}
		first, second := values[0], values[1]
		switch first := first.(type) {
		case NumV:
			if second, ok := second.(NumV); ok {
				return BoolV{b: first == second}, nil
			}
		case StrV:
			if second, ok := second.(StrV); ok {
				return BoolV{b: first == second}, nil
			}
		case BoolV:
			if second, ok := second.(BoolV); ok {
				return BoolV{b: first == second}, nil
			}
		}
		return BoolV{b: false}, nil
	case "error":
		if len(values) != 1 {
			return nil, fmt.Errorf("AAQZ `error` expects 1 argument, got: %v", values)
		}
		return nil, fmt.Errorf("AAQZ user-error: %v", values[0])
	default:
		if len(values) != 2 {
			return nil, fmt.Errorf("AAQZ Arithmetic operations expect 2 arguments, got: %v", len(values))
		}
		first, second := values[0], values[1]
		switch first := first.(type) {
		case NumV:
			second, ok := second.(NumV)
			if !ok {
				return nil, fmt.Errorf("AAQZ Expected numV for arithmetic operation, got: %v", []Value{first, second})
			}
			switch op {
			case "+":
				return NumV{n: first.n + second.n}, nil
			case "-":
				return NumV{n: first.n - second.n}, nil
			case "*":
				return NumV{n: first.n * second.n}, nil
			case "/":
				if second.n == 0 {
					return nil, fmt.Errorf("AAQZ Division by zero, got: %v", values)
				}
				return NumV{n: first.n / second.n}, nil
			case "<=":
				return BoolV{b: first.n <= second.n}, nil
			}
		default:
			return nil, fmt.Errorf("AAQZ Expected numV for arithmetic operation, got: %v", []Value{first, second})
		}
	}
	return nil, fmt.Errorf("AAQZ Unknown operation, got: %v", values)
}

func checkDupArgs(args []string) error {
	seen := map[string]bool{}
	for _, arg := range args {
		if seen[arg] {
			return errors.New("duplicate argument: " + arg)
		}
		seen[arg] = true
	}
	return nil
}

func isValidId(id interface{}) bool {
	// Check if `id` is a Symbol
	str, ok := id.(Symbol)
	if !ok {
		return false
	}

	reserved := map[Symbol]bool{
		"if":   true,
		"=":    true,
		"bind": true,
		"=>":   true,
	}

	if reserved[str] {
		return false
	}
	return true
}

func parse(s interface{}) (ExprC, error) {
	switch v := s.(type) {
	case float64: // <num>
		return NumC{n: v}, nil
	case string: 
		return StrC{str : v}, nil
	case Symbol: // <id>
		if isValidId(v) {
			return IdC{name: v}, nil
		}
		return nil, fmt.Errorf("invalid identifier: %v", v)
	case []interface{}: // List handling
		if len(v) == 0 {
			return nil, fmt.Errorf("invalid syntax: empty list")
		}

		switch head := v[0].(type) {
		case interface{}:
			switch head {
			case "if": // { if <expr> <expr> <expr> }
				if len(v) != 4 {
					return nil, fmt.Errorf("invalid if expression: %v", v)
				}
				cond, err := parse(v[1])
				if err != nil {
					return nil, err
				}
				thenExpr, err := parse(v[2])
				if err != nil {
					return nil, err
				}
				elseExpr, err := parse(v[3])
				if err != nil {
					return nil, err
				}
				return IfC{cond: cond, True: thenExpr, False: elseExpr}, nil
				// idk how to do bind

			default: // { <expr> <expr>* }
				funcExpr, err := parse(head)
				if err != nil {
					return nil, err
				}
				args := []ExprC{}
				for _, arg := range v[1:] {
					argExpr, err := parse(arg)
					if err != nil {
						return nil, err
					}
					args = append(args, argExpr)
				}
				return AppC{fun: funcExpr, args: args}, nil
			}
		default:
			return nil, fmt.Errorf("invalid expression head: %T", head)
		}

	default:
		return nil, fmt.Errorf("invalid syntax: %v", s)
	}
}

func main() {
	fmt.Println("Hi Chenyi, Katy, and Sharon")
	// testing println
	println_expr := AppC{
		fun:  IdC{name: "println"},
		args: []ExprC{StrC{str: "hello"}},
	}
	_, err1 := interp(println_expr, topEnv)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
	}

	// testing readstr
	readstr_expr := AppC{
		fun:  IdC{name: "read-str"},
		args: []ExprC{},
	}
	_, err2 := interp(readstr_expr, topEnv)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
	}

	// testing readnum
	readnum_expr := AppC{
		fun:  IdC{name: "read-num"},
		args: []ExprC{},
	}
	_, err3 := interp(readnum_expr, topEnv)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
	}

	// seq
	seq_expr := AppC{
		fun: IdC{name: "seq"},
		args: []ExprC{
			AppC{
				fun: IdC{name: "println"},
				args: []ExprC{
					StrC{"huh"},
				},
			},
			AppC{
				fun: IdC{name: "+"},
				args: []ExprC{
					NumC{n: 4},
					NumC{n: 5},
				},
			},
			AppC{
				fun: IdC{name: "println"},
				args: []ExprC{
					StrC{"added 4 + 5"},
				},
			},
		},
	}
	_, err4 := interp(seq_expr, topEnv)
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
	}

}