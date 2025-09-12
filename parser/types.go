package parser

import (
	"pegasus/scanner"
)

type Parser struct {
	nodeChan *chan INode
}

type TokenError struct {
	Expected scanner.TokenType
	Found    scanner.Token
}
