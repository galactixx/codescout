package validation

import "testing"

func TestTypeValidation(t *testing.T) {
	typeValidation := TypeValidation{
		TypeMap:       map[string]int{"string": 2, "bool": 1, "int": 1},
		NamedTypesMap: map[string]string{"name": "string", "age": "int", "sex": "string", "married": "bool"},
	}

	if typeValidation.GetParamType("age") != "int" {
		t.Errorf("Expected %s, got %s", "int", typeValidation.GetParamType("age"))
	}

	typeValidation.SetParamName("name")
	if typeValidation.CurParamName != "name" {
		t.Errorf("Expected %s, got %s", "name", typeValidation.CurParamName)
	}

	typeValidation.SetParamType("string")
	if typeValidation.CurParamType != "string" {
		t.Errorf("Expected %s, got %s", "string", typeValidation.CurParamType)
	}

	typeValidation.SetNameInParams("age")
	if !typeValidation.CurNameInParams {
		t.Errorf("Expected %v, got %v", true, false)
	}

	typeValidation.SetNameInParams("lastName")
	if typeValidation.CurNameInParams {
		t.Errorf("Expected %v, got %v", false, true)
	}

	typeValidation.SetParamInfo("age", "int")
	if typeValidation.CurParamName != "age" {
		t.Errorf("Expected %s, got %s", "age", typeValidation.CurParamName)
	}

	if typeValidation.CurParamType != "int" {
		t.Errorf("Expected %s, got %s", "int", typeValidation.CurParamType)
	}

	typeValidation.SetNameInParams("age")
	if !typeValidation.TypeExists() {
		t.Errorf("Expected %v, got %v", true, false)
	}

	if !typeValidation.TypeExclExists() {
		t.Errorf("Expected %v, got %v", true, false)
	}

	if exhausted := typeValidation.HasExhausted("bool"); exhausted {
		t.Errorf("Expected %v, got %v", false, true)
	}

	if exhausted := typeValidation.HasExhausted("bool"); !exhausted {
		t.Errorf("Expected %v, got %v", true, false)
	}
}
