package parser

type INode interface {
	nodeTag()
}

type Node struct{}

func (node *Node) nodeTag() {}

type HasName struct {
	name string
}

func (hasName *HasName) Name() string {
	return hasName.name
}

type File struct {
	HasName
	Node

	definitions []*IDefinition
}

type IDefinition interface {
	definitionTag()
	Name() string
}

type Definition struct {
	Node
	HasName
}

func (*Definition) definitionTag() {}

type VarDefinition struct {
	Definition
}

type TypeDefinition struct {
	Definition
}

type FunctionDefinition struct {
	Definition
}
