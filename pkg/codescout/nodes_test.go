package codescout

import (
	"fmt"
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldListToNamedTypes(t *testing.T) {
	fset := token.NewFileSet()
	fields := &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{{Name: "foo"}},
				Type:  &ast.Ident{Name: "int"},
			},
			{
				Names: []*ast.Ident{{Name: "bar"}},
				Type:  &ast.Ident{Name: "string"},
			},
		},
	}
	result := fieldListToNamedTypes(fields, fset)

	assert.Len(t, result, 2)
	assert.Equal(t, "foo", result[0].Name)
	assert.Equal(t, "int", result[0].Type)
	assert.Equal(t, "bar", result[1].Name)
	assert.Equal(t, "string", result[1].Type)
}

func TestCallableOps(t *testing.T) {
	fset := token.NewFileSet()
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("MyFunc"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("x")},
						Type:  &ast.Ident{Name: "int"},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("y")},
						Type:  &ast.Ident{Name: "string"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "int"}},
					{Type: &ast.Ident{Name: "string"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BasicLit{
							Kind:  token.INT,
							Value: "0",
						},
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"default"`,
						},
					},
				},
			},
		},
	}
	c := CallableOps{node: funcDecl, fset: fset}
	assert.Equal(t, []string{"int", "string"}, c.ReturnTypes())
	assert.Equal(t, "func MyFunc(x int, y string) (int, string)", c.Signature())
	assert.Equal(t, []NamedType{{Name: "x", Type: "int"}, {Name: "y", Type: "string"}}, c.Parameters())
	assert.Equal(t, "func MyFunc(x int, y string) (int, string) {\n    return 0, \"default\"\n}", c.Code())
	assert.Equal(t, "(int, string)", c.ReturnType())
	assert.Equal(t, "", c.Comments())
}

func TestFuncNode(t *testing.T) {
	fset := token.NewFileSet()
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("AddAndLabel"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("a")},
						Type:  &ast.Ident{Name: "int"},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("b")},
						Type:  &ast.Ident{Name: "int"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "int"}},
					{Type: &ast.Ident{Name: "string"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BinaryExpr{
							X:  ast.NewIdent("a"),
							Op: token.ADD,
							Y:  ast.NewIdent("b"),
						},
						&ast.BasicLit{
							Kind:  token.STRING,
							Value: `"sum"`,
						},
					},
				},
			},
		},
	}

	base := BaseNode{
		Name:       "AddAndLabel",
		Path:       "somePath",
		Line:       1,
		Characters: 1,
		Exported:   true,
		Comment:    "",
	}

	funcNode := FuncNode{
		Node: base,
		CallableOps: CallableOps{
			node: funcDecl,
			fset: fset,
		},
	}

	assert.Equal(t, "AddAndLabel", funcNode.Name())
	assert.Equal(t, "func AddAndLabel(a int, b int) (int, string) {\n    return a + b, \"sum\"\n}", funcNode.Code())
	assert.Equal(t, []string{"int", "string"}, funcNode.CallableOps.ReturnTypes())
	assert.Equal(t, "func AddAndLabel(a int, b int) (int, string)", funcNode.CallableOps.Signature())
	assert.Equal(t, []NamedType{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}}, funcNode.CallableOps.Parameters())
	assert.Equal(t, "(int, string)", funcNode.CallableOps.ReturnType())
	assert.Equal(t, "", funcNode.CallableOps.Comments())
}

func TestMethodNode(t *testing.T) {
	fset := token.NewFileSet()
	funcDecl := &ast.FuncDecl{
		Name: ast.NewIdent("Greet"),
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent("g")},
					Type: &ast.StarExpr{
						X: ast.NewIdent("Greeter"),
					},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("name")},
						Type:  &ast.Ident{Name: "string"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: &ast.Ident{Name: "string"}},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.BinaryExpr{
							X:  &ast.BasicLit{Kind: token.STRING, Value: `"Hello, "`},
							Op: token.ADD,
							Y: &ast.SelectorExpr{
								X:   ast.NewIdent("g"),
								Sel: ast.NewIdent("Name"),
							},
						},
					},
				},
			},
		},
	}

	base := BaseNode{
		Name:       "Greet",
		Path:       "somePath",
		Line:       1,
		Characters: 1,
		Exported:   true,
		Comment:    "",
	}

	methodNode := MethodNode{
		Node: base,
		CallableOps: CallableOps{
			node: funcDecl,
			fset: fset,
		},
		fieldsAccessed: make(map[string]*int),
		methodsCalled:  make(map[string]*int),
	}

	assert.Equal(t, "Greet", methodNode.Name())
	assert.Equal(t, "func (g *Greeter) Greet(name string) string {\n    return \"Hello, \" + g.Name\n}", methodNode.Code())
	assert.True(t, methodNode.HasPointerReceiver())
	assert.Equal(t, []string{}, methodNode.FieldsAccessed())
	assert.Equal(t, []string{}, methodNode.MethodsCalled())
	assert.Equal(t, "Greeter", methodNode.ReceiverType())
	assert.Equal(t, "g", methodNode.ReceiverName())
}

func TestStructNode(t *testing.T) {
	fset := token.NewFileSet()
	fields := &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("ID")},
				Type:  &ast.Ident{Name: "int"},
			},
			{
				Names: []*ast.Ident{ast.NewIdent("Name")},
				Type:  &ast.Ident{Name: "string"},
			},
			{
				Names: []*ast.Ident{ast.NewIdent("Email")},
				Type:  &ast.Ident{Name: "string"},
			},
		},
	}

	structType := &ast.StructType{Fields: fields}
	typeSpec := &ast.TypeSpec{Name: ast.NewIdent("User"), Type: structType}
	genDecl := &ast.GenDecl{Tok: token.TYPE, Specs: []ast.Spec{typeSpec}}

	base := BaseNode{
		Name:       "User",
		Path:       "somePath",
		Line:       1,
		Characters: 1,
		Exported:   true,
		Comment:    "",
	}

	structNode := StructNode{
		Node:    base,
		node:    structType,
		spec:    typeSpec,
		genNode: genDecl,
		fset:    fset,
	}

	fmt.Println(structNode.Code())

	assert.Equal(t, "User", structNode.Name())
	assert.Equal(t, "type User struct {\n    ID    int\n    Name  string\n    Email string\n}", structNode.Code())
	assert.Equal(t, "{\n    ID    int\n    Name  string\n    Email string\n}", structNode.Body())
	assert.Equal(t, "", structNode.Comments())
	assert.Equal(
		t,
		[]NamedType{{Name: "ID", Type: "int"}, {Name: "Name", Type: "string"}, {Name: "Email", Type: "string"}},
		structNode.Fields(),
	)
	assert.Equal(t, "User", structNode.Signature())
}
