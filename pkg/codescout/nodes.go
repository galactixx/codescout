package codescout

import (
	"bytes"
	"fmt"
	"go/ast"

	"github.com/galactixx/codescout/internal/pkgutils"
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
	return pkgutils.NodeToCode(s.node)
}

func (s StructNode) PrintNode() {
	fmt.Println(s.Code())
}

type MethodNode struct {
	Node        BaseNode
	CallableOps CallableOps

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
	return pkgutils.FromEmptyMapKeysToSlice(f.fieldsAccessed)
}

func (f MethodNode) MethodsCalled() []string {
	return pkgutils.FromEmptyMapKeysToSlice(f.methodsCalled)
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
		// -> m MyStruct
		return expr.Name
	case *ast.StarExpr:
		switch selExpr := expr.X.(type) {
		case *ast.SelectorExpr:
			// -> m *pkg.MyStruct
			return pkgutils.FormatStructName(selExpr)
		case *ast.Ident:
			// -> m *MyStruct
			return selExpr.Name
		default:
			return ""
		}
	case *ast.SelectorExpr:
		// -> m pkg.MyStruct
		return pkgutils.FormatStructName(expr)
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
	CallableOps CallableOps
}

func (f FuncNode) PrintNode() {
	fmt.Println(f.CallableOps.Code())
}

func (f FuncNode) PrintComments() {
	fmt.Println(f.CallableOps.Comments())
}

type CallableOps struct {
	node *ast.FuncDecl
}

func (f CallableOps) PrintReturnType() {
	fmt.Println(f.ReturnType())
}

func (f CallableOps) PrintBody() {
	fmt.Println(f.Body())
}

func (f CallableOps) PrintSignature() {
	fmt.Println(f.Signature())
}

func (f CallableOps) parameterTypeMap() map[string]int {
	var parameterTypes []string
	for _, parameter := range f.Parameters() {
		parameterTypes = append(parameterTypes, parameter.Type)
	}
	return pkgutils.DefaultTypeMap(parameterTypes)
}

func (f CallableOps) ParametersMap() map[string]string {
	parameters := make(map[string]string)
	for _, parameter := range f.Parameters() {
		parameters[parameter.Name] = parameter.Type
	}
	return parameters
}

func (f CallableOps) Parameters() []Parameter {
	parameters := make([]Parameter, 0, 5)
	for _, parameter := range f.node.Type.Params.List {
		for _, name := range parameter.Names {
			parameter := Parameter{Name: name.Name, Type: pkgutils.NodeToCode(parameter.Type)}
			parameters = append(parameters, parameter)
		}
	}
	return parameters
}

func (f CallableOps) Comments() string {
	var buf bytes.Buffer
	for _, comment := range f.node.Doc.List {
		buf.WriteString(comment.Text)
	}
	return buf.String()
}

func (f CallableOps) Body() string {
	return pkgutils.NodeToCode(f.node.Body)
}

func (f CallableOps) Code() string {
	nodeOriginalDoc := f.node.Doc
	f.node.Doc = nil
	codeString := pkgutils.NodeToCode(f.node)
	f.node.Doc = nodeOriginalDoc
	if f.node.Doc == nil {
		return codeString
	} else {
		return f.Comments() + "\n" + codeString
	}
}

func (f CallableOps) Signature() string {
	return pkgutils.NodeToCode(&ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
}

func (f CallableOps) returnTypesMap() map[string]int {
	return pkgutils.DefaultTypeMap(f.ReturnTypes())
}

func (f CallableOps) ReturnType() string {
	return pkgutils.NodeToCode(f.node.Type.Results)
}

func (f CallableOps) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if f.node.Type.Results != nil {
		for _, returnType := range f.node.Type.Results.List {
			returnTypes = append(returnTypes, pkgutils.NodeToCode(returnType.Type))
		}
	}
	return returnTypes
}
