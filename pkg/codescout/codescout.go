package codescout

import (
	"errors"
	"fmt"
)

// NamedType represents a named parameter or field with its associated type.
type NamedType struct {
	Name string
	Type string
}

// FuncConfig holds configuration for scouting a function in source code.
type FuncConfig struct {
	// Name of the function.
	Name string
	// Expected parameter types with names.
	ParamTypes []NamedType
	// Expected return types.
	ReturnTypes []string
	// If true, function should have no parameters.
	NoParams *bool
	// If true, function should have no return values.
	NoReturn *bool
	// If true, match must be exact on parameters/returns.
	Exact bool
}

// MethodConfig holds configuration for scouting a method in source code.
type MethodConfig struct {
	// Name of the method.
	Name string
	// Expected parameter types.
	ParamTypes []NamedType
	// Expected return types.
	ReturnTypes []string
	// Type of the receiver.
	Receiver string
	// If true, method must have pointer receiver.
	IsPointerRec *bool
	// Fields that must be present in receiver struct.
	Fields []string
	// Methods that must be present in receiver struct.
	Methods []string
	// If true, method should have no parameters.
	NoParams *bool
	// If true, method should have no return values.
	NoReturn *bool
	// If true, receiver struct must not have fields.
	NoFields *bool
	// If true, receiver struct must not have methods.
	NoMethods *bool
	// If true, all config criteria must match exactly.
	Exact bool
}

// StructConfig holds configuration for scouting a struct type in source code.
type StructConfig struct {
	// Name of the struct.
	Name string
	// Expected field names and types.
	FieldTypes []NamedType
	// If true, struct should not have fields.
	NoFields *bool
	// If true, match must be exact.
	Exact bool
}

// getFirstOccurrence returns the first matching node found by the inspector.
func getFirstOccurrence[T any](preScout preScoutSetup[T], symbol string) (*T, error) {
	inspector, err := preScout.initializeInspect()
	if err != nil {
		return nil, err
	}

	inspector.inspect()
	if len(inspector.getNodes()) == 0 {
		errMsg := fmt.Sprintf("no %s was found based on configuration", symbol)
		err := errors.New(errMsg)
		return nil, err
	}
	return inspector.getNodes()[0], nil
}

// getAllOccurrences returns all matching nodes found by the inspector.
func getAllOccurrences[T any](preScout preScoutSetup[T]) ([]*T, error) {
	inspector, err := preScout.initializeInspect()
	if err != nil {
		return nil, err
	}
	inspector.inspect()
	return inspector.getNodes(), nil
}

// ScoutFunction returns the first function in the given path matching the config.
func ScoutFunction(path string, config FuncConfig) (*FuncNode, error) {
	return getFirstOccurrence(funcScoutSetup{Path: path, Config: config}, "function")
}

// ScoutFunctions returns all functions in the given path matching the config.
func ScoutFunctions(path string, config FuncConfig) ([]*FuncNode, error) {
	return getAllOccurrences(funcScoutSetup{Path: path, Config: config})
}

// ScoutStruct returns the first struct in the given path matching the config.
func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	return getFirstOccurrence(structScoutSetup{Path: path, Config: config}, "struct")
}

// ScoutStructs returns all structs in the given path matching the config.
func ScoutStructs(path string, config StructConfig) ([]*StructNode, error) {
	return getAllOccurrences(structScoutSetup{Path: path, Config: config})
}

// ScoutMethod returns the first method in the given path matching the config.
func ScoutMethod(path string, config MethodConfig) (*MethodNode, error) {
	return getFirstOccurrence(methodScoutSetup{Path: path, Config: config}, "method")
}

// ScoutMethods returns all methods in the given path matching the config.
func ScoutMethods(path string, config MethodConfig) ([]*MethodNode, error) {
	return getAllOccurrences(methodScoutSetup{Path: path, Config: config})
}
