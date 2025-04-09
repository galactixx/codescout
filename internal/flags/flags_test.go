package flags

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestStringBoolToPointer(t *testing.T) {
	assert.Nil(t, StringBoolToPointer(""))

	truePtr := StringBoolToPointer("true")
	falsePtr := StringBoolToPointer("false")
	assert.NotNil(t, truePtr)
	assert.NotNil(t, falsePtr)
	assert.True(t, *truePtr)
	assert.False(t, *falsePtr)
}

func TestStringVarP(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[string]{Name: "name"}
	StringVarP(cmd, &flag, "n", "default", "usage")

	err := cmd.ParseFlags([]string{"--name=test"})
	assert.NoError(t, err)
	assert.Equal(t, "test", flag.Variable)
}

func TestStringSliceVarP(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[[]string]{Name: "tags"}
	StringSliceVarP(cmd, &flag, "t", make([]string, 0), "usage")

	err := cmd.ParseFlags([]string{"--tags=x,y"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"x", "y"}, flag.Variable)
}

func TestBoolVarP(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[bool]{Name: "active"}
	BoolVarP(cmd, &flag, "a", false, "usage")

	err := cmd.ParseFlags([]string{"--active"})
	assert.NoError(t, err)
	assert.True(t, flag.Variable)
}

func TestEmptyValidator_String(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[string]{Name: "email"}
	StringVarP(cmd, &flag, "e", "", "usage")

	assert.False(t, flag.emptyValidator(cmd))

	_ = cmd.ParseFlags([]string{"--email="})
	assert.True(t, flag.emptyValidator(cmd))
}

func TestEmptyValidator_Slice(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[[]string]{Name: "items"}
	StringSliceVarP(cmd, &flag, "i", make([]string, 0), "usage")

	_ = cmd.ParseFlags([]string{"--items="})
	assert.True(t, flag.emptyValidator(cmd))

	_ = cmd.ParseFlags([]string{"--items=x"})
	assert.False(t, flag.emptyValidator(cmd))
}

func TestBatchValidator_StringBoolValidator(t *testing.T) {
	cmd := &cobra.Command{}
	flag := CommandFlag[string]{Name: "enabled"}
	StringVarP(cmd, &flag, "e", "", "usage")

	_ = cmd.ParseFlags([]string{"--enabled=maybe"})
	validator := BatchValidator{
		StringBoolValidators: []*CommandFlag[string]{&flag},
	}
	err := validator.Validate(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be: true or false")
}

func TestBatchValidator_EmptyValidator(t *testing.T) {
	cmd := &cobra.Command{}
	strFlag := CommandFlag[string]{Name: "name"}
	sliceFlag := CommandFlag[[]string]{Name: "list"}
	StringVarP(cmd, &strFlag, "n", "", "usage")
	StringSliceVarP(cmd, &sliceFlag, "l", make([]string, 0), "usage")

	_ = cmd.ParseFlags([]string{"--name=", "--list="})

	validator := BatchValidator{
		EmptyValidators: []FlagValidator{&strFlag, &sliceFlag},
	}
	err := validator.Validate(cmd)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestBatchValidator_AllValid(t *testing.T) {
	cmd := &cobra.Command{}
	strFlag := CommandFlag[string]{Name: "verbose"}
	StringVarP(cmd, &strFlag, "v", "true", "usage")

	_ = cmd.ParseFlags([]string{"--verbose=true"})

	validator := BatchValidator{
		StringBoolValidators: []*CommandFlag[string]{&strFlag},
	}
	err := validator.Validate(cmd)
	assert.NoError(t, err)
}
