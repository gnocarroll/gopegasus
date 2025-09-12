package scanner

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	TOK_EOF TokenType = iota
	TOK_FAILURE
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
	TOK_PERIOD
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
	TOK_PERIOD:    ".",
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

func (scanner *Scanner) token(ttype TokenType) Token {
	tstr := ""
	tstrlen := 0
	ttypeInt := int(ttype)

	if ttypeInt >= 0 && ttypeInt < len(tokStrings) {
		tstr = tokStrings[ttype]
		tstrlen = len(tstr)
	}

	return Token{
		TType:  ttype,
		Line:   scanner.line,
		Column: scanner.column,
		Width:  tstrlen,
		Text:   tstr,
	}
}

func (scanner *Scanner) Advance() Token {
	if scanner.isEof {
		select {
		case scanner.peek = <-scanner.tChan:
			scanner.isEof = false
		default: // still EOF
			return scanner.peek
		}
	} else {
		scanner.peek = <-scanner.tChan
	}

	if scanner.peek.TType == TOK_EOF {
		scanner.isEof = true
	}

	return scanner.peek
}

func (scanner *Scanner) Peek() Token {
	// if upcoming token's type is EOF, try to Advance
	if scanner.peek.TType == TOK_EOF {
		scanner.Advance()
	}

	return scanner.peek
}

func (scanner *Scanner) Tokenize(s string) {
	sLen := len(s)

	for i := 0; i < sLen; {
		nbytes := consumeIgnored(s[i:])

		i += nbytes

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

		ttype, tstr, ok := tryTokFunctions(s[i:])

		if ok {
			tstrlen := len(tstr)

			scanner.tChan <- Token{
				TType:  ttype,
				Line:   scanner.line,
				Column: scanner.column,
				Width:  tstrlen,
				Text:   tstr,
			}

			scanner.column += tstrlen

			continue
		}

		// failed to parse token

		scanner.tChan <- scanner.token(TOK_FAILURE)
		return
	}
}

// consume ignored characters (comments, whitespace)
func consumeIgnored(s string) int {
	i := 0
	sLen := len(s)

	for i < sLen {
		r, bytes := utf8.DecodeRuneInString(s[i:])

		nextR, nextBytes := ' ', 1

		if i+bytes < sLen {
			nextR, nextBytes = utf8.DecodeRuneInString(s[i+bytes:])
		}

		// skip past comment
		if r == '/' && nextR == '/' {
			i += (bytes + nextBytes)

			for i < sLen {
				r, bytes := utf8.DecodeRuneInString(s[i:])

				i += bytes

				if r == '\n' {
					break
				}
			}

			continue
		}

		if !unicode.IsSpace(r) {
			break
		}

		// is space, add width of rune
		i += bytes
	}

	return i
}

// see if next token matches fixed-width string tokens
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

type ScanFunc func(string) (TokenType, string, bool)

// see if next token can be found by scanner funcs for variable-length tokens
// (e.g. integers, identifiers, etc.)
func tryTokFunctions(s string) (TokenType, string, bool) {
	scanFuncs := [...]ScanFunc{scanInteger, scanIdent, scanFloat}

	for _, scanFunc := range scanFuncs {
		ttype, s, ok := scanFunc(s)

		if !ok {
			continue
		}

		return ttype, s, ok
	}

	return TOK_EOF, "", false
}

func scanInteger(s string) (TokenType, string, bool) {
	width := 0

	for _, r := range s {
		if unicode.IsDigit(r) {
			width += utf8.RuneLen(r)
		} else {
			break
		}
	}

	if width == 0 {
		return TOK_EOF, "", false
	}

	return TOK_INTEGER, s[:width], true
}

func scanIdent(s string) (TokenType, string, bool) {
	width := 0

	r, bytes := utf8.DecodeRuneInString(s)

	if r != '_' && !unicode.IsLetter(r) {
		return TOK_EOF, "", false
	}

	width += bytes

	for _, r := range s[width:] {
		if r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			width += utf8.RuneLen(r)
		} else {
			break
		}
	}

	return TOK_IDENT, s[:width], true
}

func scanFloat(s string) (TokenType, string, bool) {
	return TOK_EOF, "", false
}
