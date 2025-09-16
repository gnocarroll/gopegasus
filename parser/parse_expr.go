package parser

import (
	"pegasus/scanner"
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

	var subExpr IExpr

	if foundOp {
		parser.scan.Advance()
		subExpr = parser.parseUnaryExpr()
	} else {
		return parser.parsePostfixExpr()
	}

	ret := &UnaryExpr{
		Operator: nextTok,
		SubExpr:  subExpr,
	}
	ret.SetPosition(nextTok.Line, nextTok.Column)

	return ret
}

// e.g. function call, member access
func (parser *Parser) parsePostfixExpr() IExpr {
	subExpr := parser.parsePrimaryExpr()

	if subExpr == nil {
		return nil
	}

	ret := subExpr

	nextTok := parser.scan.Peek()

	switch nextTok.TType {
	// Function Call or Template Expansion
	case scanner.TOK_L_PAREN, scanner.TOK_L_BRACK:
		parser.scan.Advance()

		isTemplateCall := (nextTok.TType == scanner.TOK_L_BRACK)

		ret = &FunctionCallExpr{
			IsTemplateCall: isTemplateCall,
			Function:       subExpr,
			Args:           parser.parseCallArgs(),
		}
		ret.SetPosition(subExpr.Line(), subExpr.Column())

		if isTemplateCall {
			parser.accept(scanner.TOK_R_BRACK)
		} else {
			parser.accept(scanner.TOK_R_PAREN)
		}
	case scanner.TOK_PERIOD: // Member Access
		parser.scan.Advance()

		tok, _ := parser.accept(scanner.TOK_IDENT)

		member := ""

		if tok != nil {
			member = tok.Text
		}

		ret = &MemberAccessExpr{
			Instance: subExpr,
			Member:   member,
		}
		ret.SetPosition(subExpr.Line(), subExpr.Column())
	default: // No operation to parse here
	}

	return ret
}

func (parser *Parser) parsePrimaryExpr() IExpr {
	var ret IExpr = nil
	var err error = nil

	nextTok := parser.scan.Peek()

	switch nextTok.TType {
	case scanner.TOK_L_PAREN:
		parser.scan.Advance()

		ret = parser.parseExpr()

		parser.accept(scanner.TOK_R_PAREN)
	case scanner.TOK_INTEGER:
		ret, err = IntegerLiteralFromTok(&nextTok)
		parser.scan.Advance()
	case scanner.TOK_FLOAT:
		ret, err = FloatLiteralFromTok(&nextTok)
		parser.scan.Advance()
	case scanner.TOK_STRING:
		ret, err = StringLiteralFromTok(&nextTok)
		parser.scan.Advance()
	case scanner.TOK_IDENT:
		ret = parser.parseIdentExpr()
	default:
		return nil
	}

	if err != nil {
		parser.malformed(&nextTok)
		ret = &ErrorExpr{}
	}

	if nextTok.TType != scanner.TOK_L_PAREN {
		// => should set position

		ret.SetPosition(nextTok.Line, nextTok.Column)
	}

	return ret
}

func (parser *Parser) parseIdentExpr() IExpr {
	tok := parser.scan.Peek()

	if tok.TType != scanner.TOK_IDENT {
		return nil
	}

	parser.scan.Advance()

	ret := &IdentExpr{
		Names: []string{tok.Text},
	}

	for {
		tok = parser.scan.Peek()

		// Namespace Separator: "::"
		if tok.TType != scanner.TOK_COLON_COLON {
			break
		}

		parser.scan.Advance()

		tokRef, err := parser.accept(scanner.TOK_IDENT)

		if err != nil {
			break
		}

		ret.Names = append(ret.Names, tokRef.Text)
	}

	return ret
}
