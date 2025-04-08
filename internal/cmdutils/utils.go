package cmdutils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/galactixx/codescout/internal/flags"
	"github.com/galactixx/codescout/pkg/codescout"
	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func JoinAttrs(attrs []string) string { return fmt.Sprintf("[ %v ]", strings.Join(attrs, ",")) }
func stripANSI(input string) string   { return ansiRegexp.ReplaceAllString(input, "") }

type OutputOptions[T any] struct {
	Options map[string]func(T) string
}

func (o OutputOptions[T]) validation(cmd *cobra.Command, flag flags.CommandFlag[string]) error {
	_, outputValid := o.Options[flag.Variable]
	if cmd.Flags().Changed("output") && !outputValid {
		return fmt.Errorf("%v flag must be one of: %v", flag.Name, o.ToOptionString())
	}
	return nil
}

func (o OutputOptions[T]) ToOptionString() string {
	optionsSlice := make([]string, 0, 5)
	for option := range o.Options {
		optionsSlice = append(optionsSlice, option)
	}
	return strings.Join(optionsSlice, ", ")
}

func (o OutputOptions[T]) getOutputCallable(option string) func(T) string {
	return o.Options[option]
}

func argsToNamedTypes(argTypes []string, parameterTypes *[]codescout.NamedType) error {
	for _, parameter := range argTypes {
		if strings.Count(parameter, ":") != 1 {
			return errors.New("there must be only one colon separating out the name and type")
		}

		paramDestruct := strings.SplitN(parameter, ":", 2)
		paramName := strings.TrimSpace(paramDestruct[0])
		paramType := strings.TrimSpace(paramDestruct[1])

		if paramName == "" && paramType == "" {
			return errors.New("at least one of the type or name must be defined")
		}
		param := codescout.NamedType{Name: paramName, Type: paramType}
		*parameterTypes = append(*parameterTypes, param)
	}
	return nil
}

type CobraCommandVlidation[T any] struct {
	Validator      flags.BatchValidator
	NamedTypesFlag *flags.CommandFlag[[]string]
	OutputTypeFlag *flags.CommandFlag[string]
	OutputOptions  OutputOptions[T]

	namedTypes []codescout.NamedType
}

func (v *CobraCommandVlidation[T]) GetNamedTypes() []codescout.NamedType {
	if v.namedTypes == nil {
		log.Fatal("named types field is returning nil, should never occur")
	}
	namedTypes := v.namedTypes
	v.namedTypes = nil
	return namedTypes
}

func (v *CobraCommandVlidation[T]) CommandValidation(cmd *cobra.Command) error {
	validationErr := v.Validator.Validate(cmd)
	if validationErr != nil {
		return validationErr
	}

	namedTypes := make([]codescout.NamedType, 0, 5)
	err := argsToNamedTypes(v.NamedTypesFlag.Variable, &namedTypes)
	if err != nil {
		return err
	}
	v.namedTypes = namedTypes

	outputErr := v.OutputOptions.validation(cmd, *v.OutputTypeFlag)
	if outputErr != nil {
		return outputErr
	}
	return nil
}

func NewScoutContainer[T any, C any](
	scoutFirst func(path string, config C) (*T, error),
	scoutAll func(path string, config C) ([]*T, error),
	path string,
	options OutputOptions[*T],
	config C,
	defType string,
	outputType string,
) ScoutContainer[T, C] {
	var boxWidth int
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		boxWidth = 80
	} else {
		boxWidth = width / 2
	}

	separatorSection := color.New(color.FgHiBlack, color.Bold)
	return ScoutContainer[T, C]{
		ScoutFirst:     scoutFirst,
		ScoutAll:       scoutAll,
		Path:           path,
		Options:        options,
		BoxWidth:       boxWidth,
		SeparatorColor: separatorSection,
		Config:         config,
		DefType:        defType,
		OutputType:     outputType,
	}
}

type ScoutContainer[T any, C any] struct {
	ScoutFirst     func(path string, config C) (*T, error)
	ScoutAll       func(path string, config C) ([]*T, error)
	Path           string
	Options        OutputOptions[*T]
	BoxWidth       int
	SeparatorColor *color.Color
	Config         C
	DefType        string
	OutputType     string
}

func (c ScoutContainer[T, C]) Display(verbose bool) error {
	if verbose {
		occurrences, err := c.ScoutAll(c.Path, c.Config)
		if err != nil {
			return err
		}
		for idx, occurrence := range occurrences {
			name := getNameFromNodes(occurrence)
			c.printOutput(
				name,
				idx == 0,
				c.Options.getOutputCallable(c.OutputType)(occurrence),
			)
		}
	} else {
		occurrence, err := c.ScoutFirst(c.Path, c.Config)
		if err != nil {
			return err
		}
		name := getNameFromNodes(occurrence)
		c.printOutput(
			name,
			true,
			c.Options.getOutputCallable(c.OutputType)(occurrence),
		)
	}
	return nil
}

func (c ScoutContainer[T, C]) displaySeparator(separator string) {
	c.SeparatorColor.Println(separator)
}

func (c ScoutContainer[T, C]) constructHeader(name string, boxWidth int) string {
	fieldNameColor := color.New(color.FgCyan, color.Bold, color.Underline)
	fieldValueColor := color.New(color.FgCyan)

	header := fmt.Sprintf(
		"%s: %s   |   %s: %s   |   %s: %s",
		fieldNameColor.Sprint("Type"),
		fieldValueColor.Sprint(c.DefType),
		fieldNameColor.Sprint("Name"),
		fieldValueColor.Sprint(name),
		fieldNameColor.Sprint("Output"),
		fieldValueColor.Sprint(capitalizeString(c.OutputType)),
	)

	visibleWidth := runewidth.StringWidth(stripANSI(header))
	padding := boxWidth - 4 - visibleWidth
	paddedHeader := header + strings.Repeat(" ", padding)
	return paddedHeader
}

func (c ScoutContainer[T, C]) printOutput(name string, showSeparator bool, output string) {
	boxWidth := getMax(c.BoxWidth, findLengthOfOutput(output))
	separator := strings.Repeat("═", boxWidth)

	if showSeparator {
		c.displaySeparator(separator)
	}
	boxOuterLine := strings.Repeat("─", boxWidth-2)
	header := c.constructHeader(name, boxWidth)

	titleSection := color.New(color.FgCyan, color.Bold)

	titleSection.Println("╭" + boxOuterLine + "╮")
	titleSection.Print("│ ")
	titleSection.Print(header)
	titleSection.Println(" │")
	titleSection.Println("╰" + boxOuterLine + "╯")

	codeSection := color.New(color.FgWhite, color.Bold)
	codeBorders := color.New(color.FgGreen, color.Bold)

	wrapped := wordwrap.WrapString(output, uint(boxWidth-4))

	codeBorders.Println("╭" + boxOuterLine + "╮")
	for _, line := range strings.Split(wrapped, "\n") {
		codeBorders.Print("│ ")
		codeSection.Printf("%-*s", boxWidth-4, strings.Replace(line, "\t", "    ", -1))
		codeBorders.Println(" │")
	}
	codeBorders.Println("╰" + boxOuterLine + "╯")
	c.displaySeparator(separator)
}
