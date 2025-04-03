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
	TrueMessage() error
	FalseMessage() error
	BoolNotNil() bool
	TrueValidate() bool
	FalseValidate() bool
}

type SlicePairToValidate[T any] struct {
	Slice Argument[[]T]
	Bool  Argument[*bool]
}

func (p SlicePairToValidate[T]) TrueMessage() error {
	return fmt.Errorf("%v cannot be specified if %v is set to true", p.Slice.Name, p.Bool.Name)
}

func (p SlicePairToValidate[T]) FalseMessage() error {
	return fmt.Errorf("no need to specify %v if %v is set to false", p.Slice.Name, p.Bool.Name)
}

func (p SlicePairToValidate[T]) BoolNotNil() bool {
	return p.Bool.Variable != nil
}

func (p SlicePairToValidate[T]) TrueValidate() bool {
	return len(p.Slice.Variable) > 0 && *p.Bool.Variable
}

func (p SlicePairToValidate[T]) FalseValidate() bool {
	return len(p.Slice.Variable) > 0 && !*p.Bool.Variable
}

type BatchConfigValidation struct {
	SliceValidators []SliceValidator
}

func (v BatchConfigValidation) Validate() error {
	for _, validator := range v.SliceValidators {
		if validator.BoolNotNil() {
			if validator.TrueValidate() {
				return validator.TrueMessage()
			}

			if validator.FalseValidate() {
				return validator.FalseMessage()
			}
		}
	}
	return nil
}
