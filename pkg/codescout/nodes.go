package codescout

import (
	"bytes"
	"fmt"
	"go/ast"
)

type NodeInfo interface {
	Code() string
	PrintNode()
	PrintComments()
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

type MethodNode struct {
	Node        BaseNode
	CallableOps CallableNodeOps

	fieldsAccessed map[string]*int
	methodsCalled  map[string]*int
}

func (f *MethodNode) addMethodField(field string) {
	if _, seenField := f.fieldsAccessed[field]; !seenField {
		f.fieldsAccessed[field] = nil
	}
}

func (f *MethodNode) addMethodCall(method string) {
	if _, seenMethod := f.methodsCalled[method]; !seenMethod {
		f.methodsCalled[method] = nil
	}
}

func (f MethodNode) HasPointerReceiver() bool {
	_, isPointer := f.CallableOps.node.Recv.List[0].Type.(*ast.StarExpr)
	return isPointer
}

func (f MethodNode) FieldsAccessed() []string {
	return fromEmptyMapKeysToSlice(f.fieldsAccessed)
}

func (f MethodNode) MethodsCalled() []string {
	return fromEmptyMapKeysToSlice(f.methodsCalled)
}

func (f MethodNode) PrintNode() {
	fmt.Println(f.CallableOps.Code())
}

func (f MethodNode) PrintComments() {
	fmt.Println(f.CallableOps.Comments())
}

func (f MethodNode) ReceiverType() string {
	structType := f.CallableOps.node.Recv.List[0].Type
	switch expr := structType.(type) {
	case *ast.Ident:
		// e.g. (m MyStruct)
		return expr.Name
	case *ast.StarExpr:
		switch selExpr := expr.X.(type) {
		case *ast.SelectorExpr:
			// e.g. (m *pkg.MyStruct)
			return formatStructName(selExpr)
		case *ast.Ident:
			// e.g. (m *MyStruct)
			return selExpr.Name
		default:
			return ""
		}
	case *ast.SelectorExpr:
		// e.g. (m pkg.MyStruct)
		return formatStructName(expr)
	default:
		return ""
	}
}

func (f MethodNode) ReceiverName() string {
	if len(f.CallableOps.node.Recv.List[0].Names) > 0 {
		return f.CallableOps.node.Recv.List[0].Names[0].Name
	}
	return ""
}

type FuncNode struct {
	Node        BaseNode
	CallableOps CallableNodeOps
}

func (f FuncNode) PrintNode() {
	fmt.Println(f.CallableOps.Code())
}

func (f FuncNode) PrintComments() {
	fmt.Println(f.CallableOps.Comments())
}

type CallableNodeOps struct {
	node *ast.FuncDecl
}

func (f CallableNodeOps) PrintReturnType() {
	fmt.Println(f.ReturnType())
}

func (f CallableNodeOps) PrintBody() {
	fmt.Println(f.Body())
}

func (f CallableNodeOps) PrintSignature() {
	fmt.Println(f.Signature())
}

func (f CallableNodeOps) parameterTypeMap() map[string]int {
	var parameterTypes []string
	for _, parameter := range f.Parameters() {
		parameterTypes = append(parameterTypes, parameter.Type)
	}
	return defaultTypeMap(parameterTypes)
}

func (f CallableNodeOps) ParametersMap() map[string]string {
	parameters := make(map[string]string)
	for _, parameter := range f.Parameters() {
		parameters[parameter.Name] = parameter.Type
	}
	return parameters
}

func (f CallableNodeOps) Parameters() []Parameter {
	parameters := make([]Parameter, 0, 5)
	for _, parameter := range f.node.Type.Params.List {
		for _, name := range parameter.Names {
			parameter := Parameter{Name: name.Name, Type: nodeToCode(parameter.Type)}
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}

func (f CallableNodeOps) Comments() string {
	var buf bytes.Buffer
	for _, comment := range f.node.Doc.List {
		buf.WriteString(comment.Text)
	}
	return buf.String()
}

func (f CallableNodeOps) Body() string {
	return nodeToCode(f.node.Body)
}

func (f CallableNodeOps) Code() string {
	nodeOriginalDoc := f.node.Doc
	f.node.Doc = nil
	codeString := nodeToCode(f.node)
	f.node.Doc = nodeOriginalDoc
	if f.node.Doc == nil {
		return codeString
	} else {
		return f.Comments() + "\n" + codeString
	}
}

func (f CallableNodeOps) Signature() string {
	return nodeToCode(&ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
}

func (f CallableNodeOps) returnTypesMap() map[string]int {
	return defaultTypeMap(f.ReturnTypes())
}

func (f CallableNodeOps) ReturnType() string {
	return nodeToCode(f.node.Type.Results)
}

func (f CallableNodeOps) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if f.node.Type.Results != nil {
		for _, returnType := range f.node.Type.Results.List {
			returnTypes = append(returnTypes, nodeToCode(returnType.Type))
		}
	}
	return returnTypes
}
