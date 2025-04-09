package flags

import (
	"fmt"
	"log"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
)

func StringBoolToPointer(strBool string) *bool {
	if strBool == "" {
		return nil
	}
	newBool, _ := strconv.ParseBool(strBool)
	return &newBool
}

func StringVarP(cmd *cobra.Command, flag *CommandFlag[string], s string, v string, u string) {
	cmd.Flags().StringVarP(&flag.Variable, flag.Name, s, v, u)
}

func StringSliceVarP(cmd *cobra.Command, flag *CommandFlag[[]string], s string, v []string, u string) {
	cmd.Flags().StringSliceVarP(&flag.Variable, flag.Name, s, v, u)
}

func BoolVarP(cmd *cobra.Command, flag *CommandFlag[bool], s string, v bool, u string) {
	cmd.Flags().BoolVarP(&flag.Variable, flag.Name, s, v, u)
}

type FlagVariable interface {
	bool | string | []string
}

type FlagValidator interface {
	emptyValidator(command *cobra.Command) bool
	isEmptyMessage() error
}

type CommandFlag[T FlagVariable] struct {
	Name     string
	Variable T
}

func (v CommandFlag[T]) emptyValidator(command *cobra.Command) bool {
	reflected := reflect.ValueOf(v.Variable)
	var varIsValid bool
	if reflected.Kind() == reflect.String {
		varIsValid = reflected.String() == ""
	} else if reflected.Kind() == reflect.Slice {
		varIsValid = reflected.Len() == 0
	} else {
		log.Fatal("invalid flag type for validator check")
	}
	return command.Flags().Changed(v.Name) && varIsValid
}

func (v CommandFlag[T]) isEmptyMessage() error {
	return fmt.Errorf("if name %s is specified it must not be empty", v.Name)
}

func (v CommandFlag[T]) stringBoolMessage() error {
	return fmt.Errorf("if %s flag is specified it must be: true or false", v.Name)
}

type BatchValidator struct {
	EmptyValidators      []FlagValidator
	StringBoolValidators []*CommandFlag[string]
}

func (v BatchValidator) batchEmptyValidate(cmd *cobra.Command) error {
	for _, cmdFlag := range v.EmptyValidators {
		if validatorErr := cmdFlag.emptyValidator(cmd); validatorErr {
			return cmdFlag.isEmptyMessage()
		}
	}
	return nil
}

func (v BatchValidator) batchStringBoolValidator(cmd *cobra.Command) error {
	for _, cmdFlag := range v.StringBoolValidators {
		_, ok := map[string]*int{"true": nil, "false": nil}[cmdFlag.Variable]
		if cmd.Flags().Changed(cmdFlag.Name) && !ok {
			return cmdFlag.stringBoolMessage()
		}
	}
	return nil
}

func (v BatchValidator) Validate(cmd *cobra.Command) error {
	stringBoolErr := v.batchStringBoolValidator(cmd)
	if stringBoolErr != nil {
		return stringBoolErr
	}

	emptyErr := v.batchEmptyValidate(cmd)
	if emptyErr != nil {
		return emptyErr
	}

	return nil
}
