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
	structName       = flags.CommandFlag[string]{Name: "name"}
	structOutputType = flags.CommandFlag[string]{Name: "output"}
	structFieldTypes = flags.CommandFlag[[]string]{Name: "fields"}
	structNoFields   = flags.CommandFlag[string]{Name: "no-fields"}
)

var structEnumOptions = cmdutils.EnumOptions[*codescout.StructNode]{Options: map[string]func(*codescout.StructNode) any{
	"definition": func(node *codescout.StructNode) any { return node.Code() },
	"body":       func(node *codescout.StructNode) any { return node.Body() },
	"signature":  func(node *codescout.StructNode) any { return node.Signature() },
	"comment":    func(node *codescout.StructNode) any { return node.Comments() },
	"methods":    func(node *codescout.StructNode) any { return "" },
}}

var structBatchValidator = flags.BatchValidator{
	EmptyValidators: []flags.FlagValidator{
		&structName,
	},
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Find a single struct in a file",
	Long:  "Locate and display a specific struct definition within a given source file",
	RunE:  structCmdRun,
}

func init() {
	rootCmd.AddCommand(structCmd)

	flags.StringVarP(structCmd, &structName, "n", "", "The struct name")
	flags.StringSliceVarP(structCmd, &structFieldTypes, "f", make([]string, 0), "Field names and types of struct")
	flags.StringVarP(structCmd, &structNoFields, "s", "", "If the struct has no fields (true/false)")
	flags.StringVarP(
		structCmd,
		&structOutputType,
		"o",
		"definition",
		fmt.Sprintf("Part of struct to output, must be one of: %v", structEnumOptions.ToOptionString()),
	)
}

func structCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if cmdutils.CountFlagsSet(cmd) == 0 {
		return errors.New("at least one flag must be set for the struct command")
	}

	validationErr := structBatchValidator.Validate(cmd)
	if validationErr != nil {
		return validationErr
	}

	outputErr := structEnumOptions.EnumValidation(cmd, "output", structOutputType.Variable)
	if outputErr != nil {
		return outputErr
	}

	structConfig := codescout.StructConfig{
		Name: structName.Variable,
	}
	structure, err := codescout.ScoutStruct(filePath, structConfig)
	if err != nil {
		return err
	}
	fmt.Println(structure.Comments())
	fmt.Println(structEnumOptions.GetOutputCallable(structOutputType.Variable)(structure))
	return nil
}
