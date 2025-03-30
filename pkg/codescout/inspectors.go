package codescout

import (
	"go/ast"
	"go/token"
)

type Inspector interface {
	isNodeMatch(name string) bool
	appendNode(node BaseNode)
	inspector(n ast.Node) bool
	inspect(n ast.Node)
}

type baseInspector struct {
	Path  string
	Fset  *token.FileSet
	Count int
}

func (i *baseInspector) inspect(node ast.Node, inspector func(n ast.Node) bool) {
	ast.Inspect(node, func(n ast.Node) bool {
		if i.Count > 0 {
			return false
		}
		return inspector(n)
	})
}

func (i *baseInspector) newNode(name string, node ast.Node, comment string) BaseNode {
	line, characters := i.getPos(node)
	return BaseNode{
		Name:       name,
		Path:       i.Path,
		Line:       line,
		Characters: characters,
		Exported:   token.IsExported(name),
		Comment:    comment,
	}
}

func (i *baseInspector) increment() {
	i.Count += 1
}

func (i *baseInspector) getPos(node ast.Node) (int, int) {
	pos := i.Fset.Position(node.Pos())
	line := pos.Line
	characters := pos.Column
	return line, characters
}

type structInspector struct {
	Nodes  []StructNode
	Config StructConfig
	Base   baseInspector
}

func (i *structInspector) isNodeMatch(node StructNode) bool {
	return true
}

func (i *structInspector) appendNode(node StructNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *structInspector) inspect(node ast.Node) {
	i.Base.inspect(node, i.inspector)
}

func (i *structInspector) newStruct(node ast.Node, spec *ast.TypeSpec) StructNode {
	var comment string = ""
	if spec.Doc != nil {
		comment = spec.Doc.Text()
	}
	baseNode := i.Base.newNode(spec.Name.Name, node, comment)
	structNode := node.(*ast.StructType)
	return StructNode{Node: baseNode, node: structNode}
}

func (i *structInspector) inspector(node ast.Node) bool {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return true
	}

	for _, spec := range genDecl.Specs {
		if typeSpec, ok := spec.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				structNode := i.newStruct(structType, typeSpec)
				if i.isNodeMatch(structNode) {
					i.appendNode(structNode)
				}
			}
		}
	}
	return true
}

type funcInspector struct {
	Nodes  []FuncNode
	Config FunctionConfig
	Base   baseInspector
}

func (i *funcInspector) isNodeMatch(node FuncNode) bool {
	if i.Config.Name != "" && i.Config.Name != node.Node.Name {
		return false
	}

	returnTypesValidation := typeValidation{TypeMap: node.returnTypesMap()}
	for _, returnType := range i.Config.ReturnTypes {
		if returnTypesValidation.typeExclusiveNotInParams(returnType) {
			return false
		}
		if returnTypesValidation.hasExhaustedTypes(returnType) {
			return false
		}
	}

	typesValidation := typeValidation{
		TypeMap: node.parameterTypeMap(), ParameterMap: node.ParametersMap(),
	}
	var parameterType string

	for _, funcParameter := range i.Config.Types {
		paramType := funcParameter.Type
		name := funcParameter.Name
		typesValidation.setNameInParams(name)
		if !typesValidation.NameInParams {
			return false
		}

		if typesValidation.typeNotInParams(name, paramType) {
			return false
		}

		if typesValidation.typeExclusiveNotInParams(paramType) {
			return false
		}

		if name != "" {
			parameterType = typesValidation.ParameterMap[name]
		} else {
			parameterType = paramType
		}
		if typesValidation.hasExhaustedTypes(parameterType) {
			return false
		}
	}
	return true
}

func (i *funcInspector) appendNode(node FuncNode) {
	i.Nodes = append(i.Nodes, node)
	i.Base.increment()
}

func (i *funcInspector) inspect(node ast.Node) {
	i.Base.inspect(node, i.inspector)
}

func (i *funcInspector) newFunction(name string, node ast.Node, comment string) FuncNode {
	baseNode := i.Base.newNode(name, node, comment)
	funcNode := node.(*ast.FuncDecl)
	return FuncNode{Node: baseNode, node: funcNode}
}

func (i *funcInspector) inspector(n ast.Node) bool {
	funcDecl, ok := n.(*ast.FuncDecl)
	if !ok {
		return true
	}
	var comment string = ""
	if funcDecl.Doc != nil {
		comment = funcDecl.Doc.Text()
	}
	name := funcDecl.Name.Name
	funcNode := i.newFunction(name, funcDecl, comment)

	if i.isNodeMatch(funcNode) {
		i.appendNode(funcNode)
	}
	return true
}
