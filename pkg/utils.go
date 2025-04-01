package codescout

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

func defaultTypeMap(nodeTypes []string) map[string]int {
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

func nodeToCode(node any) string {
	var buf bytes.Buffer
	printer.Fprint(&buf, token.NewFileSet(), node)
	return buf.String()
}

func parseFile(src string, fset *token.FileSet) *ast.File {
	node, _ := parser.ParseFile(fset, src, nil, parser.ParseComments)
	return node
}

func parseSource(src string, fset *token.FileSet) *ast.File {
	node, _ := parser.ParseFile(fset, "", src, parser.ParseComments)
	return node
}

func filePathExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		err := errors.New("an existing file path must be passed")
		return err
	}
	return nil
}

func inspectorGetNode[T any](inspector Inspector[T], symbol string) (*T, error) {
	if len(inspector.getNodes()) == 0 {
		errMsg := fmt.Sprintf("no %s was found based on configuration", symbol)
		err := errors.New(errMsg)
		return nil, err
	}
	return &(inspector.getNodes())[0], nil
}
