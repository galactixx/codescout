package cmdutils

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/galactixx/codescout/internal/flags"
	"github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type OutputOptions[T any] struct {
	Options map[string]func(T) any
}

func (o OutputOptions[T]) EnumValidation(cmd *cobra.Command, flag flags.CommandFlag[string]) error {
	_, outputValid := o.Options[flag.Variable]
	if cmd.Flags().Changed("output") && !outputValid {
		return fmt.Errorf("%v flag must be one of: %v", flag.Name, o.ToOptionString())
	}
	return nil
}

func (o OutputOptions[T]) ToOptionString() string {
	optionsSlice := make([]string, 0, 5)
	for option := range o.Options {
		optionsSlice = append(optionsSlice, option)
	}
	return strings.Join(optionsSlice, ", ")
}

func (o OutputOptions[T]) GetOutputCallable(option string) func(T) any {
	return o.Options[option]
}

func CountFlagsSet(cmd *cobra.Command) int {
	count := 0
	cmd.Flags().Visit(func(f *pflag.Flag) {
		count++
	})
	return count
}

func FromStringToBool(stringBool string) *bool {
	newBool, err := strconv.ParseBool(stringBool)
	if err == nil {
		return nil
	}
	return &newBool
}

func ArgsToNamedTypes(argTypes []string, parameterTypes *[]codescout.NamedType) error {
	for _, parameter := range argTypes {
		if strings.Count(parameter, ":") != 1 {
			return errors.New("there must be only one colon separating out the name and type")
		}

		paramDestruct := strings.SplitN(parameter, ":", 2)
		paramName := strings.TrimSpace(paramDestruct[0])
		paramType := strings.TrimSpace(paramDestruct[1])

		if paramName == "" && paramType == "" {
			return errors.New("at least one of the type or name must be defined")
		}
		param := codescout.NamedType{Name: paramName, Type: paramType}
		*parameterTypes = append(*parameterTypes, param)
	}
	return nil
}

type CobraCommandVlidation[T any] struct {
	Validator      flags.BatchValidator
	NamedTypesFlag flags.CommandFlag[[]string]
	OutputTypeFlag flags.CommandFlag[string]
	OutputOptions  OutputOptions[T]

	namedTypes []codescout.NamedType
}

func (v *CobraCommandVlidation[T]) GetNamedTypes() []codescout.NamedType {
	if v.namedTypes == nil {
		log.Fatal("named types field is returning nil, should never occur")
	}
	namedTypes := v.namedTypes
	v.namedTypes = nil
	return namedTypes
}

func (v *CobraCommandVlidation[T]) CommandValidation(cmd *cobra.Command) error {
	if CountFlagsSet(cmd) == 0 {
		return errors.New("at least one flag must be set for this command")
	}

	validationErr := v.Validator.Validate(cmd)
	if validationErr != nil {
		return validationErr
	}

	namedTypes := make([]codescout.NamedType, 0, 5)
	err := ArgsToNamedTypes(v.NamedTypesFlag.Variable, &namedTypes)
	if err != nil {
		return err
	}
	v.namedTypes = namedTypes

	outputErr := v.OutputOptions.EnumValidation(cmd, v.OutputTypeFlag)
	if outputErr != nil {
		return outputErr
	}
	return nil
}
