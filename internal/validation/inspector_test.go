package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeValidation(t *testing.T) {
	typeValidation := TypeValidation{
		TypeMap:       map[string]int{"string": 2, "bool": 1, "int": 1},
		NamedTypesMap: map[string]string{"name": "string", "age": "int", "sex": "string", "married": "bool"},
	}

	assert.Equal(t, "int", typeValidation.GetParamType("age"))

	typeValidation.SetParamName("name")
	assert.Equal(t, "name", typeValidation.CurParamName)

	typeValidation.SetParamType("string")
	assert.Equal(t, "string", typeValidation.CurParamType)

	typeValidation.SetNameInParams("age")
	assert.True(t, typeValidation.CurNameInParams)

	typeValidation.SetNameInParams("lastName")
	assert.False(t, typeValidation.CurNameInParams)

	typeValidation.SetParamInfo("age", "int")
	assert.Equal(t, "age", typeValidation.CurParamName)
	assert.Equal(t, "int", typeValidation.CurParamType)

	typeValidation.SetNameInParams("age")
	assert.True(t, typeValidation.TypeExists())
	assert.True(t, typeValidation.TypeExclExists())

	assert.False(t, typeValidation.HasExhausted("bool"))
	assert.True(t, typeValidation.HasExhausted("bool"))
}
