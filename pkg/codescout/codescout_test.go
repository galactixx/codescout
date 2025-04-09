package codescout

import (
	"path/filepath"
	"reflect"
	"testing"
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

		eValue := expectedV.FieldByName(name).Interface()
		aValue := actualV.FieldByName(name).Interface()
		if aValue != eValue {
			t.Errorf("field %s - expected %s, got %s", name, eValue, aValue)
		}
	}
}

func TestScoutFunction(t *testing.T) {
	path := filepath.Join("testdata", "scout_single.go")
	type FuncTestCaseExpected struct {
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
	tests := []FuncTestCase{
		{
			Name:   "test Greet function",
			Config: FuncConfig{Name: "Greet"},
			Expected: FuncTestCaseExpected{
				Path:       path,
				Line:       20,
				Characters: 1,
				Exported:   true,
				Comment:    "Above above function\nAbove function\nFunction",
			},
		},
		{
			Name:   "test main function",
			Config: FuncConfig{Name: "main"},
			Expected: FuncTestCaseExpected{
				Path:       path,
				Line:       44,
				Characters: 1,
				Exported:   false,
				Comment:    "",
			},
		},
		{
			Name:   "test Factorial function",
			Config: FuncConfig{Name: "Factorial"},
			Expected: FuncTestCaseExpected{
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
			if err != nil {
				t.Errorf("got %v, expected %v", err, nil)
			}
			compareStructValues(t, funcNode.Node, tt.Expected)
		})
	}
}
