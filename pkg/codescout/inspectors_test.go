package codescout

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseInspectorNewNode(t *testing.T) {
	fset := token.NewFileSet()
	src := `
package main

// comment here
func X() {}
`
	node, _ := parser.ParseFile(fset, "", src, parser.ParseComments)
	inspector := baseInspector{Path: "mock/path.go", Fset: fset}
	baseNode := inspector.newNode("X", node.Decls[0], "comment here")

	assert.Equal(t, "X", baseNode.Name)
	assert.Equal(t, "mock/path.go", baseNode.Path)
	assert.True(t, baseNode.Exported)
	assert.Contains(t, baseNode.Comment, "comment")
	assert.Greater(t, baseNode.Line, 0)
}

func TestBaseInspectorGetCallableNodes(t *testing.T) {
	fset := token.NewFileSet()
	src := `package main; // doc
	func Do() error { return nil }`
	node, _ := parser.ParseFile(fset, "", src, parser.ParseComments)

	inspector := baseInspector{Path: "f.go", Fset: fset}
	baseNode, fn := inspector.getCallableNodes("Do", node.Decls[0], "comment")

	assert.Equal(t, "Do", baseNode.Name)
	assert.NotNil(t, fn)
}

func TestFuncInspectorAppendAndGetNodes(t *testing.T) {
	fset := token.NewFileSet()
	src := `package main; func Hello(name string) string { return "Hi " + name }`
	file, _ := parser.ParseFile(fset, "", src, 0)

	fi := funcInspector{Base: baseInspector{Path: "demo.go", Fset: fset}}
	fn := fi.newFunction("Hello", file.Decls[0], "")
	fi.appendNode(fn)

	nodes := fi.getNodes()
	assert.Len(t, nodes, 1)
	assert.Equal(t, "Hello", nodes[0].Node.Name)
	assert.Equal(t, `func Hello(name string) string { return "Hi " + name }`, nodes[0].Code())
}

func TestMethodInspectorReceiverAndAdd(t *testing.T) {
	fset := token.NewFileSet()
	src := `
package main
type Greeter struct{}
func (g *Greeter) Greet() string { return "Hello" }`
	file, _ := parser.ParseFile(fset, "", src, 0)

	mi := methodInspector{Base: baseInspector{Path: "greeter.go", Fset: fset}}
	method := mi.newMethod("Greet", file.Decls[1], "")
	assert.Equal(t, "Greeter", method.ReceiverType())
	assert.True(t, method.HasPointerReceiver())
}

func TestStructInspectorNewStruct(t *testing.T) {
	fset := token.NewFileSet()
	src := `package main
// Person doc
type Person struct {
	Name string
	Age int 
}`
	file, _ := parser.ParseFile(fset, "", src, parser.ParseComments)

	gen := file.Decls[0].(*ast.GenDecl)
	spec := gen.Specs[0].(*ast.TypeSpec)

	si := structInspector{Base: baseInspector{Path: "p.go", Fset: fset}}
	sNode := si.newStruct(spec.Type, gen, spec)

	assert.Equal(t, "Person", sNode.Name())
	assert.Equal(t, "// Person doc", sNode.Comments())
	assert.Equal(t, "{\n    Name string\n    Age  int\n}", sNode.Body())
}

func TestFuncInspectorInspectorAndMatch(t *testing.T) {
	fset := token.NewFileSet()
	src := `package main; func Hello(name string) string { return "Hi " + name }`
	file, _ := parser.ParseFile(fset, "", src, 0)

	funcConfig := FuncConfig{Name: "Hello", ReturnTypes: []string{"string"}}
	fi := funcInspector{Config: funcConfig, Base: baseInspector{Path: "demo.go", Fset: fset}}
	fn := fi.newFunction("Hello", file.Decls[0], "")

	_ = fi.inspector(file.Decls[0])
	assert.Len(t, fi.Nodes, 1)
	assert.True(t, fi.isNodeMatch(fn))
}

func TestMethodInspectorInspectorAndMatch(t *testing.T) {
	fset := token.NewFileSet()
	src := `
package main
type Greeter struct{}
func (g *Greeter) Greet() string { return "Hello" }`
	file, _ := parser.ParseFile(fset, "", src, 0)

	methodConfig := MethodConfig{Name: "Greet", ReturnTypes: []string{"string"}, Receiver: "Greeter"}
	mi := methodInspector{Config: methodConfig, Base: baseInspector{Path: "greeter.go", Fset: fset}}
	method := mi.newMethod("Greet", file.Decls[1], "")

	_ = mi.inspector(file.Decls[1])
	assert.Len(t, mi.Nodes, 1)
	assert.True(t, mi.isNodeMatch(method))
}

func TestStructInspectorInspectorAndMatch(t *testing.T) {
	fset := token.NewFileSet()
	src := `package main
// Person doc
type Person struct {
	Name string
	Age int 
}`
	file, _ := parser.ParseFile(fset, "", src, parser.ParseComments)

	gen := file.Decls[0].(*ast.GenDecl)
	spec := gen.Specs[0].(*ast.TypeSpec)

	noFields := false
	structConfig := StructConfig{Name: "Person", NoFields: &noFields}
	si := structInspector{Nodes: map[string]*StructNode{}, Config: structConfig, Base: baseInspector{Path: "p.go", Fset: fset}}
	sNode := si.newStruct(spec.Type, gen, spec)

	_ = si.inspector(file.Decls[0])
	assert.Len(t, si.Nodes, 1)
	assert.True(t, si.isNodeMatch(sNode))
}
