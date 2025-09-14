package parser

import (
	"errors"
	"log"
	"pegasus/scanner"
	"strconv"
)

func IntegerLiteralFromTok(tok *scanner.Token) (*IntegerLiteral, error) {
	if tok.TType != scanner.TOK_INTEGER {
		return nil, errors.New("expected integer token")
	}

	value, err := strconv.ParseUint(
		tok.Text,
		0,
		64,
	)

	if err != nil {
		log.Printf("Expected integer in \"%s\"", tok.Text)
		value = 0
	}

	ret := &IntegerLiteral{
		Value: value,
	}
	ret.SetPosition(tok.Line, tok.Column)

	return ret, nil
}

func FloatLiteralFromTok(tok *scanner.Token) (*FloatLiteral, error) {
	if tok.TType != scanner.TOK_FLOAT {
		return nil, errors.New("expected float token")
	}

}

func StringLiteralFromTok(tok *scanner.Token) (*StringLiteral, error) {
	if tok.TType != scanner.TOK_STRING {
		return nil, errors.New("expected float token")
	}
}
