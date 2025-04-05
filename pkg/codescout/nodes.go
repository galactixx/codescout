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
	Name() string
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

func (s StructNode) Code() string   { return pkgutils.NodeToCode(s.fset, s.genNode) }
func (s StructNode) PrintNode()     { fmt.Println(s.Code()) }
func (s StructNode) PrintComments() { fmt.Println(s.Comments()) }
func (s StructNode) Name() string   { return s.Node.Name }

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

func (s StructNode) Comments() string { return pkgutils.CommentGroupToString(s.genNode.Doc) }

type MethodNode struct {
	Node        BaseNode
	CallableOps CallableOps

	fieldsAccessed map[string]*int
	methodsCalled  map[string]*int
}

func (m *MethodNode) addMethodField(field string) {
	if _, seenField := m.fieldsAccessed[field]; !seenField {
		m.fieldsAccessed[field] = nil
	}
}

func (m *MethodNode) addMethodCall(method string) {
	if _, seenMethod := m.methodsCalled[method]; !seenMethod {
		m.methodsCalled[method] = nil
	}
}

func (m MethodNode) HasPointerReceiver() bool {
	_, isPointer := m.CallableOps.node.Recv.List[0].Type.(*ast.StarExpr)
	return isPointer
}

func (m MethodNode) FieldsAccessed() []string {
	return pkgutils.FromEmptyMapKeysToSlice(m.fieldsAccessed)
}

func (m MethodNode) MethodsCalled() []string {
	return pkgutils.FromEmptyMapKeysToSlice(m.methodsCalled)
}

func (m MethodNode) Code() string   { return m.CallableOps.Code() }
func (m MethodNode) PrintNode()     { fmt.Println(m.Code()) }
func (m MethodNode) PrintComments() { fmt.Println(m.CallableOps.Comments()) }
func (m MethodNode) Name() string   { return m.Node.Name }

func (m MethodNode) ReceiverType() string {
	structType := m.CallableOps.node.Recv.List[0].Type
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

func (m MethodNode) ReceiverName() string {
	if len(m.CallableOps.node.Recv.List[0].Names) > 0 {
		return m.CallableOps.node.Recv.List[0].Names[0].Name
	}
	return ""
}

type FuncNode struct {
	Node        BaseNode
	CallableOps CallableOps
}

func (f FuncNode) Code() string   { return f.CallableOps.Code() }
func (f FuncNode) PrintNode()     { fmt.Println(f.Code()) }
func (f FuncNode) PrintComments() { fmt.Println(f.CallableOps.Comments()) }
func (f FuncNode) Name() string   { return f.Node.Name }

type CallableOps struct {
	node *ast.FuncDecl

	fset *token.FileSet
}

func (c CallableOps) PrintReturnType() { fmt.Println(c.ReturnType()) }
func (c CallableOps) PrintBody()       { fmt.Println(c.Body()) }
func (c CallableOps) PrintSignature()  { fmt.Println(c.Signature()) }

func (c CallableOps) Parameters() []NamedType {
	return fieldListToNamedTypes(*c.node.Type.Params, c.fset)
}

func (c CallableOps) Comments() string { return pkgutils.CommentGroupToString(c.node.Doc) }
func (c CallableOps) Body() string     { return pkgutils.NodeToCode(c.fset, c.node.Body) }

func (c CallableOps) Code() string {
	nodeOriginalDoc := c.node.Doc
	c.node.Doc = nil
	codeString := pkgutils.NodeToCode(c.fset, c.node)
	c.node.Doc = nodeOriginalDoc
	if c.node.Doc == nil {
		return codeString
	} else {
		return c.Comments() + "\n" + codeString
	}
}

func (c CallableOps) Signature() string {
	return pkgutils.NodeToCode(c.fset, &ast.FuncDecl{
		Name: c.node.Name,
		Type: c.node.Type,
	})
}

func (c CallableOps) returnTypesMap() map[string]int {
	return pkgutils.DefaultTypeMap(c.ReturnTypes())
}

func (c CallableOps) ReturnType() string {
	nodeReturnTypes := c.ReturnTypes()
	switch len(nodeReturnTypes) {
	case 0:
		return ""
	case 1:
		return nodeReturnTypes[0]
	default:
		return "(" + strings.Join(nodeReturnTypes, ", ") + ")"
	}
}

func (c CallableOps) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if c.node.Type.Results != nil {
		for _, returnType := range c.node.Type.Results.List {
			returnTypes = append(returnTypes, pkgutils.NodeToCode(c.fset, returnType.Type))
		}
	}
	return returnTypes
}
