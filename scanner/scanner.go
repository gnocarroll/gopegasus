package scanner

import (
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const MAX_BUFFERED_TOKENS = 1000000

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
	TOK_STRUCT
	TOK_CLASS
	TOK_ENUM
	TOK_VARIANT
	TOK_IF
	TOK_FOR
	TOK_WHILE
	TOK_BEGIN
	TOK_END
	TOK_GT
	TOK_LT
	TOK_GE
	TOK_LE
	TOK_EQ
	TOK_EQ_EQ
	TOK_LT_LT
	TOK_GT_GT
	TOK_PERCENT
	TOK_AMPERSAND
	TOK_CARROT
	TOK_TILDE
	TOK_PIPE
	TOK_COLON_EQ
	TOK_COLON_COLON
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

var TokStrings = [...]string{
	TOK_L_PAREN:     "(",
	TOK_R_PAREN:     ")",
	TOK_L_BRACK:     "[",
	TOK_R_BRACK:     "]",
	TOK_SEMI:        ";",
	TOK_FUNCTION:    "function",
	TOK_IF:          "if",
	TOK_WHILE:       "while",
	TOK_AMPERSAND:   "&",
	TOK_CARROT:      "^",
	TOK_TILDE:       "~",
	TOK_PIPE:        "|",
	TOK_LT_LT:       "<<",
	TOK_GT_GT:       ">>",
	TOK_PERCENT:     "%",
	TOK_LAMBDA:      "lambda",
	TOK_STRUCT:      "struct",
	TOK_CLASS:       "class",
	TOK_ENUM:        "enum",
	TOK_VARIANT:     "variant",
	TOK_COLON_COLON: "::",
	TOK_FOR:         "for",
	TOK_BEGIN:       "begin",
	TOK_END:         "end",
	TOK_STAR:        "*",
	TOK_F_SLASH:     "/",
	TOK_PERIOD:      ".",
	TOK_STAR_STAR:   "**",
	TOK_PLUS:        "+",
	TOK_MINUS:       "-",
	TOK_LT:          "<",
	TOK_GT:          ">",
	TOK_LE:          "<=",
	TOK_GE:          ">=",
	TOK_EQ:          "=",
	TOK_COLON:       ":",
	TOK_COLON_EQ:    ":=",
	TOK_EQ_EQ:       "==",
	TOK_BANG_EQ:     "!=",
	TOK_NOT:         "not",
	TOK_AND:         "and",
	TOK_OR:          "or",
}

var TokDescs = [...]string{
	TOK_EOF:       "End of File",
	TOK_FAILURE:   "Scanner Failure (Unable to Determine Token)",
	TOK_L_PAREN:   "Left Paren ('(')",
	TOK_R_PAREN:   "Right Paren (')')",
	TOK_L_BRACK:   "Left Bracket ('[')",
	TOK_R_BRACK:   "Right Bracket (']')",
	TOK_AMPERSAND: "Ampersand ('&')",
	TOK_CARROT:    "Carrot ('^')",
	TOK_TILDE:     "Tilde ('~')",
	TOK_PIPE:      "Pipe ('|')",
	TOK_LT_LT:     "Left Shift ('<<')",
	TOK_GT_GT:     "Right Shift ('>>')",
	TOK_PERCENT:   "Percent ('%')",
	TOK_SEMI:      "Semicolon (';')",
	TOK_FUNCTION:  "Function ('function')",
	TOK_LAMBDA:    "Lambda ('lambda')",
	TOK_IF:        "If ('if')",
	TOK_FOR:       "For ('for')",
	TOK_WHILE:     "While ('while')",
	TOK_STRUCT:    "Struct ('struct')",
	TOK_CLASS:     "Class ('class')",
	TOK_ENUM:      "Enum ('enum')",
	TOK_VARIANT:   "Variant ('variant')",
	TOK_BEGIN:     "Begin ('begin')",
	TOK_END:       "End ('end')",
	TOK_GT:        "Greater Than ('>')",
	TOK_LT:        "Less Than ('<')",
	TOK_GE:        "Greater Than Or Equal To ('>=')",
	TOK_LE:        "Less Than Or Equal To ('<=')",
	TOK_EQ:        "Assign ('=')",
	TOK_EQ_EQ:     "Equal To ('==')",
	TOK_COLON_EQ:  "Assign + Infer Type (':=')",
	TOK_BANG_EQ:   "Not Equal To ('!=')",
	TOK_NOT:       "Not ('not')",
	TOK_AND:       "And ('and')",
	TOK_OR:        "Or ('or')",
	TOK_PLUS:      "Plus ('+')",
	TOK_MINUS:     "Minus ('-')",
	TOK_STAR:      "Star ('*')",
	TOK_STAR_STAR: "Two Stars ('**')",
	TOK_PERIOD:    "Period ('.')",
	TOK_F_SLASH:   "Forward Slash ('/')",
	TOK_COLON:     "Colon (':')",
	TOK_STRING:    "String Literal",
	TOK_IDENT:     "Identifier",
	TOK_INTEGER:   "Integer",
	TOK_FLOAT:     "Floating-Point Number",
}

func (ttype TokenType) Text() string {
	if ttype < 0 || int(ttype) >= len(TokStrings) {
		return ""
	}

	return TokStrings[ttype]
}

func (ttype TokenType) Desc() string {
	if ttype < 0 || int(ttype) >= len(TokDescs) {
		return ""
	}

	return TokDescs[ttype]
}

func (scanner *Scanner) initScanner() {
	tChan := make(chan Token, MAX_BUFFERED_TOKENS)

	scanner.tChan = &tChan
	scanner.isEof = true

	scanner.line = 1
	scanner.column = 1
}

func NewScanner() *Scanner {
	var ret Scanner

	ret.initScanner()

	return &ret
}

func (scanner *Scanner) token(ttype TokenType) Token {
	tstr := ""
	tstrlen := 0
	ttypeInt := int(ttype)

	if ttypeInt >= 0 && ttypeInt < len(TokStrings) {
		tstr = TokStrings[ttype]
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

func (scanner *Scanner) fillCache() {
	if scanner.isCacheFilled {
		return
	}
	if scanner.tChan == nil {
		scanner.initScanner()
	}

	scanner.isCacheFilled = true

	// should always get at least one TOK_EOF
	scanner.peek = <-*scanner.tChan

	if scanner.peek.TType == TOK_EOF {
		scanner.isEof = true
	} else {
		// if did not find eof, then should be another token
		// available at some point
		scanner.peek2 = <-*scanner.tChan
	}
}

func (scanner *Scanner) Advance() Token {
	// attempt to load first two tokens if not already done

	scanner.fillCache()

	ret := scanner.peek

	var chanTok Token

	if scanner.isEof {
		select {
		case chanTok = <-*scanner.tChan:
			scanner.isEof = false
		default: // still EOF
			return ret
		}
	} else if scanner.peek2.TType != TOK_EOF {
		// if second cached tok is not EOF then attempt to get another
		chanTok = <-*scanner.tChan
	}

	scanner.peek = scanner.peek2
	scanner.peek2 = chanTok

	if scanner.peek.TType == TOK_EOF {
		scanner.isEof = true
	}

	return ret
}

func (scanner *Scanner) Peek() Token {
	scanner.fillCache()

	return scanner.peek
}

// returns the token after the one Peek() returns
func (scanner *Scanner) PeekSecond() Token {
	scanner.fillCache()

	return scanner.peek2
}

func (scanner *Scanner) tokenize(s string) {
	sLen := len(s)

	for i := 0; i < sLen; {
		nbytes := scanner.consumeIgnored(s[i:])

		i += nbytes

		ttype, tstr, ok := tryTokFunctions(s[i:])

		if ok {
			tstrLen := len(tstr)

			*scanner.tChan <- Token{
				TType:  ttype,
				Line:   scanner.line,
				Column: scanner.column,
				Width:  tstrLen,
				Text:   tstr,
			}

			i += tstrLen
			scanner.column += utf8.RuneCountInString(tstr)

			continue
		}

		ttype, ok = tryTokStrings(s[i:])

		if ok {
			tstrLen := len(TokStrings[ttype])

			*scanner.tChan <- Token{
				TType:  ttype,
				Line:   scanner.line,
				Column: scanner.column,
				Width:  tstrLen,
			}

			i += tstrLen
			scanner.column += utf8.RuneCountInString(tstr)

			continue
		}

		// failed to parse token

		*scanner.tChan <- scanner.token(TOK_FAILURE)
		break
	}

	*scanner.tChan <- scanner.token(TOK_EOF)
}

func (scanner *Scanner) Tokenize(s string) {
	if scanner.tChan == nil {
		scanner.initScanner()
	}

	scanner.isEof = false

	// do tokenization work in separate goroutine

	go scanner.tokenize(s)
}

func (scanner *Scanner) TokenizeFile(filepath string) {
	bytes, err := os.ReadFile(filepath)

	if err != nil {
		log.Fatalf("TokenizeFile failed to ReadFile: %s\n", err)
		return
	}

	scanner.Tokenize(string(bytes))
}

// consume ignored characters (comments, whitespace)
func (scanner *Scanner) consumeIgnored(s string) int {
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
			scanner.column += 2 // for two forward slashes
			i += (bytes + nextBytes)

			for i < sLen {
				r, bytes := utf8.DecodeRuneInString(s[i:])

				scanner.column++
				i += bytes

				if r == '\n' {
					scanner.line++
					scanner.column = 1
					break
				}
			}

			continue
		}

		if !unicode.IsSpace(r) {
			break
		}

		if r == '\n' {
			scanner.line++
			scanner.column = 1
		} else {
			scanner.column++
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

	for tstrIdx, tstr := range TokStrings {
		tstrLen := len(tstr)

		if tstrLen <= maxMatch ||
			!strings.HasPrefix(s, tstr) {
			continue
		}

		// passed checks, update max match

		matchIdx = tstrIdx
		maxMatch = tstrLen
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
	scanFuncs := [...]ScanFunc{
		scanInteger,
		scanIdent,
		scanFloat,
		scanString,
	}

	maxMatch := 0
	retTType := TOK_EOF
	retStr := ""

	for _, scanFunc := range scanFuncs {
		ttype, foundS, ok := scanFunc(s)

		if !ok {
			continue
		}

		if len(foundS) > maxMatch {
			maxMatch = len(foundS)

			retTType = ttype
			retStr = foundS
		}
	}

	if maxMatch <= 0 {
		return TOK_EOF, "", false
	}

	return retTType, retStr, true
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

	identStr := s[:width]

	// ensure that identStr does not match any keywords

	for _, tstr := range TokStrings {
		if identStr == tstr {
			return TOK_EOF, "", false
		}
	}

	return TOK_IDENT, identStr, true
}

func scanFloat(s string) (TokenType, string, bool) {
	pointOffset := -1
	pointEndOffset := -1

	eOffset := -1
	eEndOffset := -1
	digitsAfterE := false

	totalWidth := 0
	sLen := len(s)

	for totalWidth < sLen {
		r, width := utf8.DecodeRuneInString(s[totalWidth:])

		if unicode.IsDigit(r) {
			if eOffset != -1 {
				digitsAfterE = true
			}
		} else if r == '.' {
			if pointOffset != -1 {
				break
			}

			pointOffset = totalWidth
			pointEndOffset = pointOffset + width
		} else if r == 'e' || r == 'E' {
			if eOffset != -1 {
				break
			}

			eOffset = totalWidth
			eEndOffset = eOffset + width
		} else if totalWidth == eEndOffset && (r == '+' || r == '-') {
			// okay to have plus/minus immediately after e/E
		} else {
			// character not part of float literal, exit loop
			break
		}

		totalWidth += width
	}

	// Conditions which indicate failure
	// - point not found
	// - no digits on left or right of point
	// - e/E was found but no digits after it

	if pointOffset == -1 ||
		(pointOffset == 0 && (pointEndOffset == eOffset || pointEndOffset == totalWidth)) ||
		(eOffset != -1 && !digitsAfterE) {
		return TOK_EOF, "", false
	}

	// success, return relevant slice

	return TOK_FLOAT, s[:totalWidth], true
}

func scanString(s string) (TokenType, string, bool) {
	foundEndQuote := false

	width := 0
	sLen := len(s)

	isEscaped := false

	for width < sLen {
		r, bytes := utf8.DecodeRuneInString(s[width:])

		if width == 0 && r != '"' {
			break
		}

		if isEscaped {
			isEscaped = false
		} else if width != 0 && r == '"' {
			foundEndQuote = true
			width += bytes
			break
		} else if r == '\\' {
			isEscaped = true
		}

		width += bytes
	}

	if width == 0 || !foundEndQuote {
		return TOK_EOF, "", false
	}

	return TOK_STRING, s[:width], true
}
