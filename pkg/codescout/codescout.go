package codescout

import (
	"go/token"
)

type Parameter struct {
	Name string
	Type string
}

type FuncConfig struct {
	Name        string
	Types       []Parameter
	ReturnTypes []string
}

type StructConfig struct {
	Name    string
	Types   []Parameter
	Methods []FuncConfig
}

func ScoutFunction(path string, config FuncConfig) (*FuncNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	inspector := funcInspector{
		Nodes:  []FuncNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: token.NewFileSet()},
	}
	inspector.inspect()
	return inspectorGetNode(&inspector, "function")
}

func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	inspector := structInspector{
		Nodes:  []StructNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: token.NewFileSet()},
	}
	inspector.inspect()
	return inspectorGetNode(&inspector, "struct")
}
