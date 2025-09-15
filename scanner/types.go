package scanner

type Scanner struct {
	tChan *chan Token

	line   int
	column int

	peek  Token
	peek2 Token

	// filled cached of two tokens (peek, peek2)
	isCacheFilled bool

	isEof bool
}

type Token struct {
	TType  TokenType
	Line   int
	Column int
	Width  int
	Text   string
}

type TokenType int
