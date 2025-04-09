package cmd

import (
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
	methodVerbose        = flags.CommandFlag[bool]{Name: "verbose"}
	methodExact          = flags.CommandFlag[bool]{Name: "exact"}
)

var methodOptions = cmdutils.OutputOptions[*codescout.MethodNode]{Options: map[string]func(*codescout.MethodNode) string{
	"definition":       func(node *codescout.MethodNode) string { return node.CallableOps.Code() },
	"body":             func(node *codescout.MethodNode) string { return node.CallableOps.Body() },
	"signature":        func(node *codescout.MethodNode) string { return node.CallableOps.Signature() },
	"comment":          func(node *codescout.MethodNode) string { return node.CallableOps.Comments() },
	"return":           func(node *codescout.MethodNode) string { return node.CallableOps.ReturnType() },
	"receiver":         func(node *codescout.MethodNode) string { return node.ReceiverType() },
	"receiver-fields":  func(node *codescout.MethodNode) string { return cmdutils.JoinAttrs(node.FieldsAccessed()) },
	"receiver-methods": func(node *codescout.MethodNode) string { return cmdutils.JoinAttrs(node.MethodsCalled()) },
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

var methodCommandValidation = cmdutils.CobraCommandVlidation[*codescout.MethodNode]{
	Validator:      methodBatchValidator,
	NamedTypesFlag: &methodParameterTypes,
	OutputTypeFlag: &methodOutputType,
	OutputOptions:  methodOptions,
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

	flags.StringVarP(methodCmd, &methodName, "n", "", "the method name")
	flags.StringSliceVarP(methodCmd, &methodParameterTypes, "p", make([]string, 0), "parameter names and types of method")
	flags.StringSliceVarP(methodCmd, &methodReturnTypes, "r", make([]string, 0), "return types of method")
	flags.StringVarP(methodCmd, &methodReceiver, "m", "", "receiver type of method")
	flags.StringVarP(methodCmd, &hasPointerReceiver, "t", "", "whether method has a pointer receiver (true/false)")
	flags.StringSliceVarP(methodCmd, &fieldsAccessed, "f", make([]string, 0), "struct fields accessed")
	flags.StringSliceVarP(methodCmd, &methodsCalled, "c", make([]string, 0), "struct methods called")
	flags.StringVarP(methodCmd, &methodNoParams, "s", "", "if the method has no parameters (true/false)")
	flags.StringVarP(methodCmd, &methodNoReturn, "u", "", "if the method has no return type (true/false)")
	flags.StringVarP(methodCmd, &noFieldsAccessed, "d", "", "if the method does not access struct fields (true/false)")
	flags.StringVarP(methodCmd, &noMethodsCalled, "e", "", "if the method does not call struct methods (true/false)")
	flags.BoolVarP(methodCmd, &methodVerbose, "v", false, "whether to print all occurrences or just the first (true/false)")
	flags.BoolVarP(methodCmd, &methodExact, "x", false, "if an exact match should occur with slice flags (true/false)")
	flags.StringVarP(
		methodCmd,
		&methodOutputType,
		"o",
		"definition",
		fmt.Sprintf("part of method to output, must be one of: %s", methodOptions.ToOptionString()),
	)
}

func methodCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	validationErr := methodCommandValidation.CommandValidation(cmd)
	if validationErr != nil {
		return validationErr
	}

	methodConfig := codescout.MethodConfig{
		Name:         methodName.Variable,
		ParamTypes:   methodCommandValidation.GetNamedTypes(),
		ReturnTypes:  methodReturnTypes.Variable,
		Receiver:     methodReceiver.Variable,
		IsPointerRec: flags.StringBoolToPointer(hasPointerReceiver.Variable),
		Fields:       fieldsAccessed.Variable,
		Methods:      methodsCalled.Variable,
		NoParams:     flags.StringBoolToPointer(methodNoParams.Variable),
		NoReturn:     flags.StringBoolToPointer(methodNoReturn.Variable),
		NoFields:     flags.StringBoolToPointer(noFieldsAccessed.Variable),
		NoMethods:    flags.StringBoolToPointer(noMethodsCalled.Variable),
		Exact:        methodExact.Variable,
	}
	scoutContainer := cmdutils.NewScoutContainer(
		codescout.ScoutMethod,
		codescout.ScoutMethods,
		filePath,
		methodOptions,
		methodConfig,
		"Method",
		methodOutputType.Variable,
	)
	return scoutContainer.Display(methodVerbose.Variable)
}
