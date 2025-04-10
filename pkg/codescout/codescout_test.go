package codescout

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func compareStructValues(t *testing.T, actual any, expected any) {
	actualV := reflect.ValueOf(actual)
	expectedV := reflect.ValueOf(expected)
	expectedType := expectedV.Type()

	if actualV.Kind() == reflect.Ptr {
		actualV = actualV.Elem()
	}

	if expectedV.Kind() == reflect.Ptr {
		expectedV = expectedV.Elem()
	}

	for i := 0; i < expectedType.NumField(); i++ {
		name := expectedType.Field(i).Name
		assert.Equal(
			t,
			expectedV.FieldByName(name).Interface(),
			actualV.FieldByName(name).Interface(),
		)
	}
}

func TestScoutFunction(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type FuncTestCaseExpected struct {
		Name       string
		Path       string
		Line       int
		Characters int
		Exported   bool
		Comment    string
	}
	type FuncTestCase struct {
		Name     string
		Config   FuncConfig
		Expected FuncTestCaseExpected
	}
	noParamsBool := true
	tests := []FuncTestCase{
		{
			Name:   "test Greet function",
			Config: FuncConfig{ReturnTypes: []string{"string"}, Exact: true},
			Expected: FuncTestCaseExpected{
				Name:       "Greet",
				Path:       path,
				Line:       20,
				Characters: 1,
				Exported:   true,
				Comment:    "Above above function\nAbove function\nFunction",
			},
		},
		{
			Name:   "test main function",
			Config: FuncConfig{NoParams: &noParamsBool},
			Expected: FuncTestCaseExpected{
				Name:       "main",
				Path:       path,
				Line:       44,
				Characters: 1,
				Exported:   false,
				Comment:    "",
			},
		},
		{
			Name:   "test Factorial function",
			Config: FuncConfig{ParamTypes: []NamedType{{Type: "int"}}, ReturnTypes: []string{"int"}},
			Expected: FuncTestCaseExpected{
				Name:       "Factorial",
				Path:       path,
				Line:       54,
				Characters: 1,
				Exported:   true,
				Comment:    "Factorial function",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			funcNode, err := ScoutFunction(path, tt.Config)
			assert.NoError(t, err)
			compareStructValues(t, funcNode.Node, tt.Expected)
		})
	}
}

func TestScoutMethod(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type MethodTestCaseExpected struct {
		Name       string
		Path       string
		Line       int
		Characters int
		Exported   bool
		Comment    string
	}
	type MethodTestCase struct {
		Name     string
		Config   MethodConfig
		Expected MethodTestCaseExpected
	}
	noReturnBool := true
	isPointerReceiver := true
	tests := []MethodTestCase{
		{
			Name:   "test Birthday method",
			Config: MethodConfig{NoReturn: &noReturnBool, IsPointerRec: &isPointerReceiver},
			Expected: MethodTestCaseExpected{
				Name:       "Birthday",
				Path:       path,
				Line:       27,
				Characters: 1,
				Exported:   true,
				Comment:    "Method on Person struct",
			},
		},
		{
			Name:   "test DisplayDetails method",
			Config: MethodConfig{ReturnTypes: []string{"string"}, Receiver: "Car"},
			Expected: MethodTestCaseExpected{
				Name:       "DisplayDetails",
				Path:       path,
				Line:       32,
				Characters: 1,
				Exported:   true,
				Comment:    "Method on Car struct",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mehodNode, err := ScoutMethod(path, tt.Config)
			assert.NoError(t, err)
			compareStructValues(t, mehodNode.Node, tt.Expected)
		})
	}
}

func TestScoutStruct(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type StructTestCaseExpected struct {
		Name       string
		Path       string
		Line       int
		Characters int
		Exported   bool
		Comment    string
	}
	type StructTestCase struct {
		Name     string
		Config   StructConfig
		Expected StructTestCaseExpected
	}
	tests := []StructTestCase{
		{
			Name: "test Person struct",
			Config: StructConfig{
				FieldTypes: []NamedType{
					{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"},
				}, Exact: true,
			},
			Expected: StructTestCaseExpected{
				Name:       "Person",
				Path:       path,
				Line:       6,
				Characters: 13,
				Exported:   true,
				Comment:    "Structs",
			},
		},
		{
			Name:   "test Car struct",
			Config: StructConfig{FieldTypes: []NamedType{{Name: "Make", Type: "string"}}},
			Expected: StructTestCaseExpected{
				Name:       "Car",
				Path:       path,
				Line:       11,
				Characters: 10,
				Exported:   true,
				Comment:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			mehodNode, err := ScoutStruct(path, tt.Config)
			assert.NoError(t, err)
			compareStructValues(t, mehodNode.Node, tt.Expected)
		})
	}
}
