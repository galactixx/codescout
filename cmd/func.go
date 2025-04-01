package cmd

import (
	"errors"
	"strings"

	codescout "github.com/galactixx/codescout/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	funcName       string
	returnTypes    []string
	parameterTypes []string
	outputType     string
)

var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "Find a single function in a file",
	Long:  "Locate and display a specific function definition by name within a given source file.",
	Args:  cobra.ExactArgs(1),
	RunE:  funcCmdRun,
}

func init() {
	rootCmd.AddCommand(funcCmd)

	funcCmd.Flags().StringVarP(&funcName, "name", "n", "", "Function name to use")
	funcCmd.Flags().StringSliceVarP(&parameterTypes, "params", "p", make([]string, 0), "Parameter names and types of function")
	funcCmd.Flags().StringSliceVarP(&returnTypes, "return", "r", make([]string, 0), "Return types of function")
	funcCmd.Flags().StringVarP(
		&outputType,
		"output",
		"o",
		"declaration",
		"Part of function to output, must be one of: declaration, body, signature, comment, return",
	)
}

func countFlagsSet(cmd *cobra.Command) int {
	count := 0
	cmd.Flags().Visit(func(f *pflag.Flag) {
		count++
	})
	return count
}

func displayOutput(function *codescout.FuncNode) {
	switch outputType {
	case "declaration":
		function.PrintNode()
	case "return":
		function.PrintReturnType()
	case "signature":
		function.PrintSignature()
	case "comment":
		function.PrintComments()
	default:
		function.PrintBody()
	}
}

func funcCmdRun(cmd *cobra.Command, args []string) error {
	numFlagsSet := countFlagsSet(cmd)
	filePath := args[0]

	if numFlagsSet == 0 {
		return errors.New("at least one flag must be set for the func command")
	}

	if cmd.Flags().Changed("name") && funcName == "" {
		return errors.New("if name flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("paramtypes") && len(parameterTypes) == 0 {
		return errors.New("if paramtypes flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("return") && len(returnTypes) == 0 {
		return errors.New("if return flag is specified it must not be empty")
	}

	outputAllowed := map[string]*int{
		"declaration": nil,
		"body":        nil,
		"signature":   nil,
		"comment":     nil,
		"return":      nil,
	}
	_, outputValid := outputAllowed[outputType]
	if cmd.Flags().Changed("output") && !outputValid {
		return errors.New("output flag must be one of: declaration, body, signature, comment, return")
	}

	functionTypes := make([]codescout.Parameter, 0, 5)
	for _, parameter := range parameterTypes {
		if strings.Count(parameter, ":") != 1 {
			return errors.New("")
		}

		paramDestruct := strings.SplitN(parameter, ":", 2)
		paramName := strings.TrimSpace(paramDestruct[0])
		paramType := strings.TrimSpace(paramDestruct[1])
		param := codescout.Parameter{Name: paramName, Type: paramType}
		functionTypes = append(functionTypes, param)
	}

	functionConfig := codescout.FuncConfig{
		Name:        funcName,
		Types:       functionTypes,
		ReturnTypes: returnTypes,
	}
	function, err := codescout.ScoutFunction(filePath, functionConfig)
	if err != nil {
		return err
	}
	displayOutput(function)
	return nil
}
