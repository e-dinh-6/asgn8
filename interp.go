package main

import (
	"fmt"
	"strconv"
)

func Hello() string {
	return "hi"
}

func lookup(name string, env Env) (Value, error) {
	for _, binding := range env {
		if binding.name == name {
			return binding.val, nil
		}
	}
	return nil, fmt.Errorf("AAQZ: name not found %s", name)
}

// e:= expr.(type) is a type assertion switch case
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
	return nil, fmt.Errorf("AAQZ Unknown expression, got: %v", expr)
}

func applyPrim(op string, values []Value) (Value, error) {
	switch op{
	case "println":
		if len(values) != 1 {
			return nil, fmt.Errorf("AAQZ: `println` expects 1 argument, got: %v", values)
		}
		val := values[0]
		switch val := val.(type) {
		case StrV:
			fmt.Println(val)
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
			return nil, fmt.Errorf("AAQZ Arithmetic operations expect 2 arguments, got: %v", values)
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