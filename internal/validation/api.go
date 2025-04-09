package validation

import "fmt"

type Argument[T any] struct {
	Name     string
	Variable T
}

func Arg[T any](name string, variable T) Argument[T] {
	return Argument[T]{Name: name, Variable: variable}
}

type SliceValidator interface {
	trueMessage() error
	falseMessage() error
	boolNotNil() bool
	trueValidate() bool
	falseValidate() bool
	nonEmptySlice() bool
}

type SlicePairToValidate[T any] struct {
	Slice Argument[[]T]
	Bool  Argument[*bool]
}

func (p SlicePairToValidate[T]) trueMessage() error {
	return fmt.Errorf("%s cannot be specified if %s is set to true", p.Slice.Name, p.Bool.Name)
}

func (p SlicePairToValidate[T]) falseMessage() error {
	return fmt.Errorf("no need to specify %s if %s is set to false", p.Slice.Name, p.Bool.Name)
}

func (p SlicePairToValidate[T]) boolNotNil() bool {
	return p.Bool.Variable != nil
}

func (p SlicePairToValidate[T]) nonEmptySlice() bool {
	return len(p.Slice.Variable) > 0
}

func (p SlicePairToValidate[T]) trueValidate() bool {
	return p.nonEmptySlice() && *p.Bool.Variable
}

func (p SlicePairToValidate[T]) falseValidate() bool {
	return p.nonEmptySlice() && !*p.Bool.Variable
}

type BatchConfigValidation struct {
	SliceValidators []SliceValidator
	Exact           bool
}

func (p BatchConfigValidation) exactMessage() error {
	return fmt.Errorf("exact should not be true if no slices are passed")
}

func (v BatchConfigValidation) Validate() error {
	existingNonEmptySlice := false
	for _, validator := range v.SliceValidators {
		if validator.boolNotNil() {
			if validator.trueValidate() {
				return validator.trueMessage()
			}

			if validator.falseValidate() {
				return validator.falseMessage()
			}
		}

		if !existingNonEmptySlice && validator.nonEmptySlice() {
			existingNonEmptySlice = true
		}
	}

	if !existingNonEmptySlice && v.Exact {
		return v.exactMessage()
	}

	return nil
}
