package codescout

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/galactixx/codescout/internal/pkgutils"
)

func fieldListToNamedTypes(fields ast.FieldList, fset *token.FileSet) []NamedType {
	fieldList := make([]NamedType, 0, len(fields.List))
	for _, field := range fields.List {
		for _, name := range field.Names {
			named := NamedType{Name: name.Name, Type: pkgutils.NodeToCode(fset, field.Type)}
			fieldList = append(fieldList, named)
		}
	}
	return fieldList
}

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
	Node    BaseNode
	Methods []MethodNode

	node    *ast.StructType
	spec    *ast.TypeSpec
	genNode *ast.GenDecl
	fset    *token.FileSet
}

func (s StructNode) Fields() []NamedType {
	return fieldListToNamedTypes(*s.node.Fields, s.fset)
}

func (s StructNode) Code() string {
	return pkgutils.NodeToCode(s.fset, s.genNode)
}

func (s StructNode) PrintNode() {
	fmt.Println(s.Code())
}

func (s StructNode) Body() string {
	structFields := pkgutils.NodeToCode(s.fset, s.node)
	structFields = strings.Replace(structFields, "struct", "", 1)
	structFields = strings.TrimSpace(structFields)
	return structFields
}

func (s StructNode) Signature() string {
	signature := s.Node.Name
	if s.spec.TypeParams != nil {
		var params []string
		for _, field := range s.spec.TypeParams.List {
			for _, name := range field.Names {
				params = append(params, name.Name)
			}
		}
		signature += "[" + strings.Join(params, ", ") + "]"
	}
	return signature
}

func (s StructNode) Comments() string {
	return pkgutils.CommentGroupToString(s.genNode.Doc)
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

	fset *token.FileSet
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

func (f CallableOps) Parameters() []NamedType {
	return fieldListToNamedTypes(*f.node.Type.Params, f.fset)
}

func (f CallableOps) Comments() string {
	return pkgutils.CommentGroupToString(f.node.Doc)
}

func (f CallableOps) Body() string {
	return pkgutils.NodeToCode(f.fset, f.node.Body)
}

func (f CallableOps) Code() string {
	nodeOriginalDoc := f.node.Doc
	f.node.Doc = nil
	codeString := pkgutils.NodeToCode(f.fset, f.node)
	f.node.Doc = nodeOriginalDoc
	if f.node.Doc == nil {
		return codeString
	} else {
		return f.Comments() + "\n" + codeString
	}
}

func (f CallableOps) Signature() string {
	return pkgutils.NodeToCode(f.fset, &ast.FuncDecl{
		Name: f.node.Name,
		Type: f.node.Type,
	})
}

func (f CallableOps) returnTypesMap() map[string]int {
	return pkgutils.DefaultTypeMap(f.ReturnTypes())
}

func (f CallableOps) ReturnType() string {
	nodeReturnTypes := f.ReturnTypes()
	switch len(nodeReturnTypes) {
	case 0:
		return ""
	case 1:
		return nodeReturnTypes[0]
	default:
		return "(" + strings.Join(nodeReturnTypes, ", ") + ")"
	}
}

func (f CallableOps) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if f.node.Type.Results != nil {
		for _, returnType := range f.node.Type.Results.List {
			returnTypes = append(returnTypes, pkgutils.NodeToCode(f.fset, returnType.Type))
		}
	}
	return returnTypes
}
