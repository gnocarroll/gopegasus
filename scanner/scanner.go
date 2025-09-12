package scanner

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	TOK_EOF TokenType = iota
	TOK_L_PAREN
	TOK_R_PAREN
	TOK_L_BRACK
	TOK_R_BRACK
	TOK_SEMI
	TOK_FUNCTION
	TOK_LAMBDA
	TOK_IF
	TOK_FOR
	TOK_WHILE
	TOK_END
	TOK_GT
	TOK_LT
	TOK_GE
	TOK_LE
	TOK_EQ
	TOK_EQ_EQ
	TOK_COLON_EQ
	TOK_BANG_EQ
	TOK_NOT
	TOK_AND
	TOK_OR
	TOK_PLUS
	TOK_MINUS
	TOK_STAR
	TOK_STAR_STAR
	TOK_F_SLASH
	TOK_COLON
	TOK_STRING
	TOK_IDENT
	TOK_INTEGER
	TOK_FLOAT
)

var tokStrings = [...]string{
	TOK_L_PAREN:   "(",
	TOK_R_PAREN:   ")",
	TOK_L_BRACK:   "[",
	TOK_R_BRACK:   "]",
	TOK_SEMI:      ";",
	TOK_FUNCTION:  "function",
	TOK_IF:        "if",
	TOK_WHILE:     "while",
	TOK_LAMBDA:    "lambda",
	TOK_FOR:       "for",
	TOK_END:       "end",
	TOK_STAR:      "*",
	TOK_F_SLASH:   "/",
	TOK_STAR_STAR: "**",
	TOK_PLUS:      "+",
	TOK_MINUS:     "-",
	TOK_LT:        "<",
	TOK_GT:        ">",
	TOK_LE:        "<=",
	TOK_GE:        ">=",
	TOK_EQ:        "=",
	TOK_COLON:     ":",
	TOK_COLON_EQ:  ":=",
	TOK_EQ_EQ:     "==",
	TOK_BANG_EQ:   "!=",
	TOK_NOT:       "not",
	TOK_AND:       "and",
	TOK_OR:        "or",
}

func (scanner Scanner) tokenize(s string) {
	sLen := len(s)
	runeIdx := 0

	for i := 0; i < sLen; {
		ttype, ok := tryTokStrings(s[i:])

		if ok {
			tstrLen := len(tokStrings[ttype])

			scanner.tChan <- Token{
				TType:  ttype,
				Line:   scanner.line,
				Column: scanner.column,
				Width:  tstrLen,
			}

			scanner.column += tstrLen

			continue
		}

		ttype, s, ok := tryTokFunctions(s, i)
	}
}

func tryTokStrings(s string) (TokenType, bool) {
	matchIdx := -1
	maxMatch := 0

	for tstrIdx, tstr := range tokStrings {
		tstrLen := len(tstr)

		if tstrLen <= maxMatch {
			continue
		}

		if strings.HasPrefix(s, tstr) {
			nextRune, _ := utf8.DecodeRuneInString(s[tstrLen:])

			if !unicode.IsSpace(nextRune) {
				matchIdx = tstrIdx
				maxMatch = tstrLen
			}
		}
	}

	if matchIdx == -1 {
		return TOK_EOF, false
	}

	return TokenType(matchIdx), true
}

func tryTokFunctions(s string, i int) (TokenType, string, bool) {
	// TODO: implement

	return TOK_EOF, "", false
}
