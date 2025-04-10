package codescout

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScoutFunction(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type FuncTestCase struct {
		Name     string
		Config   FuncConfig
		Expected BaseNode
	}
	noParamsBool := true
	tests := []FuncTestCase{
		{
			Name:   "test Greet function",
			Config: FuncConfig{ReturnTypes: []string{"string"}, Exact: true},
			Expected: BaseNode{
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
			Expected: BaseNode{
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
			Expected: BaseNode{
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
			assert.Equal(t, tt.Expected, funcNode.Node)
		})
	}
}

func TestScoutMethod(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type MethodTestCase struct {
		Name     string
		Config   MethodConfig
		Expected BaseNode
	}
	noReturnBool := true
	isPointerReceiver := true
	tests := []MethodTestCase{
		{
			Name:   "test Birthday method",
			Config: MethodConfig{NoReturn: &noReturnBool, IsPointerRec: &isPointerReceiver},
			Expected: BaseNode{
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
			Expected: BaseNode{
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
			assert.Equal(t, tt.Expected, mehodNode.Node)
		})
	}
}

func TestScoutStruct(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type StructTestCase struct {
		Name     string
		Config   StructConfig
		Expected BaseNode
	}
	tests := []StructTestCase{
		{
			Name: "test Person struct",
			Config: StructConfig{
				FieldTypes: []NamedType{
					{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"},
				}, Exact: true,
			},
			Expected: BaseNode{
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
			Expected: BaseNode{
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
			structNode, err := ScoutStruct(path, tt.Config)
			assert.NoError(t, err)
			assert.Equal(t, tt.Expected, structNode.Node)
		})
	}
}
