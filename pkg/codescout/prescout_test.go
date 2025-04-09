package codescout

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncScoutSetupInvalidPath(t *testing.T) {
	config := FuncConfig{}
	scouter := funcScoutSetup{Path: "non_existent_file.go", Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestFuncScoutSetupInvalidBatch(t *testing.T) {
	noParamsConfig := true
	config := FuncConfig{ParamTypes: []NamedType{{Name: "name", Type: "string"}}, NoParams: &noParamsConfig}
	f, _ := os.CreateTemp("", "testfile")
	scouter := funcScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestFuncScoutSetupValid(t *testing.T) {
	f, _ := os.CreateTemp("", "testfile")
	config := FuncConfig{ParamTypes: []NamedType{{Name: "name", Type: "string"}}}
	scouter := funcScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.NoError(t, err)
}

func TestStructScoutSetupInvalidPath(t *testing.T) {
	config := StructConfig{}
	scouter := structScoutSetup{Path: "non_existent_file.go", Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestStructScoutSetupInvalidBatch(t *testing.T) {
	noFieldsConfig := true
	config := StructConfig{FieldTypes: []NamedType{{Name: "name", Type: "string"}}, NoFields: &noFieldsConfig}
	f, _ := os.CreateTemp("", "testfile")
	scouter := structScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestStructScoutSetupValid(t *testing.T) {
	f, _ := os.CreateTemp("", "testfile")
	config := StructConfig{FieldTypes: []NamedType{{Name: "name", Type: "string"}}}
	scouter := structScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.NoError(t, err)
}

func TestMethodScoutSetupInvalidPath(t *testing.T) {
	config := StructConfig{}
	scouter := structScoutSetup{Path: "non_existent_file.go", Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestMethodScoutSetupInvalidBatch(t *testing.T) {
	noParamsConfig := true
	config := MethodConfig{ParamTypes: []NamedType{{Name: "name", Type: "string"}}, NoParams: &noParamsConfig}
	f, _ := os.CreateTemp("", "testfile")
	scouter := methodScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.Error(t, err)
}

func TestMethodScoutSetupValid(t *testing.T) {
	f, _ := os.CreateTemp("", "testfile")
	config := MethodConfig{ParamTypes: []NamedType{{Name: "name", Type: "string"}}}
	scouter := methodScoutSetup{Path: f.Name(), Config: config}
	_, err := scouter.initializeInspect()
	assert.NoError(t, err)
}
