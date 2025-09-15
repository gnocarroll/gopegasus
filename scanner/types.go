package scanner

type Scanner struct {
	tChan *chan Token

	line   int
	column int

	peek  Token
	peek2 Token

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
