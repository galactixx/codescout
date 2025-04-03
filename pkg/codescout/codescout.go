package codescout

import (
	"go/token"

	"github.com/galactixx/codescout/internal/pkgutils"
	"github.com/galactixx/codescout/internal/validation"
)

type Parameter struct {
	Name string
	Type string
}

type FuncConfig struct {
	Name        string
	Types       []Parameter
	ReturnTypes []string
	NoParams    *bool
	NoReturn    *bool
}

type MethodConfig struct {
	Name             string
	Types            []Parameter
	ReturnTypes      []string
	Receiver         string
	IsPointerRec     *bool
	FieldsAccessed   []string
	MethodsCalled    []string
	NoParams         *bool
	NoReturn         *bool
	NoFieldsAccessed *bool
	NoMethodsCalled  *bool
}

type StructConfig struct {
	Name    string
	Types   []Parameter
	Methods []FuncConfig
}

func ScoutFunction(path string, config FuncConfig) (*FuncNode, error) {
	if fileExistsErr := pkgutils.FilePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[Parameter]{
				Slice: validation.Argument[[]Parameter]{Name: "Types", Variable: config.Types},
				Bool:  validation.Argument[*bool]{Name: "NoParams", Variable: config.NoParams},
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Argument[[]string]{Name: "ReturnTypes", Variable: config.ReturnTypes},
				Bool:  validation.Argument[*bool]{Name: "NoReturn", Variable: config.NoReturn},
			},
		},
	}

	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	inspector := funcInspector{
		Nodes:  []FuncNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: token.NewFileSet()},
	}
	inspector.inspect()
	return inspectorGetNode(&inspector, "function")
}

func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	if fileExistsErr := pkgutils.FilePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	inspector := structInspector{
		Nodes:  []StructNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: token.NewFileSet()},
	}
	inspector.inspect()
	return inspectorGetNode(&inspector, "struct")
}

func ScoutMethod(path string, config MethodConfig) (*MethodNode, error) {
	if fileExistsErr := pkgutils.FilePathExists(path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[string]{
				Slice: validation.Argument[[]string]{Name: "FieldsAccessed", Variable: config.FieldsAccessed},
				Bool:  validation.Argument[*bool]{Name: "NoFieldsAccessed", Variable: config.NoFieldsAccessed},
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Argument[[]string]{Name: "MethodsCalled", Variable: config.MethodsCalled},
				Bool:  validation.Argument[*bool]{Name: "NoMethodsCalled", Variable: config.NoMethodsCalled},
			},
			validation.SlicePairToValidate[Parameter]{
				Slice: validation.Argument[[]Parameter]{Name: "Types", Variable: config.Types},
				Bool:  validation.Argument[*bool]{Name: "NoParams", Variable: config.NoParams},
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Argument[[]string]{Name: "ReturnTypes", Variable: config.ReturnTypes},
				Bool:  validation.Argument[*bool]{Name: "NoReturn", Variable: config.NoReturn},
			},
		},
	}

	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	inspector := methodInspector{
		Nodes:  []MethodNode{},
		Config: config,
		Base:   baseInspector{Path: path, Fset: token.NewFileSet()},
	}
	inspector.inspect()
	return inspectorGetNode(&inspector, "method")
}
