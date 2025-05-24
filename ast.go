package main

import (
	"fmt"
	"math"
)

type Expr interface {
	Eval() (float64, error)
}

type BinaryOpNode struct {
	Op TokenType
	Left Expr
	Right Expr
}

func (bon BinaryOpNode) Eval() (float64, error) {
	var ret float64
	left, lerr := bon.Left.Eval()
	if lerr != nil {
		return 0, lerr
	}
	right, rerr := bon.Right.Eval()
	if rerr != nil {
		return 0, rerr
	}
	switch bon.Op {
	case PLUS:
		ret = left + right
	case MINUS:
		ret = left - right
	case MUL:
		ret = left * right
	case DIV:
		ret = left / right
	}
	return ret, nil
}

type ArgumentError string

func (ae ArgumentError) Error() string {
	return string(ae)
}

type Function func([]Expr) (float64, error)

var FunctionTable map[string]Function = map[string]Function {
	"sqrt": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			arg1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Sqrt(arg1), nil
		}
		return 0, ArgumentError("sqrt takes one argument")
	},
	"abs": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			arg1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Abs(arg1), nil
		}
		return 0, ArgumentError("abs takes one argument")
	},
	"pow": func(args []Expr) (float64, error) {
		if len(args) == 2 {
			args1, err1 := args[0].Eval()
			if err1 != nil {
				return 0, err1
			}
			args2, err2 := args[1].Eval()
			if err2 != nil {
				return 0, err2
			}
			return math.Pow(args1, args2), nil
		}
		return 0, ArgumentError("pow takes two arguments")
	},
	"cos": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Cos(args1), nil
		}
		return 0, ArgumentError("cos takes one argument")
	},
	"sin": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Sin(args1), nil
		}
		return 0, ArgumentError("sin takes one argument")
	},
	"tan": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Tan(args1), nil
		}
		return 0, ArgumentError("tan takes one argument")
	},
	"acos": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Acos(args1), nil
		}
		return 0, ArgumentError("acos takes one argument")
	},
	"asin": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Asin(args1), nil
		}
		return 0, ArgumentError("asin takes one argument")
	},
	"atan": func(args []Expr) (float64, error) {
		if len(args) == 1 {
			args1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Atan(args1), nil
		}
		return 0, ArgumentError("atan takes one argument")
	},
	"max": func(args []Expr) (float64, error) {
		if len(args) == 2 {
			arg1, err1 := args[0].Eval()
			if err1 != nil {
				return 0, err1
			}
			arg2, err2 := args[1].Eval()
			if err2 != nil {
				return 0, err2
			}
			return math.Max(arg1, arg2), nil
		}
		return 0, ArgumentError("max takes two arguments")
	},
	"min": func(args []Expr) (float64, error) {
		if len(args) == 2 {
			arg2, err := args[1].Eval()
			arg1, err := args[0].Eval()
			if err != nil {
				return 0, err
			}
			return math.Min(arg1, arg2), nil
		}
		return 0, ArgumentError("min takes two arguments")
	},
}

type FunctionNode struct {
	Name string
	Args []Expr
}

func (fn FunctionNode) Eval() (float64, error) {
	funcPtr, ok := FunctionTable[fn.Name]
	if ok {
		ans, err := funcPtr(fn.Args)
		if err != nil {
			return 0, err
		}
		return ans, nil
	}
	return 0, fmt.Errorf("function %v does not exist", fn.Name)
}

type NumberNode float64

func (nn NumberNode) Eval() (float64, error) {
	return float64(nn), nil
}
