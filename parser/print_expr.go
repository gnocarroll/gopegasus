package parser

import "strconv"

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
	default:
	}

	if s == "" {
		return "ERR"
	}

	return s
}
