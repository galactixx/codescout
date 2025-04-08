package codescout

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/galactixx/codescout/internal/pkgutils"
)

// inspector is a generic interface for inspecting and collecting AST nodes of type T.
type inspector[T any] interface {
	isNodeMatch(name *T) bool
	appendNode(node *T)
	inspector(n ast.Node) bool
	inspect()
	getNodes() []*T
}

// baseInspector provides shared utilities for AST traversal and node metadata extraction.
type baseInspector struct {
	Path string
	Fset *token.FileSet
}

// inspect traverses the AST and applies a list of inspector functions to each node.
func (i baseInspector) inspect(node ast.Node, inspectors []func(n ast.Node) bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		for _, inspector := range inspectors {
			inspector(n)
		}
		return true
	})
}

// newNode constructs a BaseNode with name, position, export status, and comment.
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

// getCallableNodes extracts metadata and casts the node to *ast.FuncDecl.
func (i baseInspector) getCallableNodes(name string, node ast.Node, comment string) (BaseNode, *ast.FuncDecl) {
	baseNode := i.newNode(name, node, comment)
	funcNode := node.(*ast.FuncDecl)
	return baseNode, funcNode
}

// getPos returns the line and column position of the given AST node.
func (i baseInspector) getPos(node ast.Node) (int, int) {
	pos := i.Fset.Position(node.Pos())
	line := pos.Line
	characters := pos.Column
	return line, characters
}

// structInspector inspects struct declarations and associates their methods.
type structInspector struct {
	Nodes  map[string]*StructNode
	Config StructConfig
	Base   baseInspector
}

// isNodeMatch determines whether a StructNode matches the struct inspection configuration.
func (i structInspector) isNodeMatch(node *StructNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	matchFields := astMatch(i.Config.FieldTypes, node.Fields(), i.Config.Exact, i.Config.NoFields, namedTypesMatch)
	return nameEquals && matchFields.validate()
}

// appendNode stores a matched StructNode in the inspector's map.
func (i *structInspector) appendNode(node *StructNode) {
	i.Nodes[node.Node.Name] = node
}

// inspect performs the struct inspection and attaches discovered methods to their respective structs.
func (i *structInspector) inspect() {
	methodsInspect := methodInspector{
		Nodes:  []*MethodNode{},
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

// getNodes returns a slice of all matched StructNode instances.
func (i *structInspector) getNodes() []*StructNode {
	structNodes := make([]*StructNode, 0, len(i.Nodes))
	for _, node := range i.Nodes {
		structNodes = append(structNodes, node)
	}
	return structNodes
}

// newStruct constructs a StructNode from its AST components.
func (i structInspector) newStruct(node ast.Node, gen *ast.GenDecl, spec *ast.TypeSpec) *StructNode {
	var comment string = ""
	if gen.Doc != nil {
		comment = gen.Doc.Text()
	}
	baseNode := i.Base.newNode(spec.Name.Name, node, comment)
	structNode := node.(*ast.StructType)
	return &StructNode{
		Node: baseNode, node: structNode, spec: spec, genNode: gen, fset: i.Base.Fset,
	}
}

// inspector checks if the current AST node is a struct declaration, then stores it if matched.
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

// methodInspector inspects methods and captures metadata such as accessed fields and called methods.
type methodInspector struct {
	Nodes  []*MethodNode
	Config MethodConfig
	Base   baseInspector
}

// isNodeMatch determines whether a MethodNode matches method inspection criteria.
func (i methodInspector) isNodeMatch(node *MethodNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	matchReturn := astMatch(i.Config.ReturnTypes, node.CallableOps.ReturnTypes(), i.Config.Exact, i.Config.NoReturn, returnMatch)
	matchParams := astMatch(i.Config.ParamTypes, node.CallableOps.Parameters(), i.Config.Exact, i.Config.NoParams, namedTypesMatch)
	validReceiver := !(i.Config.Receiver != "" && i.Config.Receiver != node.ReceiverType())

	validPtr := i.Config.IsPointerRec == nil || *i.Config.IsPointerRec == node.HasPointerReceiver()
	return nameEquals && matchReturn.validate() && matchParams.validate() && validReceiver && validPtr
}

// isAttrsMatch validates the fields accessed and methods called by the method node.
func (i methodInspector) isAttrsMatch(node *MethodNode) bool {
	matchAccessed := astMatch(i.Config.Fields, node.FieldsAccessed(), i.Config.Exact, i.Config.NoFields, accessedMatch)
	matchCalled := astMatch(i.Config.Methods, node.MethodsCalled(), i.Config.Exact, i.Config.NoMethods, accessedMatch)
	return matchAccessed.validate() && matchCalled.validate()
}

// appendNode stores a matched MethodNode.
func (i *methodInspector) appendNode(node *MethodNode) { i.Nodes = append(i.Nodes, node) }

// inspect parses and traverses the file to extract method nodes.
func (i *methodInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, []func(n ast.Node) bool{i.inspector})
}

// getNodes returns all matched MethodNode instances.
func (i methodInspector) getNodes() []*MethodNode { return i.Nodes }

// newMethod constructs a MethodNode with tracking fields and associated callable operations.
func (i methodInspector) newMethod(name string, node ast.Node, comment string) *MethodNode {
	baseNode, funcNode := i.Base.getCallableNodes(name, node, comment)
	return &MethodNode{
		Node:           baseNode,
		CallableOps:    CallableOps{node: funcNode, fset: i.Base.Fset},
		fieldsAccessed: make(map[string]*int),
		methodsCalled:  make(map[string]*int),
	}
}

// inspector traverses AST to identify method declarations and captures method interactions.
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

// funcInspector inspects top-level functions in the Go source.
type funcInspector struct {
	Nodes  []*FuncNode
	Config FuncConfig
	Base   baseInspector
}

// isNodeMatch determines whether a function matches the criteria defined in FuncConfig.
func (i funcInspector) isNodeMatch(node *FuncNode) bool {
	nameEquals := !(i.Config.Name != "" && i.Config.Name != node.Node.Name)
	matchReturn := astMatch(i.Config.ReturnTypes, node.CallableOps.ReturnTypes(), i.Config.Exact, i.Config.NoReturn, returnMatch)
	matchParams := astMatch(i.Config.ParamTypes, node.CallableOps.Parameters(), i.Config.Exact, i.Config.NoParams, namedTypesMatch)
	return nameEquals && matchReturn.validate() && matchParams.validate()
}

// appendNode stores a matched FuncNode.
func (i *funcInspector) appendNode(node *FuncNode) {
	i.Nodes = append(i.Nodes, node)
}

// inspect parses and traverses the file to find matching function declarations.
func (i *funcInspector) inspect() {
	node := pkgutils.ParseFile(i.Base.Path, i.Base.Fset)
	i.Base.inspect(node, []func(n ast.Node) bool{i.inspector})
}

// getNodes returns all matched FuncNode instances.
func (i funcInspector) getNodes() []*FuncNode { return i.Nodes }

// newFunction constructs a FuncNode from the given AST node.
func (i funcInspector) newFunction(name string, node ast.Node, comment string) *FuncNode {
	baseNode, funcNode := i.Base.getCallableNodes(name, node, comment)
	return &FuncNode{Node: baseNode, CallableOps: CallableOps{node: funcNode, fset: i.Base.Fset}}
}

// inspector visits each AST node and extracts top-level (non-method) functions that match the config.
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
