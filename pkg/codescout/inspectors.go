package codescout

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/galactixx/codescout/internal/pkgutils"
)

type inspector[T any] interface {
	isNodeMatch(name T) bool
	appendNode(node T)
	inspector(n ast.Node) bool
	inspect()
	getNodes() []T
}

type baseInspector struct {
	Path string
	Fset *token.FileSet
}

func (i baseInspector) inspect(node ast.Node, inspectors []func(n ast.Node) bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		for _, inspector := range inspectors {
			inspector(n)
		}
		return true
	})
}

func (i baseInspector) newNode(name string, node ast.Node, comment string) BaseNode {
	line, characters := i.getPos(node)
	return BaseNode{
		Name:       name,
		Path:       i.Path,
		Line:       line,
		Characters: characters,
		Exported:   token.IsExported(name),
		Comment:    strings.TrimSpace(comment),
	}
}

func (i baseInspector) getCallableNodes(name string, node ast.Node, comment string) (BaseNode, *ast.FuncDecl) {
	baseNode := i.newNode(name, node, comment)
	funcNode := node.(*ast.FuncDecl)
	return baseNode, funcNode
}

func (i baseInspector) getPos(node ast.Node) (int, int) {
	pos := i.Fset.Position(node.Pos())
	line := pos.Line
	characters := pos.Column
	return line, characters
}

type structInspector struct {
	Nodes  map[string]*StructNode
	Config StructConfig
	Base   baseInspector
}

func (i structInspector) isNodeMatch(node StructNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	validFields := fullTypesMatch(i.Config.FieldTypes, i.Config.NoFields, node.Fields())
	return nameEquals && validFields
}

func (i *structInspector) appendNode(node StructNode) {
	i.Nodes[node.Node.Name] = &node
}

func (i *structInspector) inspect() {
	methodsInspect := methodInspector{
		Nodes:  []MethodNode{},
		Config: MethodConfig{},
		Base:   i.Base,
	}

	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, []func(n ast.Node) bool{i.inspector, methodsInspect.inspector})

	for _, methodNode := range methodsInspect.Nodes {
		if structNode, ok := i.Nodes[methodNode.ReceiverType()]; ok {
			structNode.Methods = append(structNode.Methods, methodNode)
		}
	}
}

func (i *structInspector) getNodes() []StructNode {
	structNodes := make([]StructNode, 0, len(i.Nodes))
	for _, node := range i.Nodes {
		structNodes = append(structNodes, *node)
	}
	return structNodes
}

func (i structInspector) newStruct(node ast.Node, gen *ast.GenDecl, spec *ast.TypeSpec) StructNode {
	var comment string = ""
	if gen.Doc != nil {
		comment = gen.Doc.Text()
	}
	baseNode := i.Base.newNode(spec.Name.Name, node, comment)
	structNode := node.(*ast.StructType)
	return StructNode{
		Node: baseNode, node: structNode, spec: spec, genNode: gen, fset: i.Base.Fset,
	}
}

func (i *structInspector) inspector(node ast.Node) bool {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return true
	}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				structNode := i.newStruct(structType, genDecl, typeSpec)
				if i.isNodeMatch(structNode) {
					i.appendNode(structNode)
				}
			}
		}
	}
	return true
}

type methodInspector struct {
	Nodes  []MethodNode
	Config MethodConfig
	Base   baseInspector
}

func (i methodInspector) isNodeMatch(node MethodNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	validReturns := fullReturnMatch(i.Config.ReturnTypes, i.Config.NoReturn, node.CallableOps)
	validParams := fullTypesMatch(i.Config.ParamTypes, i.Config.NoParams, node.CallableOps.Parameters())
	validReceiver := !(i.Config.Receiver != "" && i.Config.Receiver != node.ReceiverType())

	validPtr := i.Config.IsPointerRec == nil || *i.Config.IsPointerRec == node.HasPointerReceiver()
	return nameEquals && validReturns && validParams && validReceiver && validPtr
}

func (i methodInspector) isAttrsMatch(node MethodNode) bool {
	accessed := fullAccessedMatch(i.Config.Fields, i.Config.NoFields, node)
	called := fullCalledMatch(i.Config.Methods, i.Config.NoMethods, node)
	return accessed && called
}

func (i *methodInspector) appendNode(node MethodNode) { i.Nodes = append(i.Nodes, node) }
func (i *methodInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, []func(n ast.Node) bool{i.inspector})
}

func (i methodInspector) getNodes() []MethodNode { return i.Nodes }
func (i methodInspector) newMethod(name string, node ast.Node, comment string) MethodNode {
	baseNode, funcNode := i.Base.getCallableNodes(name, node, comment)
	return MethodNode{
		Node:           baseNode,
		CallableOps:    CallableOps{node: funcNode, fset: i.Base.Fset},
		fieldsAccessed: make(map[string]*int),
		methodsCalled:  make(map[string]*int),
	}
}

func (i *methodInspector) inspector(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok || funcDecl.Recv == nil {
		return true
	}

	var comment string = ""
	if funcDecl.Doc != nil {
		comment = funcDecl.Doc.Text()
	}
	name := funcDecl.Name.Name
	methodNode := i.newMethod(name, funcDecl, comment)
	receiverName := methodNode.ReceiverName()

	if i.isNodeMatch(methodNode) {
		var parentStack []ast.Node
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			if n == nil {
				parentStack = parentStack[:len(parentStack)-1]
				return true
			}

			sel, ok := n.(*ast.SelectorExpr)
			if ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == receiverName {
					curParent := parentStack[len(parentStack)-1]
					if call, isCall := curParent.(*ast.CallExpr); isCall && call.Fun == sel {
						methodNode.addMethodCall(sel.Sel.Name)
					} else {
						methodNode.addMethodField(sel.Sel.Name)
					}
				}
			}
			parentStack = append(parentStack, n)
			return true
		})
		if i.isAttrsMatch(methodNode) {
			i.appendNode(methodNode)
		}
	}

	return true
}

type funcInspector struct {
	Nodes  []FuncNode
	Config FuncConfig
	Base   baseInspector
}

func (i funcInspector) isNodeMatch(node FuncNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	validReturns := fullReturnMatch(i.Config.ReturnTypes, i.Config.NoReturn, node.CallableOps)
	validParams := fullTypesMatch(i.Config.ParamTypes, i.Config.NoParams, node.CallableOps.Parameters())
	return nameEquals && validReturns && validParams
}

func (i *funcInspector) appendNode(node FuncNode) {
	i.Nodes = append(i.Nodes, node)
}

func (i *funcInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, []func(n ast.Node) bool{i.inspector})
}

func (i funcInspector) getNodes() []FuncNode { return i.Nodes }
func (i funcInspector) newFunction(name string, node ast.Node, comment string) FuncNode {
	baseNode, funcNode := i.Base.getCallableNodes(name, node, comment)
	return FuncNode{Node: baseNode, CallableOps: CallableOps{node: funcNode, fset: i.Base.Fset}}
}

func (i *funcInspector) inspector(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}

	if funcDecl.Recv != nil {
		return true
	}

	var comment string = ""
	if funcDecl.Doc != nil {
		comment = funcDecl.Doc.Text()
	}
	name := funcDecl.Name.Name
	funcNode := i.newFunction(name, funcDecl, comment)

	if i.isNodeMatch(funcNode) {
		i.appendNode(funcNode)
	}
	return true
}
