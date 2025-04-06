package cmd

import (
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
	structVerbose    = flags.CommandFlag[bool]{Name: "verbose"}
	structExact      = flags.CommandFlag[bool]{Name: "exact"}
)

var structOptions = cmdutils.OutputOptions[*codescout.StructNode]{Options: map[string]func(*codescout.StructNode) string{
	"definition": func(node *codescout.StructNode) string { return node.Code() },
	"body":       func(node *codescout.StructNode) string { return node.Body() },
	"signature":  func(node *codescout.StructNode) string { return node.Signature() },
	"comment":    func(node *codescout.StructNode) string { return node.Comments() },
}}

var structBatchValidator = flags.BatchValidator{
	EmptyValidators:      []flags.FlagValidator{&structName},
	StringBoolValidators: []*flags.CommandFlag[string]{&structNoFields},
}

var structCommandValidation = cmdutils.CobraCommandVlidation[*codescout.StructNode]{
	Validator:      structBatchValidator,
	NamedTypesFlag: &structFieldTypes,
	OutputTypeFlag: &structOutputType,
	OutputOptions:  structOptions,
}

var structCmd = &cobra.Command{
	Use:   "struct",
	Short: "Find a single struct in a file",
	Long:  "Locate and display a specific struct definition within a given source file",
	RunE:  structCmdRun,
}

func init() {
	rootCmd.AddCommand(structCmd)

	flags.StringVarP(structCmd, &structName, "n", "", "the struct name")
	flags.StringSliceVarP(structCmd, &structFieldTypes, "f", make([]string, 0), "field names and types of struct")
	flags.StringVarP(structCmd, &structNoFields, "s", "", "if the struct has no fields (true/false)")
	flags.BoolVarP(structCmd, &structVerbose, "v", false, "whether to print all occurrences or just the first")
	flags.BoolVarP(structCmd, &structExact, "x", false, "if an exact match should occur with slice flags")
	flags.StringVarP(
		structCmd,
		&structOutputType,
		"o",
		"definition",
		fmt.Sprintf("part of struct to output, must be one of: %v", structOptions.ToOptionString()),
	)
}

func structCmdRun(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	validationErr := structCommandValidation.CommandValidation(cmd)
	if validationErr != nil {
		return validationErr
	}

	structConfig := codescout.StructConfig{
		Name:       structName.Variable,
		FieldTypes: structCommandValidation.GetNamedTypes(),
		NoFields:   flags.StringBoolToPointer(structNoFields.Variable),
		Exact:      structExact.Variable,
	}
	scoutContainer := cmdutils.NewScoutContainer(
		codescout.ScoutStruct,
		codescout.ScoutStructs,
		filePath,
		structOptions,
		structConfig,
		"Struct",
		structOutputType.Variable,
	)
	return scoutContainer.Display(structVerbose.Variable)
}
