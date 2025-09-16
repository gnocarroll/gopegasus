package parser

import (
	"strconv"
	"strings"
)

func callArgToString(arg *CallArg) string {
	ret := ""

	if arg.Name != "" {
		ret += (arg.Name + "=")
	}

	ret += ExprToString(arg.Value)

	return ret
}

func ExprToString(expr IExpr) string {
	s := ""

	switch expr := expr.(type) {
	case *BinaryExpr:
		s = "(" + expr.Operator.TType.Text() + " " +
			ExprToString(expr.Lhs) + " " +
			ExprToString(expr.Rhs) + ")"
	case *UnaryExpr:
		s = "(" + expr.Operator.TType.Text() + " " +
			ExprToString(expr.SubExpr) + ")"
	case *FloatLiteral:
		s = strconv.FormatFloat(
			expr.Value,
			'f',
			3,
			64,
		)
	case *StringLiteral:
		s = expr.Text
	case *IntegerLiteral:
		s = strconv.FormatUint(
			expr.Value,
			10,
		)
	case *IdentExpr:
		s = strings.Join(expr.Names, "::")
	case *FunctionCallExpr:
		argStrs := make([]string, len(expr.Args.ArgList))

		for i, arg := range expr.Args.ArgList {
			argStrs[i] = callArgToString(&arg)
		}

		s = "(" + ExprToString(expr.Function) + " " +
			strings.Join(argStrs, " ") + ")"
	default:
	}

	if s == "" {
		return "ERR"
	}

	return s
}
