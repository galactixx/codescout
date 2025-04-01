package codescout

import (
	"bytes"
	"fmt"
	"go/ast"
)

type NodeInfo interface {
	Code() string
	PrintNode()
}

type BaseNode struct {
	Name       string
	Path       string
	Line       int
	Characters int
	Exported   bool
	Comment    string
}

type StructNode struct {
	Node BaseNode
	node *ast.StructType
}

func (s StructNode) Code() string {
	return nodeToCode(s.node)
}

func (s StructNode) PrintNode() {
	fmt.Println(s.Code())
}

type FuncNode struct {
	Node BaseNode
	node *ast.FuncDecl
}

func (f FuncNode) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if f.node.Type.Results != nil {
		for _, returnType := range f.node.Type.Results.List {
			returnTypes = append(returnTypes, nodeToCode(returnType.Type))
		}
	}
	return returnTypes
}

func (f FuncNode) ReturnType() string {
	return nodeToCode(f.node.Type.Results)
}

func (f FuncNode) returnTypesMap() map[string]int {
	return defaultTypeMap(f.ReturnTypes())
}

func (f FuncNode) Comments() string {
	var buf bytes.Buffer
	for _, comment := range f.node.Doc.List {
		buf.WriteString(comment.Text)
	}
	return buf.String()
}

func (f FuncNode) Code() string {
	nodeOriginalDoc := f.node.Doc
	f.node.Doc = nil
	codeString := nodeToCode(f.node)
	f.node.Doc = nodeOriginalDoc
	if f.node.Doc == nil {
		return codeString
	} else {
		return f.Comments() + codeString
	}
}

func (f FuncNode) PrintNode() {
	fmt.Println(f.Code())
}

func (f FuncNode) PrintBody() {
	fmt.Println(f.Body())
}

func (f FuncNode) PrintSignature() {
	fmt.Println(f.Signature())
}

func (f FuncNode) PrintReturnType() {
	fmt.Println(f.ReturnType())
}

func (f FuncNode) PrintComments() {
	fmt.Println(f.Comments())
}

func (f FuncNode) Body() string {
	return nodeToCode(f.node.Body)
}

func (f FuncNode) Signature() string {
	return nodeToCode(&ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
}

func (f FuncNode) parameterTypeMap() map[string]int {
	var parameterTypes []string
	for _, parameter := range f.Parameters() {
		parameterTypes = append(parameterTypes, parameter.Type)
	}
	return defaultTypeMap(parameterTypes)
}

func (f FuncNode) ParametersMap() map[string]string {
	parameters := make(map[string]string)
	for _, parameter := range f.Parameters() {
		parameters[parameter.Name] = parameter.Type
	}
	return parameters
}

func (f FuncNode) Parameters() []Parameter {
	parameters := make([]Parameter, 0, 5)
	for _, parameter := range f.node.Type.Params.List {
		for _, name := range parameter.Names {
			parameter := Parameter{Name: name.Name, Type: nodeToCode(parameter.Type)}
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}
