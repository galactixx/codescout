package cmd

import (
	"errors"
	"fmt"

	"github.com/galactixx/codescout/internal/cmdutils"
	"github.com/galactixx/codescout/internal/flags"
	codescout "github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
)

var (
	funcName           = flags.CommandFlag[string]{Name: "name"}
	funcReturnTypes    = flags.CommandFlag[[]string]{Name: "return"}
	funcParameterTypes = flags.CommandFlag[[]string]{Name: "params"}
	funcOutputType     = flags.CommandFlag[string]{Name: "output"}
	funcNoParams       = flags.CommandFlag[string]{Name: "no-params"}
	funcNoReturn       = flags.CommandFlag[string]{Name: "no-return"}
)

var funcEnumOptions = cmdutils.EnumOptions[*codescout.FuncNode]{Options: map[string]func(*codescout.FuncNode) any{
	"definition": func(node *codescout.FuncNode) any { return node.CallableOps.Code() },
	"body":       func(node *codescout.FuncNode) any { return node.CallableOps.Body() },
	"signature":  func(node *codescout.FuncNode) any { return node.CallableOps.Signature() },
	"comment":    func(node *codescout.FuncNode) any { return node.CallableOps.Comments() },
	"return":     func(node *codescout.FuncNode) any { return node.CallableOps.ReturnType() },
}}

var funcBatchValidator = flags.BatchValidator{
	EmptyValidators: []flags.FlagValidator{
		&funcName,
		&funcParameterTypes,
		&funcReturnTypes,
	},
	StringBoolValidators: []*flags.CommandFlag[string]{&funcNoParams, &funcNoReturn},
}

var funcCmd = &cobra.Command{
	Use:   "func",
	Short: "Find a single function in a file",
	Long:  "Locate and display a specific function definition within a given source file",
	Args:  cobra.ExactArgs(1),
	RunE:  funcCmdRun,
}

func init() {
	rootCmd.AddCommand(funcCmd)

	flags.StringVarP(funcCmd, &funcName, "n", "", "The function name")
	flags.StringSliceVarP(funcCmd, &funcParameterTypes, "p", make([]string, 0), "Parameter names and types of function")
	flags.StringSliceVarP(funcCmd, &funcReturnTypes, "r", make([]string, 0), "Return types of function")
	flags.StringVarP(funcCmd, &funcNoParams, "s", "", "If the function has no parameters (true/false)")
	flags.StringVarP(funcCmd, &funcNoReturn, "u", "", "If the function has no return type (true/false)")
	flags.StringVarP(
		funcCmd,
		&funcOutputType,
		"o",
		"definition",
		fmt.Sprintf("Part of function to output, must be one of: %v", funcEnumOptions.ToOptionString()),
	)
}

func funcCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if cmdutils.CountFlagsSet(cmd) == 0 {
		return errors.New("at least one flag must be set for the func command")
	}

	validationErr := funcBatchValidator.Validate(cmd)
	if validationErr != nil {
		return validationErr
	}

	functionTypes := make([]codescout.NamedType, 0, 5)
	err := cmdutils.ArgsToParams(funcParameterTypes.Variable, &functionTypes)
	if err != nil {
		return err
	}

	outputErr := funcEnumOptions.EnumValidation(cmd, "output", funcOutputType.Variable)
	if outputErr != nil {
		return outputErr
	}

	functionConfig := codescout.FuncConfig{
		Name:        funcName.Variable,
		ParamTypes:  functionTypes,
		ReturnTypes: funcReturnTypes.Variable,
		NoParams:    flags.StringBoolToPointer(funcNoParams.Variable),
		NoReturn:    flags.StringBoolToPointer(funcNoReturn.Variable),
	}
	function, err := codescout.ScoutFunction(filePath, functionConfig)
	if err != nil {
		return err
	}
	fmt.Println(funcEnumOptions.GetOutputCallable(funcOutputType.Variable)(function))
	return nil
}
