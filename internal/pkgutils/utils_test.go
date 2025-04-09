package pkgutils

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodWithoutRecvList(t *testing.T) {
	funcDecl := &ast.FuncDecl{
		Recv: &ast.FieldList{List: []*ast.Field{}},
	}
	assert.True(t, MethodWithoutRecvList(funcDecl))

	funcDecl.Recv.List = []*ast.Field{{Type: &ast.Ident{Name: "MyType"}}}
	assert.False(t, MethodWithoutRecvList(funcDecl))
}

func TestMethodWithoutReceiver(t *testing.T) {
	noRecv := &ast.FuncDecl{Recv: &ast.FieldList{List: []*ast.Field{}}}
	assert.True(t, MethodWithoutReceiver(noRecv))

	nilType := &ast.FuncDecl{Recv: &ast.FieldList{List: []*ast.Field{{Type: nil}}}}
	assert.True(t, MethodWithoutReceiver(nilType))

	hasType := &ast.FuncDecl{Recv: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "Receiver"}}}}}
	assert.False(t, MethodWithoutReceiver(hasType))
}

func TestFromEmptyMapKeysToSlice(t *testing.T) {
	input := map[string]*int{"a": nil, "b": nil}
	result := FromEmptyMapKeysToSlice(input)
	assert.ElementsMatch(t, []string{"a", "b"}, result)
}

func TestDefaultTypeMap(t *testing.T) {
	input := []string{"x", "y", "x"}
	expected := map[string]int{"x": 2, "y": 1}
	assert.Equal(t, expected, DefaultTypeMap(input))
}

func TestDefaultTypeNilMap(t *testing.T) {
	input := []string{"a", "b", "a"}
	result := DefaultTypeNilMap(input)
	assert.Contains(t, result, "a")
	assert.Contains(t, result, "b")
	assert.Nil(t, result["a"])
	assert.Nil(t, result["b"])
}

func TestFormatStructName(t *testing.T) {
	expr := &ast.SelectorExpr{
		X:   &ast.Ident{Name: "pkg"},
		Sel: &ast.Ident{Name: "Type"},
	}
	assert.Equal(t, "pkg.Type", FormatStructName(expr))
}

func TestExprToString(t *testing.T) {
	expr := &ast.Ident{Name: "someVar"}
	assert.Equal(t, "someVar", ExprToString(expr))

	unknown := &ast.BasicLit{}
	assert.Contains(t, ExprToString(unknown), "*ast.BasicLit")
}

func TestNodeToCode(t *testing.T) {
	fset := token.NewFileSet()
	expr := &ast.Ident{Name: "someVar"}
	code := NodeToCode(fset, expr)
	assert.Equal(t, "someVar", code)
}

func TestParseFile(t *testing.T) {
	fset := token.NewFileSet()

	// Create a temporary Go file
	content := `package main; func main() {}`
	tmpfile := filepath.Join(os.TempDir(), "tmp.go")
	err := os.WriteFile(tmpfile, []byte(content), 0o644)
	assert.NoError(t, err)
	defer os.Remove(tmpfile)

	file := ParseFile(tmpfile, fset)
	assert.Equal(t, "main", file.Name.Name)
}

func TestParseSource(t *testing.T) {
	fset := token.NewFileSet()
	src := `package demo; func Demo() {}`
	file := ParseSource(src, fset)
	assert.Equal(t, "demo", file.Name.Name)
}

func TestFilePathExists(t *testing.T) {
	// Existing file
	f, err := os.CreateTemp("", "testfile")
	assert.NoError(t, err)
	defer os.Remove(f.Name())

	assert.NoError(t, FilePathExists(f.Name()))

	// Non-existing file
	err = FilePathExists("non_existent_file.go")
	assert.Error(t, err)
}

func TestCommentGroupToString(t *testing.T) {
	commentGroup := &ast.CommentGroup{
		List: []*ast.Comment{
			{Text: "// This is a comment. "},
			{Text: "// Another comment."},
		},
	}
	result := CommentGroupToString(commentGroup)
	assert.Contains(t, result, "This is a comment.")
	assert.Contains(t, result, "Another comment.")

	assert.Equal(t, "", CommentGroupToString(nil))
}
