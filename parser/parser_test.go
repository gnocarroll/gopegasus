package parser

import (
	"pegasus/scanner"
	"testing"
)

func createBinary(ttype scanner.TokenType, lhs IExpr, rhs IExpr) IExpr {
	return &BinaryExpr{
		Operator: scanner.Token{
			TType: ttype,
		},
		Lhs: lhs,
		Rhs: rhs,
	}
}

func createUnary(ttype scanner.TokenType, subexpr IExpr) IExpr {
	return &UnaryExpr{
		Operator: scanner.Token{
			TType: ttype,
		},
		SubExpr: subexpr,
	}
}

func TestPrintExpr(t *testing.T) {
	exprs := [...]IExpr{
		createBinary(scanner.TOK_PLUS, &IntegerLiteral{Value: 5}, &IntegerLiteral{Value: 5}),
	}
	strs := [...]string{
		"(+ 5 5)",
	}

	nLoops := min(len(exprs), len(strs))

	for i := 0; i < nLoops; i++ {
		got := ExprToString(exprs[i])
		expected := strs[i]

		if got != expected {
			t.Errorf(
				"Expected \"%s\", got \"%s\"",
				expected,
				got,
			)
		}
	}
}

func TestParseExpr(t *testing.T) {

}
