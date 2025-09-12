package parser

type File struct {
}

type Definition interface {
	definitionTag()
}

type VarDefinition struct {
}

func (VarDefinition) definitionTag() {}

type TypeDefinition struct {
}

func (TypeDefinition) definitionTag() {}

type FunctionDefinition struct {
}

func (FunctionDefinition) definitionTag() {}
