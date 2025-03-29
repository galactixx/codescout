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
	inspector(n ast.Node) bool
	inspect(n ast.Node)
}

type baseInspector struct {
	Path  string
	Fset  *token.FileSet
	Count int
}

func (i *baseInspector) inspect(node ast.Node, inspector func(n ast.Node) bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		if i.Count > 0 {
			return false
		}
		return inspector(n)
	})
}

func (i *baseInspector) newNode(name string, node ast.Node, comment string) BaseNode {
	line, characters := i.getPos(node)
	return BaseNode{
		Name:       name,
		Path:       i.Path,
		Line:       line,
		Characters: characters,
		Exported:   token.IsExported(name),
		Comment:    comment,
	}
}

func (i *baseInspector) increment() {
	i.Count += 1
}

func (i *baseInspector) getPos(node ast.Node) (int, int) {
	pos := i.Fset.Position(node.Pos())
	line := pos.Line
	characters := pos.Column
	return line, characters
}

type structInspector struct {
	Nodes  []StructNode
	Config StructConfig
	Base   baseInspector
}

func (i *structInspector) nameEquals(name string) bool {
	return i.Config.Name == name
}

func (i *structInspector) appendNode(node StructNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *structInspector) inspect(node ast.Node) {
	i.Base.inspect(node, i.inspector)
}

func (i *structInspector) newStruct(node ast.Node, spec *ast.TypeSpec) StructNode {
	var comment string = ""
	if spec.Doc != nil {
		comment = spec.Doc.Text()
	}
	baseNode := i.Base.newNode(spec.Name.Name, node, comment)
	structNode := node.(*ast.StructType)
	return StructNode{Node: baseNode, node: structNode}
}

func (i *structInspector) inspector(node ast.Node) bool {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return true
	}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				if i.nameEquals(typeSpec.Name.Name) {
					structNode := i.newStruct(structType, typeSpec)
					i.appendNode(structNode)
				}
			}
		}
	}
	return true
}

type funcInspector struct {
	Nodes  []FuncNode
	Config FunctionConfig
	Base   baseInspector
}

func (i *funcInspector) nameEquals(name string) bool {
	return i.Config.Name == name
}

func (i *funcInspector) appendNode(node FuncNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *funcInspector) inspect(node ast.Node) {
	i.Base.inspect(node, i.inspector)
}

func (i *funcInspector) newFunction(name string, node ast.Node, comment string) FuncNode {
	baseNode := i.Base.newNode(name, node, comment)
	funcNode := node.(*ast.FuncDecl)
	return FuncNode{Node: baseNode, node: funcNode}
}

func (i *funcInspector) inspector(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}

	if i.nameEquals(funcDecl.Name.Name) {
		var comment string = ""
		if funcDecl.Doc != nil {
			comment = funcDecl.Doc.Text()
		}
		name := funcDecl.Name.Name
		funcNode := i.newFunction(name, funcDecl, comment)
		i.appendNode(funcNode)
	}
	return true
}

type Parameter struct {
	Name string
	Type string
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

func (f *FuncNode) Parameters() []Parameter {
	parameters := make([]Parameter, 0, 5)
	for _, parameter := range f.node.Type.Params.List {
		for _, name := range parameter.Names {
			parameter := Parameter{Name: name.Name, Type: nodeToCode(parameter.Type)}
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}

func (f *FuncNode) Comments() string {
	var buf bytes.Buffer
	for _, comment := range f.node.Doc.List {
		buf.WriteString(comment.Text + "\n")
	}
	return buf.String()
}

func (f *FuncNode) Code() string {
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

func (f *FuncNode) PrintNode() {
	fmt.Println(f.Code())
}

func (f *FuncNode) Body() string {
	return nodeToCode(f.node.Body)
}

func (f *FuncNode) Signature() string {
	return nodeToCode(&ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
}

type StructNode struct {
	Node BaseNode
	node *ast.StructType
}

func (s *StructNode) Code() string {
	return nodeToCode(s.node)
}

func (s *StructNode) PrintNode() {
	fmt.Println(s.Code())
}

func nodeToCode(node any) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), node)
	return buf.String()
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
