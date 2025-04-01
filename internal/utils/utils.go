package utils

import (
	"errors"
	"strings"

	"github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func CountFlagsSet(cmd *cobra.Command) int {
	count := 0
	cmd.Flags().Visit(func(f *pflag.Flag) {
		count++
	})
	return count
}

func ArgsToParams(argTypes []string, parameterTypes *[]codescout.Parameter) error {
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
		param := codescout.Parameter{Name: paramName, Type: paramType}
		*parameterTypes = append(*parameterTypes, param)
	}
	return nil
}
