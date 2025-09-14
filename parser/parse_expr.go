package parser

import (
	"log"
	"pegasus/scanner"
	"strconv"
)

var binaryOps [][]scanner.TokenType = [][]scanner.TokenType{
	{scanner.TOK_OR},
	{scanner.TOK_AND},

	// will differ from C by making bit ops more tightly binding than
	// comparison ops
	{scanner.TOK_EQ_EQ, scanner.TOK_BANG_EQ},
	{scanner.TOK_LT, scanner.TOK_LE, scanner.TOK_GT, scanner.TOK_GE},
	{scanner.TOK_PIPE},
	{scanner.TOK_CARROT},
	{scanner.TOK_AMPERSAND},
	{scanner.TOK_LT_LT, scanner.TOK_GT_GT},
	{scanner.TOK_PLUS, scanner.TOK_MINUS},
	{scanner.TOK_STAR, scanner.TOK_F_SLASH, scanner.TOK_PERCENT},

	// adding built in base ** exponent
	{scanner.TOK_STAR_STAR},
}

var maxPrec int = len(binaryOps) - 1

var unaryOps []scanner.TokenType = []scanner.TokenType{
	scanner.TOK_PLUS,
	scanner.TOK_MINUS,
	scanner.TOK_NOT,
}

func (parser *Parser) parseExpr() IExpr {
	return parser.parseBinaryExpr()
}

func (parser *Parser) parseBinaryExpr() IExpr {
	return parser.parseBinaryExprPrec(0)
}

func (parser *Parser) parseBinaryExprPrec(prec int) IExpr {
	if prec < 0 {
		prec = 0
	}
	if prec > maxPrec {
		prec = maxPrec
	}

	parseNext := func() IExpr {
		if prec >= maxPrec {
			return parser.parseUnaryExpr()
		}

		return parser.parseBinaryExprPrec(prec + 1)
	}

	expr := parseNext()

	if expr == nil {
		return expr
	}

	line, column := expr.Position()

	for {
		foundOp := false
		nextTok := parser.scan.Peek()

		for _, op := range binaryOps[prec] {
			if nextTok.TType == op {
				foundOp = true
				break
			}
		}

		if !foundOp {
			break
		}

		parser.scan.Advance()

		expr = &BinaryExpr{
			Operator: nextTok,
			Lhs:      expr,
			Rhs:      parseNext(),
		}
		expr.SetPosition(line, column)
	}

	return expr
}

// unary operator in front
func (parser *Parser) parseUnaryExpr() IExpr {
	foundOp := false
	nextTok := parser.scan.Peek()

	for _, op := range unaryOps {
		if nextTok.TType == op {
			foundOp = true
			break
		}
	}

	if foundOp {
		parser.scan.Advance()
	}

	subExpr := parser.parsePostfixExpr()

	if !foundOp {
		return subExpr
	}

	return &UnaryExpr{
		Operator: nextTok,
		SubExpr:  subExpr,
	}
}

// e.g. function call, member access
func (parser *Parser) parsePostfixExpr() IExpr {
	subExpr := parser.parsePrimaryExpr()

	if subExpr == nil {
		return nil
	}

	return nil
}

func (parser *Parser) parsePrimaryExpr() IExpr {
	var ret IExpr = nil
	nextTok := parser.scan.Peek()

	switch nextTok.TType {
	case scanner.TOK_L_PAREN:
		ret = parser.parseExpr()
		parser.accept(scanner.TOK_R_PAREN)
	case scanner.TOK_INTEGER:
		value, err := strconv.ParseUint(
			nextTok.Text,
			0,
			64,
		)

		if err != nil {
			log.Printf("Expected integer in \"%s\"", nextTok.Text)
			value = 0
		}

		ret = &IntegerLiteral{
			Value: value,
		}
	case scanner.TOK_FLOAT:
		value, err := strconv.ParseFloat(
			nextTok.Text,
			64,
		)

		if err != nil {
			log.Printf("Expected float in \"%s\"", nextTok.Text)
			value = 0.0
		}

		ret = &FloatLiteral{
			Value: value,
		}
	case scanner.TOK_STRING:

	default:
		return nil
	}

	if nextTok.TType != scanner.TOK_L_PAREN {
		// => should set position

		ret.SetPosition(nextTok.Line, nextTok.Column)
	}

	parser.scan.Advance()

	return nil
}
