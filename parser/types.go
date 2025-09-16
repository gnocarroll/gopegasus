package parser

import (
	"pegasus/scanner"
	"sync/atomic"
)

type Parser struct {
	scan *scanner.Scanner

	nodeChan *chan INode
	errChan  *chan ParseError

	errCount atomic.Uint32
}

type ParseError struct {
	Expected scanner.TokenType
	Found    scanner.Token

	ExpectedNode INode

	Message string
}

func (err *ParseError) Error() string {
	return ""
}
