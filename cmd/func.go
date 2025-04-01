package cmd

import (
	"errors"

	"github.com/galactixx/codescout/internal/utils"
	codescout "github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
)

var (
	funcName           string
	funcReturnTypes    []string
	funcParameterTypes []string
	funcOutputType     string
)

var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "Find a single function in a file",
	Long:  "Locate and display a specific function definition within a given source file",
	Args:  cobra.ExactArgs(1),
	RunE:  funcCmdRun,
}

func init() {
	rootCmd.AddCommand(funcCmd)

	funcCmd.Flags().StringVarP(&funcName, "name", "n", "", "The function name")
	funcCmd.Flags().StringSliceVarP(&funcParameterTypes, "params", "p", make([]string, 0), "Parameter names and types of function")
	funcCmd.Flags().StringSliceVarP(&funcReturnTypes, "return", "r", make([]string, 0), "Return types of function")
	funcCmd.Flags().StringVarP(
		&funcOutputType,
		"output",
		"o",
		"declaration",
		"Part of function to output, must be one of: declaration, body, signature, comment, return",
	)
}

func displayOutput(function *codescout.FuncNode) {
	switch funcOutputType {
	case "return":
		function.CallableOps.PrintReturnType()
	case "signature":
		function.CallableOps.PrintSignature()
	case "body":
		function.CallableOps.PrintBody()
	case "comment":
		function.PrintComments()
	default:
		function.PrintNode()
	}
}

func funcCmdRun(cmd *cobra.Command, args []string) error {
	numFlagsSet := utils.CountFlagsSet(cmd)
	filePath := args[0]

	if numFlagsSet == 0 {
		return errors.New("at least one flag must be set for the func command")
	}

	if cmd.Flags().Changed("name") && funcName == "" {
		return errors.New("if name flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("paramtypes") && len(funcParameterTypes) == 0 {
		return errors.New("if paramtypes flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("return") && len(funcReturnTypes) == 0 {
		return errors.New("if return flag is specified it must not be empty")
	}

	outputAllowed := map[string]*int{
		"declaration": nil,
		"body":        nil,
		"signature":   nil,
		"comment":     nil,
		"return":      nil,
	}
	_, outputValid := outputAllowed[funcOutputType]
	if cmd.Flags().Changed("output") && !outputValid {
		return errors.New("output flag must be one of: declaration, body, signature, comment, return")
	}

	functionTypes := make([]codescout.Parameter, 0, 5)
	err := utils.ArgsToParams(funcParameterTypes, &functionTypes)
	if err != nil {
		return err
	}

	functionConfig := codescout.FuncConfig{
		Name:        funcName,
		Types:       functionTypes,
		ReturnTypes: funcReturnTypes,
	}
	function, err := codescout.ScoutFunction(filePath, functionConfig)
	if err != nil {
		return err
	}
	displayOutput(function)
	return nil
}
