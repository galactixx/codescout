package validation

import (
	"fmt"
	"reflect"
	"testing"
)

func TestArg(t *testing.T) {
	tests := []struct {
		name     string
		variable []string
	}{
		{"Params", []string{"string", "int"}},
		{"Returns", []string{"bool"}},
	}

	for _, tt := range tests {
		createdArg := Arg(tt.name, tt.variable)
		if reflect.TypeOf(createdArg) != reflect.TypeOf(Argument[[]string]{}) {
			expectedType := reflect.TypeOf(Argument[[]string]{}).String()
			gotType := reflect.TypeOf(createdArg).String()
			t.Errorf("Expected %s, got %s", expectedType, gotType)
		}
	}
}

func TestSlicePairToValidate(t *testing.T) {
	boolVariable := false
	slicePairToValidate := SlicePairToValidate[string]{
		Slice: Argument[[]string]{Name: "Params", Variable: []string{"string", "int"}},
		Bool:  Argument[*bool]{Name: "NoParams", Variable: &boolVariable},
	}

	trueMessage := fmt.Errorf("Params cannot be specified if NoParams is set to true")
	if slicePairToValidate.trueMessage().Error() != trueMessage.Error() {
		t.Errorf("Expected %v, got %v", trueMessage, slicePairToValidate.trueMessage())
	}

	falseMessage := fmt.Errorf("no need to specify Params if NoParams is set to false")
	if slicePairToValidate.falseMessage().Error() != falseMessage.Error() {
		t.Errorf("Expected %v, got %v", falseMessage, slicePairToValidate.falseMessage())
	}

	if !slicePairToValidate.boolNotNil() {
		t.Errorf("Expected %v, got %v", true, false)
	}

	if !slicePairToValidate.nonEmptySlice() {
		t.Errorf("Expected %v, got %v", true, false)
	}

	if slicePairToValidate.trueValidate() {
		t.Errorf("Expected %v, got %v", false, true)
	}

	if !slicePairToValidate.falseValidate() {
		t.Errorf("Expected %v, got %v", true, false)
	}
}

func TestBatchConfigValidation(t *testing.T) {
	boolVariable := false
	batchvalidation := BatchConfigValidation{
		SliceValidators: []SliceValidator{
			SlicePairToValidate[string]{
				Slice: Argument[[]string]{Name: "Params", Variable: []string{"string", "int"}},
				Bool:  Argument[*bool]{Name: "NoParams", Variable: &boolVariable},
			},
		},
		Exact: true,
	}

	exactMessage := "exact should not be true if no slices are passed"
	if batchvalidation.exactMessage().Error() != exactMessage {
		t.Errorf("Expected %s, got %s", exactMessage, batchvalidation.exactMessage().Error())
	}

	if batchvalidation.Validate() == nil {
		t.Errorf("Expected %v, got %v", nil, batchvalidation.Validate())
	}
}
