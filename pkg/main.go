package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

type Inspector interface {
	nameEquals(name string) bool
	appendNode(node BaseNode)
}

type BaseInspector struct {
	Path  string
	Fset  *token.FileSet
	Count int
}

func (i *BaseInspector) increment() {
	i.Count += 1
}

func (i BaseInspector) getPos(node ast.Node) (int, int) {
	pos := i.Fset.Position(node.Pos())
	line := pos.Line
	characters := pos.Column
	return line, characters
}

type StructInspector struct {
	Nodes  *[]StructNode
	Config StructConfig
	Base   BaseInspector
}

func (i StructInspector) nameEquals(name string) bool {
	return i.Config.Name == name
}

func (i *StructInspector) appendNode(node StructNode) {
	*i.Nodes = append(*i.Nodes, node)
	i.Base.increment()
}

type FuncInspector struct {
	Nodes  *[]FuncNode
	Config FunctionConfig
	Base   BaseInspector
}

func (i FuncInspector) nameEquals(name string) bool {
	return i.Config.Name == name
}

func (i *FuncInspector) appendNode(node FuncNode) {
	*i.Nodes = append(*i.Nodes, node)
	i.Base.increment()
}

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

type FuncNode struct {
	Node BaseNode
	node *ast.FuncDecl
}

func (f FuncNode) Comments() string {
	var buf bytes.Buffer
	for _, comment := range f.node.Doc.List {
		buf.WriteString(comment.Text + "\n")
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

func (f FuncNode) Body() string {
	return nodeToCode(f.node.Body)
}

func (f FuncNode) Signature() string {
	return nodeToCode(&ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
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

func nodeToCode(node any) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), node)
	return buf.String()
}

func newFunction(name string, node ast.Node, insp FuncInspector, comment string) FuncNode {
	baseNode := newNode(name, node, insp.Base, comment)
	funcNode := node.(*ast.FuncDecl)
	return FuncNode{Node: baseNode, node: funcNode}
}

func newStruct(node ast.Node, spec *ast.TypeSpec, insp StructInspector) StructNode {
	var comment string = ""
	if spec.Doc != nil {
		comment = spec.Doc.Text()
	}
	baseNode := newNode(spec.Name.Name, node, insp.Base, comment)
	structNode := node.(*ast.StructType)
	return StructNode{Node: baseNode, node: structNode}
}

func newNode(
	name string, node ast.Node, inspector BaseInspector, comment string,
) BaseNode {
	line, characters := inspector.getPos(node)
	return BaseNode{
		Name:       name,
		Path:       inspector.Path,
		Line:       line,
		Characters: characters,
		Exported:   token.IsExported(name),
		Comment:    comment,
	}
}

func functionInspector(n ast.Node, inspector FuncInspector) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}

	if inspector.nameEquals(funcDecl.Name.Name) {
		var comment string = ""
		if funcDecl.Doc != nil {
			comment = funcDecl.Doc.Text()
		}
		name := funcDecl.Name.Name
		funcNode := newFunction(name, funcDecl, inspector, comment)
		inspector.appendNode(funcNode)
	}
	return true
}

func structInspector(n ast.Node, inspector StructInspector) bool {
	genDecl, ok := n.(*ast.GenDecl)
	if !ok {
		return true
	}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				if inspector.nameEquals(typeSpec.Name.Name) {
					structNode := newStruct(structType, typeSpec, inspector)
					inspector.appendNode(structNode)
				}
			}
		}
	}
	return true
}

func methodInsepector(n ast.Node, inspector Inspector) bool {
	return true
}

func parseFile(src string, fset *token.FileSet) *ast.File {
	var node *ast.File
	_, err := os.Stat(src)
	if err == nil {
		node, _ = parser.ParseFile(fset, src, nil, parser.ParseComments)
	} else {
		node, _ = parser.ParseFile(fset, "", src, parser.ParseComments)
	}
	return node
}

func filePathExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		err := errors.New("an existing file path must be passed")
		return err
	}
	return nil
}

type FunctionConfig struct {
	Name string
}

type StructConfig struct {
	Name string
}

func ScoutFunction(path string, config FunctionConfig) (*FuncNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	fset := token.NewFileSet()
	inspector := FuncInspector{
		Nodes:  &[]FuncNode{},
		Config: config,
		Base:   BaseInspector{Path: path, Fset: fset},
	}
	node := parseFile(path, fset)
	ast.Inspect(node, func(n ast.Node) bool {
		if inspector.Base.Count > 0 {
			return false
		}
		return functionInspector(n, inspector)
	})

	if len(*inspector.Nodes) == 0 {
		err := errors.New("no functions were found based on configuration")
		return nil, err
	}
	return &(*inspector.Nodes)[0], nil
}

func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	if fileExistsErr := filePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	fset := token.NewFileSet()
	inspector := StructInspector{
		Nodes:  &[]StructNode{},
		Config: config,
		Base:   BaseInspector{Path: path, Fset: fset},
	}
	node := parseFile(path, fset)
	ast.Inspect(node, func(n ast.Node) bool {
		if inspector.Base.Count > 0 {
			return false
		}
		return structInspector(n, inspector)
	})

	if len(*inspector.Nodes) == 0 {
		err := errors.New("no functions were found based on configuration")
		return nil, err
	}
	return &(*inspector.Nodes)[0], nil
}
