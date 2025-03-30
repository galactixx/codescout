package codescout

import (
	"errors"
	"go/token"
)

type Parameter struct {
	Name string
	Type string
}

type FunctionConfig struct {
	Name        string
	Types       []Parameter
	ReturnTypes []string
}

type StructConfig struct {
	Name string
}

func ScoutFunction(path string, config FunctionConfig) (*FuncNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	fset := token.NewFileSet()
	inspector := funcInspector{
		Nodes:  []FuncNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: fset},
	}
	node := parseFile(path, fset)
	inspector.inspect(node)

	if len(inspector.Nodes) == 0 {
		err := errors.New("no function was found based on configuration")
		return nil, err
	}
	return &(inspector.Nodes)[0], nil
}

func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	fset := token.NewFileSet()
	inspector := structInspector{
		Nodes:  []StructNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: fset},
	}
	node := parseFile(path, fset)
	inspector.inspect(node)

	if len(inspector.Nodes) == 0 {
		err := errors.New("no function was found based on configuration")
		return nil, err
	}
	return &(inspector.Nodes)[0], nil
}
