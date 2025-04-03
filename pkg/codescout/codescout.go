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
	Name         string
	Types        []Parameter
	ReturnTypes  []string
	Receiver     string
	IsPointerRec *bool
	Fields       []string
	Methods      []string
	NoParams     *bool
	NoReturn     *bool
	NoFields     *bool
	NoMethods    *bool
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
				Slice: validation.Arg("Types", config.Types),
				Bool:  validation.Arg("NoParams", config.NoParams),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("ReturnTypes", config.ReturnTypes),
				Bool:  validation.Arg("NoReturn", config.NoReturn),
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
				Slice: validation.Arg("Fields", config.Fields),
				Bool:  validation.Arg("NoFields", config.NoFields),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("Methods", config.Methods),
				Bool:  validation.Arg("NoMethods", config.NoMethods),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("ReturnTypes", config.ReturnTypes),
				Bool:  validation.Arg("NoReturn", config.NoReturn),
			},
			validation.SlicePairToValidate[Parameter]{
				Slice: validation.Arg("Types", config.Types),
				Bool:  validation.Arg("NoParams", config.NoParams),
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
