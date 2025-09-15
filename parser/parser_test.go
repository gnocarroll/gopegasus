package parser

import (
	"errors"
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
		createUnary(scanner.TOK_PLUS, &FloatLiteral{Value: 10}),
	}
	strs := [...]string{
		"(+ 5 5)",
		"(+ 10.000)",
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

func parseExprForTest(s string) (string, error) {
	scan := scanner.NewScanner()

	scan.Tokenize(s)

	parse := NewParser(scan)

	expr := parse.parseExpr()

	if parse.ErrorCount() > 0 {
		return "", errors.New("non-zero error count")
	}

	return ExprToString(expr), nil
}

func TestParseValidExprs(t *testing.T) {
	exprs := [...]string{
		"5 + 5",
		"5 + 3 * 7",
		"5 + 3 * \"hello\"",
		"3. + 1. * 2. ** 4.0E10",
		"x",
		"A::B::x",
		"2 ** A::B::x",
		"(1 + 2) * 3",
		"x + 2",
		"+++5",
	}
	outputs := [...]string{
		"(+ 5 5)",
		"(+ 5 (* 3 7))",
		"(+ 5 (* 3 hello))",
		"(+ 3.000 (* 1.000 (** 2.000 40000000000.000)))",
		"x",
		"A::B::x",
		"(** 2 A::B::x)",
		"(* (+ 1 2) 3)",
		"(+ x 2)",
		"(+ (+ (+ 5)))",
	}

	nLoops := min(len(exprs), len(outputs))

	for i := 0; i < nLoops; i++ {
		got, err := parseExprForTest(exprs[i])

		if err != nil {
			t.Errorf("Unexpected err while parsing expr")
			continue
		}

		expected := outputs[i]

		if got != expected {
			t.Errorf(
				"Expected \"%s\", got \"%s\"",
				expected,
				got,
			)
		}
	}
}
