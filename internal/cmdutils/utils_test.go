package cmdutils

import (
	"testing"

	"github.com/galactixx/codescout/internal/flags"
	"github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestJoinAttrs(t *testing.T) {
	attrs := []string{"one", "two", "three"}
	expected := "[ one,two,three ]"
	assert.Equal(t, expected, JoinAttrs(attrs))
}

func TestStripANSI(t *testing.T) {
	input := "\x1b[31mRed Text\x1b[0m"
	expected := "Red Text"
	assert.Equal(t, expected, stripANSI(input))
}

func TestOutputOptionsValidation(t *testing.T) {
	opt := OutputOptions[string]{
		Options: map[string]func(string) string{
			"json": func(s string) string { return s },
		},
	}
	flag := flags.CommandFlag[string]{Name: "output", Variable: "yaml"}
	cmd := &cobra.Command{}
	cmd.Flags().String("output", "", "")
	_ = cmd.ParseFlags([]string{"--output=yaml"})

	err := opt.validation(cmd, flag)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be one of")
}

func TestOutputOptionsToOptionString(t *testing.T) {
	opt := OutputOptions[string]{
		Options: map[string]func(string) string{
			"json": func(s string) string { return s },
			"text": func(s string) string { return s },
		},
	}
	optStr := opt.ToOptionString()
	assert.Contains(t, optStr, "json")
	assert.Contains(t, optStr, "text")
}

func TestArgsToNamedTypesValid(t *testing.T) {
	args := []string{"name:string", "age:int"}
	var result []codescout.NamedType
	err := argsToNamedTypes(args, &result)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "name", result[0].Name)
	assert.Equal(t, "int", result[1].Type)
}

func TestArgsToNamedTypesInvalidColonCount(t *testing.T) {
	args := []string{"badarg"}
	var result []codescout.NamedType
	err := argsToNamedTypes(args, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "one colon")
}

func TestArgsToNamedTypesEmptyParts(t *testing.T) {
	args := []string{":"}
	var result []codescout.NamedType
	err := argsToNamedTypes(args, &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be defined")
}
