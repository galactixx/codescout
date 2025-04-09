package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.IsType(t, Argument[[]string]{}, createdArg)
	}
}

func TestSlicePairToValidate(t *testing.T) {
	boolVariable := false
	slicePairToValidate := SlicePairToValidate[string]{
		Slice: Argument[[]string]{Name: "Params", Variable: []string{"string", "int"}},
		Bool:  Argument[*bool]{Name: "NoParams", Variable: &boolVariable},
	}

	trueMessage := fmt.Errorf("Params cannot be specified if NoParams is set to true")
	assert.Equal(t, trueMessage, slicePairToValidate.trueMessage())

	falseMessage := fmt.Errorf("no need to specify Params if NoParams is set to false")
	assert.Equal(t, falseMessage, slicePairToValidate.falseMessage())

	assert.True(t, slicePairToValidate.boolNotNil())
	assert.True(t, slicePairToValidate.nonEmptySlice())
	assert.False(t, slicePairToValidate.trueValidate())
	assert.True(t, slicePairToValidate.falseValidate())
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
	assert.Equal(t, exactMessage, batchvalidation.exactMessage().Error())
	assert.NotNil(t, batchvalidation.Validate())
}
