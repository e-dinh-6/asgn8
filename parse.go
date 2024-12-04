package main

import (
	"errors"
	"fmt"
	// "strconv"
	// "strings"
)

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
	// Check if `id` is a string
	str, ok := id.(string)
	if !ok {
		return false
	}

	reserved := map[string]bool{
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
		case string: // <id>
			if isValidId(v) {
				return IdC{name: v}, nil
			}
			return nil, fmt.Errorf("invalid identifier: %v", v)
		case []interface{}: // List handling
			if len(v) == 0 {
				return nil, fmt.Errorf("invalid syntax: empty list")
			}
	
			switch head := v[0].(type) {
			case string:
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
				return nil, fmt.Errorf("invalid expression head: %v", head)
			}
	
		default:
			return nil, fmt.Errorf("invalid syntax: %v", s)
		}
	}