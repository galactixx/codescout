package codescout

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func validator(cfgTypes []string, nodeTypes []string) bool { return true }

func TestAstMatch(t *testing.T) {
	noBoolVar := true
	astMatchSlice := astMatch([]string{}, []string{}, true, &noBoolVar, validator)
	assert.IsType(t, aSTNodeSliceMatch[[]string, []string]{}, astMatchSlice)
}

func TestAstMatchesReturnMatchExact(t *testing.T) {
	configTypes := []string{"string", "int"}
	nodeTypes := []string{"string", "int", "bool"}
	astMatchSlice := astMatch(configTypes, nodeTypes, true, nil, returnMatch)

	assert.False(t, astMatchSlice.runValidator())
	assert.False(t, astMatchSlice.noTypesMatch())
	assert.False(t, astMatchSlice.validate())
}

func TestAstMatchesReturnMatchNonExact(t *testing.T) {
	configTypes := []string{"string", "int"}
	nodeTypes := []string{"string", "int", "bool"}
	astMatchSlice := astMatch(configTypes, nodeTypes, false, nil, returnMatch)

	assert.True(t, astMatchSlice.runValidator())
	assert.False(t, astMatchSlice.noTypesMatch())
	assert.True(t, astMatchSlice.validate())
}

func TestAstMatchesAccessedMatchNonExact(t *testing.T) {
	configTypes := []string{"Name", "Types", "Age"}
	nodeTypes := []string{"Name", "Types"}
	astMatchSlice := astMatch(configTypes, nodeTypes, false, nil, accessedMatch)

	assert.False(t, astMatchSlice.runValidator())
	assert.False(t, astMatchSlice.noTypesMatch())
	assert.False(t, astMatchSlice.validate())
}

func TestAstMatchesAccessedMatchNoAccess(t *testing.T) {
	noBoolVar := true
	astMatchSlice := astMatch([]string{}, []string{}, false, &noBoolVar, accessedMatch)

	assert.False(t, astMatchSlice.runValidator())
	assert.True(t, astMatchSlice.noTypesMatch())
	assert.True(t, astMatchSlice.validate())
}

func TestAstMatchesNamedTypesExact(t *testing.T) {
	configTypes := []NamedType{{Name: "Name", Type: "string"}}
	nodeTypes := []NamedType{{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"}}
	astMatchSlice := astMatch(configTypes, nodeTypes, true, nil, namedTypesMatch)

	assert.False(t, astMatchSlice.runValidator())
	assert.False(t, astMatchSlice.noTypesMatch())
	assert.False(t, astMatchSlice.validate())
}

func TestAstMatchNamedTypesNoTypes(t *testing.T) {
	noBoolVar := true
	astMatchSlice := astMatch([]NamedType{}, []NamedType{}, false, &noBoolVar, namedTypesMatch)

	assert.False(t, astMatchSlice.runValidator())
	assert.True(t, astMatchSlice.noTypesMatch())
	assert.True(t, astMatchSlice.validate())
}

func TestNamedTypesMapOfTypes(t *testing.T) {
	namedTypesMapping := namedTypesMapOfTypes([]NamedType{{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"}})
	assert.True(t, namedTypesMapping["string"] == 1)
	assert.True(t, namedTypesMapping["int"] == 1)
	assert.Len(t, namedTypesMapping, 2)
}

func TestNamedTypesMap(t *testing.T) {
	namedTypesMapping := namedTypesMap([]NamedType{{Name: "Name", Type: "string"}, {Name: "Age", Type: "int"}})
	assert.True(t, namedTypesMapping["Name"] == "string")
	assert.True(t, namedTypesMapping["Age"] == "int")
	assert.Len(t, namedTypesMapping, 2)
}
