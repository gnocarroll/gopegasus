package parser

import (
	"errors"
	"log"
	"pegasus/scanner"
	"strconv"
	"unicode/utf8"
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

	value, err := strconv.ParseFloat(
		tok.Text,
		64,
	)

	if err != nil {
		log.Printf("Expected float in \"%s\"", tok.Text)
		value = 0.0
	}

	ret := &FloatLiteral{
		Value: value,
	}
	ret.SetPosition(tok.Line, tok.Column)

	return ret, nil
}

func StringLiteralFromTok(tok *scanner.Token) (*StringLiteral, error) {
	if tok.TType != scanner.TOK_STRING {
		return nil, errors.New("expected string token")
	}

	runeCount := utf8.RuneCountInString(tok.Text)

	if runeCount < 2 {
		return nil, errors.New("expected at least two runes to hold quotes")
	}

	retIdx := 0
	retRunes := make([]rune, runeCount)
	isEscaped := false

	for i, r := range tok.Text {
		// check for quote at beginning, end
		if i == 0 || i == runeCount-1 {
			if r != '"' {
				return nil, errors.New("expected quote at first and last pos")
			} else { // found quote as expected
				continue
			}
		}

		if isEscaped {
			switch r {
			case 'n': // Line Feed
				retRunes[retIdx] = '\n'
			case 't': // Tab
				retRunes[retIdx] = '\t'
			case 'r': // Carriage Return
				retRunes[retIdx] = '\r'
			case 'v': // Vertical Tab
				retRunes[retIdx] = '\v'
			case 'f': // Form Feed
				retRunes[retIdx] = '\f'
			default:
				retRunes[retIdx] = r
			}
		} else if r == '\\' {
			isEscaped = true
			retIdx-- // will be incremented before next loop
		} else {
			retRunes[retIdx] = r
		}

		retIdx++
	}

	ret := &StringLiteral{
		Text: string(retRunes[:retIdx]),
	}
	ret.SetPosition(tok.Line, tok.Column)

	return ret, nil
}
