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
	methodName           = flags.CommandFlag[string]{Name: "name"}
	methodReturnTypes    = flags.CommandFlag[[]string]{Name: "return"}
	methodParameterTypes = flags.CommandFlag[[]string]{Name: "params"}
	methodOutputType     = flags.CommandFlag[string]{Name: "output"}
	methodReceiver       = flags.CommandFlag[string]{Name: "receiver"}
	hasPointerReceiver   = flags.CommandFlag[string]{Name: "pointer"}
	fieldsAccessed       = flags.CommandFlag[[]string]{Name: "fields"}
	methodsCalled        = flags.CommandFlag[[]string]{Name: "methods"}
	methodNoParams       = flags.CommandFlag[string]{Name: "no-params"}
	methodNoReturn       = flags.CommandFlag[string]{Name: "no-return"}
	noFieldsAccessed     = flags.CommandFlag[string]{Name: "no-fields"}
	noMethodsCalled      = flags.CommandFlag[string]{Name: "no-methods"}
)

var methodEnumOptions = cmdutils.EnumOptions[*codescout.MethodNode]{Options: map[string]func(*codescout.MethodNode) any{
	"definition":       func(node *codescout.MethodNode) any { return node.CallableOps.Code() },
	"body":             func(node *codescout.MethodNode) any { return node.CallableOps.Body() },
	"signature":        func(node *codescout.MethodNode) any { return node.CallableOps.Signature() },
	"comment":          func(node *codescout.MethodNode) any { return node.CallableOps.Comments() },
	"return":           func(node *codescout.MethodNode) any { return node.CallableOps.ReturnType() },
	"receiver":         func(node *codescout.MethodNode) any { return node.ReceiverType() },
	"receiver-fields":  func(node *codescout.MethodNode) any { return node.FieldsAccessed() },
	"receiver-methods": func(node *codescout.MethodNode) any { return node.MethodsCalled() },
}}

var methodBatchValidator = flags.BatchValidator{
	EmptyValidators: []flags.FlagValidator{
		&methodName,
		&methodReceiver,
		&methodParameterTypes,
		&methodReturnTypes,
		&fieldsAccessed,
		&methodsCalled,
	},
	StringBoolValidators: []*flags.CommandFlag[string]{
		&methodNoParams,
		&methodNoReturn,
		&noFieldsAccessed,
		&noMethodsCalled,
		&hasPointerReceiver,
	},
}

var methodCmd = &cobra.Command{
	Use:   "method",
	Short: "Find a single method in a file",
	Long:  "Locate and display a specific method definition within a given source file",
	Args:  cobra.ExactArgs(1),
	RunE:  methodCmdRun,
}

func init() {
	rootCmd.AddCommand(methodCmd)

	flags.StringVarP(methodCmd, &methodName, "n", "", "The method name")
	flags.StringSliceVarP(methodCmd, &methodParameterTypes, "p", make([]string, 0), "Parameter names and types of method")
	flags.StringSliceVarP(methodCmd, &methodReturnTypes, "r", make([]string, 0), "Return types of method")
	flags.StringVarP(methodCmd, &methodReceiver, "v", "", "Receiver type of method")
	flags.StringVarP(methodCmd, &hasPointerReceiver, "t", "", "Whether method has a pointer receiver (true/false)")
	flags.StringSliceVarP(methodCmd, &fieldsAccessed, "f", make([]string, 0), "Struct fields accessed")
	flags.StringSliceVarP(methodCmd, &methodsCalled, "m", make([]string, 0), "Struct methods called")
	flags.StringVarP(methodCmd, &methodNoParams, "s", "", "If the method has no parameters (true/false)")
	flags.StringVarP(methodCmd, &methodNoReturn, "u", "", "If the method has no return type (true/false)")
	flags.StringVarP(methodCmd, &noFieldsAccessed, "d", "", "If the method does not access struct fields (true/false)")
	flags.StringVarP(methodCmd, &noMethodsCalled, "e", "", "If the method does not call struct methods (true/false)")
	flags.StringVarP(
		methodCmd,
		&methodOutputType,
		"o",
		"definition",
		fmt.Sprintf("Part of method to output, must be one of: %v", methodEnumOptions.ToOptionString()),
	)
}

func methodCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if cmdutils.CountFlagsSet(cmd) == 0 {
		return errors.New("at least one flag must be set for the method command")
	}

	validationErr := methodBatchValidator.Validate(cmd)
	if validationErr != nil {
		return validationErr
	}

	methodTypes := make([]codescout.NamedType, 0, 5)
	err := cmdutils.ArgsToParams(methodParameterTypes.Variable, &methodTypes)
	if err != nil {
		return err
	}

	outputErr := methodEnumOptions.EnumValidation(cmd, "output", methodOutputType.Variable)
	if outputErr != nil {
		return outputErr
	}

	methodConfig := codescout.MethodConfig{
		Name:         methodName.Variable,
		ParamTypes:   methodTypes,
		ReturnTypes:  methodReturnTypes.Variable,
		Receiver:     methodReceiver.Variable,
		IsPointerRec: flags.StringBoolToPointer(hasPointerReceiver.Variable),
		Fields:       fieldsAccessed.Variable,
		Methods:      methodsCalled.Variable,
		NoParams:     flags.StringBoolToPointer(methodNoParams.Variable),
		NoReturn:     flags.StringBoolToPointer(methodNoReturn.Variable),
		NoFields:     flags.StringBoolToPointer(noFieldsAccessed.Variable),
		NoMethods:    flags.StringBoolToPointer(noMethodsCalled.Variable),
	}
	method, err := codescout.ScoutMethod(filePath, methodConfig)
	if err != nil {
		return err
	}
	fmt.Println(methodEnumOptions.GetOutputCallable(methodOutputType.Variable)(method))
	return nil
}
