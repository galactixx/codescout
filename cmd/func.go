package cmd

import (
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

var funcOptions = cmdutils.OutputOptions[*codescout.FuncNode]{Options: map[string]func(*codescout.FuncNode) any{
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

var funcCommandValidation = cmdutils.CobraCommandVlidation[*codescout.FuncNode]{
	Validator:      funcBatchValidator,
	NamedTypesFlag: funcParameterTypes,
	OutputTypeFlag: funcOutputType,
	OutputOptions:  funcOptions,
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
		fmt.Sprintf("Part of function to output, must be one of: %v", funcOptions.ToOptionString()),
	)
}

func funcCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	validationErr := funcCommandValidation.CommandValidation(cmd)
	if validationErr != nil {
		return validationErr
	}

	functionConfig := codescout.FuncConfig{
		Name:        funcName.Variable,
		ParamTypes:  funcCommandValidation.GetNamedTypes(),
		ReturnTypes: funcReturnTypes.Variable,
		NoParams:    flags.StringBoolToPointer(funcNoParams.Variable),
		NoReturn:    flags.StringBoolToPointer(funcNoReturn.Variable),
	}
	function, err := codescout.ScoutFunction(filePath, functionConfig)
	if err != nil {
		return err
	}
	fmt.Println(funcOptions.GetOutputCallable(funcOutputType.Variable)(function))
	return nil
}
