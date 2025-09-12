package parser

import (
	"pegasus/scanner"
)

const MAX_BUFFERED_NODES = 1000000

func (parser *Parser) initParser() {
	nodeChan := make(chan INode, MAX_BUFFERED_NODES)

	parser.nodeChan = &nodeChan
}

func NewParser() Parser {
	var parser Parser

	parser.initParser()

	return parser
}

func (parser *Parser) parse(scan *scanner.Scanner) {
	f := parser.parseFile(scan)

	*parser.nodeChan <- f
}

func (parser *Parser) Parse(scan *scanner.Scanner) {
	if parser.nodeChan == nil {
		parser.initParser()
	}

	go parser.parse(scan)
}

func (parser *Parser) parseFile(scan *scanner.Scanner) *File {
	return nil
}
