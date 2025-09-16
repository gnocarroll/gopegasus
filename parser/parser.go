package parser

import (
	"fmt"
	"pegasus/scanner"
)

const MAX_BUFFERED_NODES = 1000000
const MAX_PARSE_ERRORS = 100

func (parser *Parser) initParser() {
	nodeChan := make(chan INode, MAX_BUFFERED_NODES)
	errChan := make(chan ParseError, MAX_PARSE_ERRORS)

	parser.nodeChan = &nodeChan
	parser.errChan = &errChan
}

func (parser *Parser) ErrorCount() int {
	return int(parser.errCount.Load())
}

func NewParser(scan *scanner.Scanner) *Parser {
	var parser Parser

	parser.initParser()

	parser.scan = scan

	return &parser
}

func (parser *Parser) SetScanner(scan *scanner.Scanner) {
	parser.scan = scan
}

func (parser *Parser) send(node INode) {
	if node != nil {
		*parser.nodeChan <- node
	}
}

func (parser *Parser) addError(err *ParseError) {
	if parser.errCount.Load() < MAX_PARSE_ERRORS {
		parser.errCount.Add(1)

		*parser.errChan <- *err
	}
}

func (parser *Parser) accept(ttype scanner.TokenType) (*scanner.Token, error) {
	t := parser.scan.Peek()

	if t.TType == ttype {
		parser.scan.Advance()
		return &t, nil
	}

	e := ParseError{
		Expected: ttype,
		Found:    t,
	}

	parser.addError(&e)

	return nil, &e
}

func (parser *Parser) expectedNode(node INode) {
	e := ParseError{
		ExpectedNode: node,
		Message: fmt.Sprintf(
			"Expected but did not find node of type %T",
			node,
		),
	}

	parser.addError(&e)
}

func (parser *Parser) malformed(tok *scanner.Token) {
	e := &ParseError{
		Expected: tok.TType,
		Found:    *tok,
		Message:  "malformed token",
	}

	parser.addError(e)
}

func (parser *Parser) parse() {
	parser.send(parser.parseFile())

}

func (parser *Parser) Parse() {
	if parser.scan == nil {
		return
	}
	if parser.nodeChan == nil {
		parser.initParser()
	}

	go parser.parse()
}

func (parser *Parser) parseFile() *File {
	var f File

	for {
		def := parser.parseDefinition()

		if def == nil {
			break
		}

		f.definitions = append(f.definitions, def)
	}

	f.SetPosition(1, 1)

	return &f
}

func (parser *Parser) parseDefinition() *Definition {
	next := parser.scan.Peek()

	switch next.TType {
	case scanner.TOK_STRUCT, scanner.TOK_CLASS:
		return parser.parseTypeDef()
	case scanner.TOK_ENUM:
		return parser.parseEnumDef()
	case scanner.TOK_FUNCTION:
		return parser.parseFunctionDef()
	case scanner.TOK_IDENT:
		return parser.parseAssignment()
	default:
	}

	return nil
}

func (parser *Parser) parseTypeDef() *Definition {
	return nil
}

func (parser *Parser) parseEnumDef() *Definition {
	next := parser.scan.Peek()

	if next.TType != scanner.TOK_ENUM {
		return nil
	}

	return nil
}

func (parser *Parser) parseFunctionDef() *Definition {
	return nil
}

func (parser *Parser) parseAssignment() *Definition {
	next := parser.scan.Peek()

	if next.TType != scanner.TOK_IDENT {
		return nil
	}

	parser.scan.Advance()

	var ret Definition
	ret.SetPosition(next.Line, next.Column)

	next = parser.scan.Peek()

	if next.TType == scanner.TOK_COLON {
		ret.Type = parser.parseExpr()

		parser.accept(scanner.TOK_EQ)
	} else {
		ret.InferType = true

		parser.accept(scanner.TOK_COLON_EQ)
	}

	ret.Value = parser.parseExpr()

	return &ret
}

func (parser *Parser) parseCallArgs() CallArgs {
	var args CallArgs

	tok := parser.scan.Peek()

	args.SetPosition(tok.Line, tok.Column)

	for {
		tok1 := parser.scan.Peek()
		tok2 := parser.scan.PeekSecond()

		if tok1.TType == scanner.TOK_R_PAREN {
			break
		}

		var arg CallArg

		arg.SetPosition(tok1.Line, tok1.Column)

		if tok1.TType == scanner.TOK_IDENT &&
			tok2.TType == scanner.TOK_EQ {
			parser.scan.Advance()
			parser.scan.Advance()

			arg.Name = tok1.Text
		}

		arg.Value = parser.parseExpr()

		if arg.Value == nil {
			var expr IExpr = &Expr{}

			tok := parser.scan.Peek()

			expr.SetPosition(tok.Line, tok.Column)

			// failed to find expression when one was expected

			parser.expectedNode(expr)

			arg.Value = expr
		}

		args.ArgList = append(args.ArgList, arg)

		comma := parser.scan.Peek()

		if comma.TType != scanner.TOK_COMMA {
			break
		}

		parser.scan.Advance()
	}

	return args
}
