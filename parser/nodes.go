package parser

import "pegasus/scanner"

type INode interface {
	nodeTag()

	Name() string
	Line() int
	Column() int
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

type File struct {
	Node

	definitions []Definition
}

type Definition struct {
	Node

	Type  *IExpr
	Value *IExpr
}

type IExpr interface {
	exprTag()
}

type Expr struct {
	Node
}

func (expr *Expr) exprTag() {}

type BinaryExpr struct {
	Expr

	Operator scanner.Token
	Lhs      *IExpr
	Rhs      *IExpr
}

type UnaryExpr struct {
	Expr

	Operator scanner.Token
	SubExpr  *IExpr
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
