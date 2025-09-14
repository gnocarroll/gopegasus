package parser

import "pegasus/scanner"

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

var nPrecLevels int = len(binaryOps)

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

	return nil
}

func (parser *Parser) parseUnaryExpr() IExpr {

}
