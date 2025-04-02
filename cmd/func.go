package cmd

import (
	"errors"
	"fmt"

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

var funcEnumOptions = utils.EnumOptions[*codescout.FuncNode]{Options: map[string]func(*codescout.FuncNode) any{
	"declaration": func(node *codescout.FuncNode) any { return node.CallableOps.Code() },
	"body":        func(node *codescout.FuncNode) any { return node.CallableOps.Body() },
	"signature":   func(node *codescout.FuncNode) any { return node.CallableOps.Signature() },
	"comment":     func(node *codescout.FuncNode) any { return node.CallableOps.Comments() },
	"return":      func(node *codescout.FuncNode) any { return node.CallableOps.ReturnType() },
}}

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
		fmt.Sprintf("Part of function to output, must be one of: %v", funcEnumOptions.ToOptionString()),
	)
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

	outputErr := funcEnumOptions.EnumValidation(cmd, "output", funcOutputType)
	if outputErr != nil {
		return outputErr
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
	fmt.Println(funcEnumOptions.GetOutputCallable(funcOutputType)(function))
	return nil
}
