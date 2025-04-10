package codescout

import (
	"go/token"

	"github.com/galactixx/codescout/internal/pkgutils"
	"github.com/galactixx/codescout/internal/validation"
)

// preScoutSetup defines an interface for initializing an inspector of a generic type T.
type preScoutSetup[T any] interface {
	initializeInspect() (inspector[T], error)
}

// funcScoutSetup holds configuration for scanning functions.
type funcScoutSetup struct {
	Path   string
	Config FuncConfig
}

// initializeInspect validates function-related configuration and returns an inspector for FuncNode.
//
//lint:ignore U1000 used via interface
func (s funcScoutSetup) initializeInspect() (inspector[FuncNode], error) {
	// Check if the provided path exists.
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	// Create validation rules for function parameters and return types.
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

	// Run batch validation and return an error if it fails.
	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	// Create and return the function inspector.
	inspector := funcInspector{
		Nodes:  []*FuncNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}

// methodScoutSetup holds configuration for scanning methods.
type methodScoutSetup struct {
	Path   string
	Config MethodConfig
}

// initializeInspect validates method-related configuration and returns an inspector for MethodNode.
//
//lint:ignore U1000 used via interface
func (s methodScoutSetup) initializeInspect() (inspector[MethodNode], error) {
	// Check if the provided path exists.
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	// Create validation rules for method fields, methods, return types, and parameters.
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

	// Run batch validation and return an error if it fails.
	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	// Create and return the method inspector.
	inspector := methodInspector{
		Nodes:  []*MethodNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}

// structScoutSetup holds configuration for scanning structs.
type structScoutSetup struct {
	Path   string
	Config StructConfig
}

// initializeInspect validates struct-related configuration and returns an inspector for StructNode.
//
//lint:ignore U1000 used via interface
func (s structScoutSetup) initializeInspect() (inspector[StructNode], error) {
	// Check if the provided path exists.
	if fileExistsErr := pkgutils.FilePathExists(s.Path); fileExistsErr != nil {
		return nil, fileExistsErr
	}

	// Create validation rules for struct fields.
	batchValidation := validation.BatchConfigValidation{
		SliceValidators: []validation.SliceValidator{
			validation.SlicePairToValidate[NamedType]{
				Slice: validation.Arg("FieldTypes", s.Config.FieldTypes),
				Bool:  validation.Arg("NoFields", s.Config.NoFields),
			},
		},
		Exact: s.Config.Exact,
	}

	// Run batch validation and return an error if it fails.
	batchErr := batchValidation.Validate()
	if batchErr != nil {
		return nil, batchErr
	}

	// Create and return the struct inspector.
	inspector := structInspector{
		Nodes:  map[string]*StructNode{},
		Config: s.Config,
		Base:   baseInspector{Path: s.Path, Fset: token.NewFileSet()},
	}
	return &inspector, nil
}
