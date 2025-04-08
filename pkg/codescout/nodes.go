package codescout

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/galactixx/codescout/internal/pkgutils"
)

// fieldListToNamedTypes converts a list of AST fields to a slice of NamedType structs,
// each containing the name and type of the field.
func fieldListToNamedTypes(fields *ast.FieldList, fset *token.FileSet) []NamedType {
	fieldList := make([]NamedType, 0)
	if fields == nil {
		return fieldList
	}

	for _, field := range fields.List {
		for _, name := range field.Names {
			named := NamedType{Name: name.Name, Type: pkgutils.NodeToCode(fset, field.Type)}
			fieldList = append(fieldList, named)
		}
	}
	return fieldList
}

// NodeInfo provides a generic interface for inspecting code entities.
type NodeInfo interface {
	Code() string
	PrintNode()
	PrintComments()
	Name() string
}

// BaseNode contains shared metadata about code elements like structs, methods, and functions.
type BaseNode struct {
	// Name of the code element (e.g., function or struct name)
	Name string
	// Path to the file where the element is defined
	Path string
	// Line number where the element starts
	Line int
	// Number of characters from the start of the file to the element
	Characters int
	// Whether the element is exported (starts with uppercase letter)
	Exported bool
	// Leading comment associated with the element
	Comment string
}

// StructNode represents a Go struct declaration in the AST.
type StructNode struct {
	// Node contains metadata such as name, path, line number, etc.
	Node BaseNode
	// Methods holds all methods associated with this struct
	Methods []*MethodNode

	node    *ast.StructType
	spec    *ast.TypeSpec
	genNode *ast.GenDecl
	fset    *token.FileSet
}

// Code returns the source code representation of the struct declaration.
func (s StructNode) Code() string { return pkgutils.NodeToCode(s.fset, s.genNode) }

// PrintNode prints the full code of the struct.
func (s StructNode) PrintNode() { fmt.Println(s.Code()) }

// PrintComments prints comments associated with the struct.
func (s StructNode) PrintComments() { fmt.Println(s.Comments()) }

// Name returns the name of the struct.
func (s StructNode) Name() string { return s.Node.Name }

// Fields extracts all named fields from the struct definition.
func (s StructNode) Fields() []NamedType { return fieldListToNamedTypes(s.node.Fields, s.fset) }

// Comments returns documentation comments associated with the struct declaration.
func (s StructNode) Comments() string { return pkgutils.CommentGroupToString(s.genNode.Doc) }

// Body returns the string representation of the struct's fields only.
func (s StructNode) Body() string {
	structFields := pkgutils.NodeToCode(s.fset, s.node)
	structFields = strings.Replace(structFields, "struct", "", 1)
	structFields = strings.TrimSpace(structFields)
	return structFields
}

// Signature returns the struct name along with any generic type parameters.
func (s StructNode) Signature() string {
	signature := s.Node.Name
	if s.spec.TypeParams == nil {
		return signature
	}

	var params []string
	for _, field := range s.spec.TypeParams.List {
		for _, name := range field.Names {
			params = append(params, name.Name)
		}
	}
	signature += "[" + strings.Join(params, ", ") + "]"
	return signature
}

// MethodNode represents a method with its associated metadata and interactions.
type MethodNode struct {
	// Node contains metadata such as name, path, line number, etc.
	Node BaseNode
	// CallableOps provides operations and data tied to the method's AST node
	CallableOps CallableOps

	fieldsAccessed map[string]*int
	methodsCalled  map[string]*int
}

// addMethodField registers that a struct field is accessed in this method.
func (m *MethodNode) addMethodField(field string) {
	if _, seenField := m.fieldsAccessed[field]; !seenField {
		m.fieldsAccessed[field] = nil
	}
}

// addMethodCall registers that another method is called from this method.
func (m *MethodNode) addMethodCall(method string) {
	if _, seenMethod := m.methodsCalled[method]; !seenMethod {
		m.methodsCalled[method] = nil
	}
}

// HasPointerReceiver checks whether the method has a pointer receiver.
func (m MethodNode) HasPointerReceiver() bool {
	if pkgutils.MethodWithoutReceiver(m.CallableOps.node) {
		return false
	}
	_, isPointer := m.CallableOps.node.Recv.List[0].Type.(*ast.StarExpr)
	return isPointer
}

