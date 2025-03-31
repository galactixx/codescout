package codescout

import (
	"fmt"
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
		fmt.Println(name)

		eValue := expectedV.FieldByName(name).Interface()
		aValue := actualV.FieldByName(name).Interface()
		if aValue != eValue {
			t.Errorf("field %v - expected %v, got %v", name, eValue, aValue)
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
			Name:   "test simple function",
			Config: FuncConfig{Name: "Greet"},
			Expected: FuncTestCaseExpected{
				Path:       path,
				Line:       20,
				Characters: 1,
				Exported:   true,
				Comment:    "Above above function\nAbove function\nFunction\n",
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
