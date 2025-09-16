package parser

import "pegasus/scanner"

type INode interface {
	nodeTag()

	Name() string
	Line() int
	Column() int

	Position() (int, int)
	SetPosition(int, int)
}

type Node struct {
	name   string
	line   int
	column int
}

func (node *Node) nodeTag() {}
func (node *Node) Name() string {
	return node.name
}
func (node *Node) Line() int {
	return node.line
}
func (node *Node) Column() int {
	return node.column
}
func (node *Node) Position() (int, int) {
	return node.line, node.column
}
func (node *Node) SetPosition(line int, column int) {
	node.line = line
	node.column = column
}

type File struct {
	Node

	definitions []*Definition
}

type Definition struct {
	Node

	InferType bool

	Type  IExpr
	Value IExpr
}

type IExpr interface {
	INode
	exprTag()
}

type Expr struct {
	Node
}

func (expr *Expr) exprTag() {}

type ErrorExpr struct {
	Expr
}

type BinaryExpr struct {
	Expr

	Operator scanner.Token
	Lhs      IExpr
	Rhs      IExpr
}

type UnaryExpr struct {
	Expr

	Operator scanner.Token
	SubExpr  IExpr
}

type IntegerLiteral struct {
	Expr

	Value uint64
}

type StringLiteral struct {
	Expr

	Text string
}

type FloatLiteral struct {
	Expr

	Value float64
}

type IdentExpr struct {
	Expr

	// Namespaces followed by final identifier
	Names []string
}

type FunctionCallExpr struct {
	Expr

	Function IExpr
	Args     CallArgs
}

// Represents argument passed to function call,
// name will be non-empty if it is keyword arg
type CallArg struct {
	Node

	Name  string
	Value IExpr
}

type CallArgs struct {
	Node

	ArgList []CallArg
}
