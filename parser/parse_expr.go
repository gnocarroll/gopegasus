package parser

import "pegasus/scanner"

var binaryOps [][]scanner.TokenType = [][]scanner.TokenType{
	{scanner.TOK_OR},
	{scanner.TOK_AND},
	{scanner.TOK_PIPE},
	{scanner.TOK_CARROT},
	{scanner.TOK_AMPERSAND},
}

func (parser *Parser) parseExpr() IExpr {
	return nil
}
