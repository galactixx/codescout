package cmdutils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type EnumOptions[T any] struct {
	Options map[string]func(T) any
}

func (o EnumOptions[T]) EnumValidation(cmd *cobra.Command, flag string, outputType string) error {
	_, outputValid := o.Options[outputType]
	if cmd.Flags().Changed("output") && !outputValid {
		return fmt.Errorf("%v flag must be one of: %v", flag, o.ToOptionString())
	}
	return nil
}

func (o EnumOptions[T]) ToOptionString() string {
	optionsSlice := make([]string, 0, 5)
	for option := range o.Options {
		optionsSlice = append(optionsSlice, option)
	}
	return strings.Join(optionsSlice, ", ")
}

func (o EnumOptions[T]) GetOutputCallable(option string) func(T) any {
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

func ArgsToParams(argTypes []string, parameterTypes *[]codescout.NamedType) error {
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
