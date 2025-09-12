package scanner

type Scanner struct {
	tChan chan Token

	line   int
	column int
}

type Token struct {
	TType  TokenType
	Line   int
	Column int
	Width  int
}

type TokenType int
