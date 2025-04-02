package cmd

import (
	"errors"
	"fmt"

	"github.com/galactixx/codescout/internal/utils"
	codescout "github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
)

var (
	methodName           string
	methodReturnTypes    []string
	methodParameterTypes []string
	methodOutputType     string
	methodReceiver       string
	hasPointerReceiver   string
)

var methodEnumOptions = utils.EnumOptions[*codescout.MethodNode]{Options: map[string]func(*codescout.MethodNode) any{
	"declaration":      func(node *codescout.MethodNode) any { return node.CallableOps.Code() },
	"body":             func(node *codescout.MethodNode) any { return node.CallableOps.Body() },
	"signature":        func(node *codescout.MethodNode) any { return node.CallableOps.Signature() },
	"comment":          func(node *codescout.MethodNode) any { return node.CallableOps.Comments() },
	"return":           func(node *codescout.MethodNode) any { return node.CallableOps.ReturnType() },
	"receiver":         func(node *codescout.MethodNode) any { return node.ReceiverType() },
	"receiver-fields":  func(node *codescout.MethodNode) any { return node.FieldsAccessed() },
	"recevier-methods": func(node *codescout.MethodNode) any { return node.MethodsCalled() },
}}

var methodCmd = &cobra.Command{
	Use:   "method",
	Short: "Find a single method in a file",
	Long:  "Locate and display a specific method definition within a given source file",
	Args:  cobra.ExactArgs(1),
	RunE:  methodCmdRun,
}

func init() {
	rootCmd.AddCommand(methodCmd)

	methodCmd.Flags().StringVarP(&methodName, "name", "n", "", "The method name")
	methodCmd.Flags().StringSliceVarP(&methodParameterTypes, "params", "p", make([]string, 0), "Parameter names and types of method")
	methodCmd.Flags().StringSliceVarP(&methodReturnTypes, "return", "r", make([]string, 0), "Return types of method")
	methodCmd.Flags().StringVarP(&methodReceiver, "receiver", "v", "", "Receiver type of method")
	methodCmd.Flags().StringVarP(&hasPointerReceiver, "pointer", "t", "", "Whether method has a pointer receiver (true/false)")
	methodCmd.Flags().StringVarP(
		&methodOutputType,
		"output",
		"o",
		"declaration",
		fmt.Sprintf("Part of method to output, must be one of: %v", methodEnumOptions.ToOptionString()),
	)
}

func methodCmdRun(cmd *cobra.Command, args []string) error {
	numFlagsSet := utils.CountFlagsSet(cmd)
	filePath := args[0]

	if numFlagsSet == 0 {
		return errors.New("at least one flag must be set for the func command")
	}

	if cmd.Flags().Changed("name") && methodName == "" {
		return errors.New("if name flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("paramtypes") && len(methodParameterTypes) == 0 {
		return errors.New("if paramtypes flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("return") && len(methodReturnTypes) == 0 {
		return errors.New("if return flag is specified it must not be empty")
	}

	if cmd.Flags().Changed("receiver") && methodReceiver == "" {
		return errors.New("if receiver flag is specified it must not be empty")
	}

	_, hasPointer := map[string]*int{"true": nil, "false": nil}[hasPointerReceiver]
	if cmd.Flags().Changed("pointer") && !hasPointer {
		return errors.New("if pointer flag is specified it must be: true or false")
	}

	methodTypes := make([]codescout.Parameter, 0, 5)
	err := utils.ArgsToParams(methodParameterTypes, &methodTypes)
	if err != nil {
		return err
	}

	outputErr := methodEnumOptions.EnumValidation(cmd, "output", methodOutputType)
	if outputErr != nil {
		return outputErr
	}

	methodConfig := codescout.MethodConfig{
		Name:         methodName,
		Types:        methodTypes,
		ReturnTypes:  methodReturnTypes,
		Receiver:     methodReceiver,
		IsPointerRec: hasPointerReceiver,
	}
	method, err := codescout.ScoutMehod(filePath, methodConfig)
	if err != nil {
		return err
	}
	fmt.Println(methodEnumOptions.GetOutputCallable(methodOutputType)(method))
	return nil
}
