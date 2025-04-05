package codescout

import (
	"errors"
	"fmt"
)

type NamedType struct {
	Name string
	Type string
}

type FuncConfig struct {
	Name        string
	ParamTypes  []NamedType
	ReturnTypes []string
	NoParams    *bool
	NoReturn    *bool
}

type MethodConfig struct {
	Name         string
	ParamTypes   []NamedType
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
	Name       string
	FieldTypes []NamedType
	NoFields   *bool
}

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
	return &(inspector.getNodes())[0], nil
}

func getAllOccurrences[T any](preScout preScoutSetup[T]) ([]T, error) {
	inspector, err := preScout.initializeInspect()
	if err != nil {
		return nil, err
	}
	inspector.inspect()
	return inspector.getNodes(), nil
}

func ScoutFunction(path string, config FuncConfig) (*FuncNode, error) {
	return getFirstOccurrence(funcScoutSetup{Path: path, Config: config}, "function")
}

func ScoutFunctions(path string, config FuncConfig) ([]FuncNode, error) {
	return getAllOccurrences(funcScoutSetup{Path: path, Config: config})
}

func ScoutStruct(path string, config StructConfig) (*StructNode, error) {
	return getFirstOccurrence(structScoutSetup{Path: path, Config: config}, "struct")
}

func ScoutStructs(path string, config StructConfig) ([]StructNode, error) {
	return getAllOccurrences(structScoutSetup{Path: path, Config: config})
}

func ScoutMethod(path string, config MethodConfig) (*MethodNode, error) {
	return getFirstOccurrence(methodScoutSetup{Path: path, Config: config}, "method")
}

func ScoutMethods(path string, config MethodConfig) ([]MethodNode, error) {
	return getAllOccurrences(methodScoutSetup{Path: path, Config: config})
}
