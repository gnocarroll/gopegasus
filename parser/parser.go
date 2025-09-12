package parser

import (
	"pegasus/scanner"
)

type Parser struct {
}

func (parser *Parser) parse(scan *scanner.Scanner) {

}

func (parser *Parser) Parse(scan *scanner.Scanner) {
	go parser.parse(scan)
}
