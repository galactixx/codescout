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

func formatStructName(expr *ast.SelectorExpr) string {
	return fmt.Sprintf("%s.%s", exprToString(expr.X), expr.Sel.Name)
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	default:
		return fmt.Sprintf("%T", e)
	}
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

func returnTypeValidation(returns []string, ops CallableNodeOps) bool {
	returnValidation := typeValidation{TypeMap: ops.returnTypesMap()}
	for _, returnType := range returns {
		returnValidation.setParamType(returnType)
		if !returnValidation.typeExclExists() ||
			returnValidation.hasExhausted(returnType) {
			return false
		}
	}
	return true
}

func parameterTypeValidation(params []Parameter, ops CallableNodeOps) bool {
	typesValidation := typeValidation{
		TypeMap:      ops.parameterTypeMap(),
		ParameterMap: ops.ParametersMap(),
	}
	var parameterType string

	for _, parameter := range params {
		typesValidation.setParamInfo(parameter.Name, parameter.Type)
		typesValidation.setNameInParams(parameter.Name)

		if !typesValidation.CurNameInParams && parameter.Name != "" ||
			!typesValidation.typeExists() {
			return false
		}

		if typesValidation.CurNameInParams {
			parameterType = typesValidation.getParamType(parameter.Name)
		} else {
			parameterType = parameter.Type
		}
		if typesValidation.hasExhausted(parameterType) {
			return false
		}
	}
	return true
}
