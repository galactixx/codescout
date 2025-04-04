package pkgutils

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func FromEmptyMapKeysToSlice(someMap map[string]*int) []string {
	fieldsAccessed := make([]string, 0, 10)
	for key := range someMap {
		fieldsAccessed = append(fieldsAccessed, key)
	}
	return fieldsAccessed
}

func DefaultTypeMap(nodeTypes []string) map[string]int {
	nodeTypeMapping := make(map[string]int)
	for _, nodeType := range nodeTypes {
		_, ok := nodeTypeMapping[nodeType]
		if !ok {
			nodeTypeMapping[nodeType] = 0
		}
		nodeTypeMapping[nodeType] += 1
	}
	return nodeTypeMapping
}

func FormatStructName(expr *ast.SelectorExpr) string {
	return fmt.Sprintf("%s.%s", ExprToString(expr.X), expr.Sel.Name)
}

func ExprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	default:
		return fmt.Sprintf("%T", e)
	}
}

func NodeToCode(fset *token.FileSet, node any) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, node)
	return buf.String()
}

func ParseFile(src string, fset *token.FileSet) *ast.File {
	node, _ := parser.ParseFile(fset, src, nil, parser.ParseComments)
	return node
}

func ParseSource(src string, fset *token.FileSet) *ast.File {
	node, _ := parser.ParseFile(fset, "", src, parser.ParseComments)
	return node
}

func FilePathExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		err := errors.New("an existing file path must be passed")
		return err
	}
	return nil
}

func CommentGroupToString(comment *ast.CommentGroup) string {
	if comment == nil {
		return ""
	}
	var buf bytes.Buffer
	for _, comment := range comment.List {
		buf.WriteString(comment.Text)
	}
	return buf.String()
}
