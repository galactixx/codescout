package codescout

import (
	"bytes"
	"errors"
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
	var node *ast.File
	_, err := os.Stat(src)
	if err == nil {
		node, _ = parser.ParseFile(fset, src, nil, parser.ParseComments)
	} else {
		node, _ = parser.ParseFile(fset, "", src, parser.ParseComments)
	}
	return node
}

func filePathExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		err := errors.New("an existing file path must be passed")
		return err
	}
	return nil
}
