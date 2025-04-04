package codescout

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/galactixx/codescout/internal/pkgutils"
)

func inspectorGetNode[T any](inspector inspector[T], symbol string) (*T, error) {
	if len(inspector.getNodes()) == 0 {
		errMsg := fmt.Sprintf("no %s was found based on configuration", symbol)
		err := errors.New(errMsg)
		return nil, err
	}
	return &(inspector.getNodes())[0], nil
}

type inspector[T any] interface {
	isNodeMatch(name T) bool
	appendNode(node T)
	inspector(n ast.Node) bool
	inspect()
	getNodes() []T
}

type baseInspector struct {
	Path  string
	Fset  *token.FileSet
	Count int
}

func (i baseInspector) inspect(node ast.Node, inspector func(n ast.Node) bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		if i.Count > 0 {
			return false
		}
		return inspector(n)
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

func (i *baseInspector) increment() {
	i.Count += 1
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
	Nodes  []StructNode
	Config StructConfig
	Base   baseInspector
}

func (i structInspector) isNodeMatch(node StructNode) bool {
	return true
}

func (i *structInspector) appendNode(node StructNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *structInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, i.inspector)
}

func (i *structInspector) getNodes() []StructNode {
	return i.Nodes
}

func (i structInspector) newStruct(node ast.Node, spec *ast.TypeSpec) StructNode {
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
				structNode := i.newStruct(structType, typeSpec)
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
	validParams := fullParamsMatch(i.Config.Types, i.Config.NoParams, node.CallableOps)
	validReceiver := !(i.Config.Receiver != "" && i.Config.Receiver != node.ReceiverType())

	validPtr := i.Config.IsPointerRec == nil || *i.Config.IsPointerRec == node.HasPointerReceiver()
	return nameEquals && validReturns && validParams && validReceiver && validPtr
}

func (i methodInspector) isAttrsMatch(node MethodNode) bool {
	accessed := fullAccessedMatch(i.Config.Fields, i.Config.NoFields, node)
	called := fullCalledMatch(i.Config.Methods, i.Config.NoMethods, node)
	return accessed && called
}

func (i *methodInspector) appendNode(node MethodNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *methodInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, i.inspector)
}

func (i methodInspector) getNodes() []MethodNode {
	return i.Nodes
}

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
	validParams := fullParamsMatch(i.Config.Types, i.Config.NoParams, node.CallableOps)
	return nameEquals && validReturns && validParams
}

func (i *funcInspector) appendNode(node FuncNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *funcInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, i.inspector)
}

func (i funcInspector) getNodes() []FuncNode {
	return i.Nodes
}

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
