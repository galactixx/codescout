package codescout

import (
	"go/token"

	"github.com/galactixx/codescout/internal/pkgutils"
	"github.com/galactixx/codescout/internal/validation"
)

type preScoutSetup[T any] interface {
	initializeInspect() (inspector[T], error)
}

type funcScoutSetup struct {
	Path   string
	Config FuncConfig
}

//lint:ignore U1000 used via interface
func (s funcScoutSetup) initializeInspect() (inspector[FuncNode], error) {
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[NamedType]{
				Slice: validation.Arg("Types", s.Config.ParamTypes),
				Bool:  validation.Arg("NoParams", s.Config.NoParams),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("ReturnTypes", s.Config.ReturnTypes),
				Bool:  validation.Arg("NoReturn", s.Config.NoReturn),
			},
		},
		Exact: s.Config.Exact,
	}

	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	inspector := funcInspector{
		Nodes:  []*FuncNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}

type methodScoutSetup struct {
	Path   string
	Config MethodConfig
}

//lint:ignore U1000 used via interface
func (s methodScoutSetup) initializeInspect() (inspector[MethodNode], error) {
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("Fields", s.Config.Fields),
				Bool:  validation.Arg("NoFields", s.Config.NoFields),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("Methods", s.Config.Methods),
				Bool:  validation.Arg("NoMethods", s.Config.NoMethods),
			},
			validation.SlicePairToValidate[string]{
				Slice: validation.Arg("ReturnTypes", s.Config.ReturnTypes),
				Bool:  validation.Arg("NoReturn", s.Config.NoReturn),
			},
			validation.SlicePairToValidate[NamedType]{
				Slice: validation.Arg("Types", s.Config.ParamTypes),
				Bool:  validation.Arg("NoParams", s.Config.NoParams),
			},
		},
		Exact: s.Config.Exact,
	}

	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	inspector := methodInspector{
		Nodes:  []*MethodNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}

type structScoutSetup struct {
	Path   string
	Config StructConfig
}

//lint:ignore U1000 used via interface
func (s structScoutSetup) initializeInspect() (inspector[StructNode], error) {
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[NamedType]{
				Slice: validation.Arg("FieldTypes", s.Config.FieldTypes),
				Bool:  validation.Arg("NoFields", s.Config.NoFields),
			},
		},
		Exact: s.Config.Exact,
	}

	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	inspector := structInspector{
		Nodes:  map[string]*StructNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}