// FieldsAccessed returns a slice of field names accessed by this method.
func (m MethodNode) FieldsAccessed() []string {
	return pkgutils.FromEmptyMapKeysToSlice(m.fieldsAccessed)
}

// MethodsCalled returns a slice of method names called by this method.
func (m MethodNode) MethodsCalled() []string {
	return pkgutils.FromEmptyMapKeysToSlice(m.methodsCalled)
}

// Code returns the full source code of the method.
func (m MethodNode) Code() string { return m.CallableOps.Code() }

// PrintNode prints the method's code.
func (m MethodNode) PrintNode() { fmt.Println(m.Code()) }

// PrintComments prints the method's associated documentation comments.
func (m MethodNode) PrintComments() { fmt.Println(m.CallableOps.Comments()) }

// Name returns the method name.
func (m MethodNode) Name() string { return m.Node.Name }

// ReceiverType returns the type name of the method's receiver.
func (m MethodNode) ReceiverType() string {
	if pkgutils.MethodWithoutReceiver(m.CallableOps.node) {
		return ""
	}

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

// ReceiverName returns the name of the receiver variable (e.g., "m" in "func (m *MyStruct) ...").
func (m MethodNode) ReceiverName() string {
	if pkgutils.MethodWithoutRecvList(m.CallableOps.node) {
		return ""
	}

	receiverName := m.CallableOps.node.Recv.List[0]
	if len(receiverName.Names) > 0 && receiverName.Names[0] != nil {
		return receiverName.Names[0].Name
	}
	return ""
}

// FuncNode represents a top-level function with its metadata and operations.
type FuncNode struct {
	// Node contains metadata such as name, path, line number, etc.
	Node BaseNode
	// CallableOps provides operations and data tied to the function's AST node
	CallableOps CallableOps
}

// Code returns the full source code of the function.
func (f FuncNode) Code() string { return f.CallableOps.Code() }

// PrintNode prints the function's code.
func (f FuncNode) PrintNode() { fmt.Println(f.Code()) }

// PrintComments prints the function's associated comments.
func (f FuncNode) PrintComments() { fmt.Println(f.CallableOps.Comments()) }

// Name returns the function name.
func (f FuncNode) Name() string { return f.Node.Name }

// CallableOps contains logic for extracting code and metadata from AST function declarations.
type CallableOps struct {
	node *ast.FuncDecl

	fset *token.FileSet
}

// PrintReturnType prints the function's return type(s).
func (c CallableOps) PrintReturnType() { fmt.Println(c.ReturnType()) }

// PrintBody prints the body of the function.
func (c CallableOps) PrintBody() { fmt.Println(c.Body()) }

// PrintSignature prints the function's signature (name + parameters).
func (c CallableOps) PrintSignature() { fmt.Println(c.Signature()) }

// Comments returns the associated documentation comments for the function.
func (c CallableOps) Comments() string { return pkgutils.CommentGroupToString(c.node.Doc) }

// Body returns the string representation of the function body.
func (c CallableOps) Body() string { return pkgutils.NodeToCode(c.fset, c.node.Body) }

// Parameters returns a slice of parameter names and types for the function.
func (c CallableOps) Parameters() []NamedType {
	if c.node.Type == nil {
		return make([]NamedType, 0)
	}
	return fieldListToNamedTypes(c.node.Type.Params, c.fset)
}

// Code returns the full source code of the function, optionally including comments.
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

// Signature returns the function's name and parameter list.
func (c CallableOps) Signature() string {
	return pkgutils.NodeToCode(c.fset, &ast.FuncDecl{
		Name: c.node.Name,
		Type: c.node.Type,
	})
}

// ReturnType returns a string representing the function's return type(s).
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

// ReturnTypes returns a slice of string representations of all return types.
func (c CallableOps) ReturnTypes() []string {
	returnTypes := make([]string, 0, 5)
	if c.node.Type != nil && c.node.Type.Results != nil {
		for _, returnType := range c.node.Type.Results.List {
			if returnType.Type != nil {
				returnTypes = append(returnTypes, pkgutils.NodeToCode(c.fset, returnType.Type))
			}
		}
	}
	return returnTypes
}
