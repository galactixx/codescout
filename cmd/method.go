package cmd

import (
	"errors"

	"github.com/galactixx/codescout/internal/utils"
	codescout "github.com/galactixx/codescout/pkg/codescout"
	"github.com/spf13/cobra"
)

var (
	methodName           string
	methodReturnTypes    []string
	methodParameterTypes []string
)

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

	methodTypes := make([]codescout.Parameter, 0, 5)
	err := utils.ArgsToParams(methodParameterTypes, &methodTypes)
	if err != nil {
		return err
	}

	methodConfig := codescout.MethodConfig{
		Name:        methodName,
		Types:       methodTypes,
		ReturnTypes: funcReturnTypes,
	}
	method, err := codescout.ScoutMehod(filePath, methodConfig)
	if err != nil {
		return err
	}
	method.PrintNode()
	return nil
}
