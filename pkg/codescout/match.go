package codescout

import (
	"github.com/galactixx/codescout/internal/pkgutils"
	"github.com/galactixx/codescout/internal/validation"
)

// astMatch initializes and returns an aSTNodeSliceMatch struct to compare config and AST node types.
func astMatch[T, C nodeMatchTypes](
	configTypes T, nodeTypes C, exact bool, notypes *bool, validator func(configTypes T, nodeTypse C) bool,
) aSTNodeSliceMatch[T, C] {
	return aSTNodeSliceMatch[T, C]{
		ConfigTypes:  configTypes,
		ASTNodeTypes: nodeTypes,
		Exact:        exact,
		NoTypes:      notypes,
		Validator:    validator,
	}
}

// nodeMatchTypes is a constraint for generic type parameters that can be a slice of strings or NamedType.
type nodeMatchTypes interface{ []string | []NamedType }

// aSTNodeSliceMatch holds the data and logic needed to compare config-defined types with AST-derived types.
type aSTNodeSliceMatch[T, C nodeMatchTypes] struct {
	ConfigTypes  T
	ASTNodeTypes C
	Exact        bool
	NoTypes      *bool
	Validator    func(configTypes T, nodeTypse C) bool
}

// nonExactMatch returns true if the validator passes and an exact match is not required.
func (m aSTNodeSliceMatch[T, C]) nonExactMatch() bool {
	return m.Validator(m.ConfigTypes, m.ASTNodeTypes) && !m.Exact
}

// noTypesMatch returns true if type matching is disabled and AST node types are empty,
// or if config types are empty when NoTypes is unset.
func (m aSTNodeSliceMatch[T, C]) noTypesMatch() bool {
	return (m.NoTypes != nil && *m.NoTypes == (len(m.ASTNodeTypes) == 0)) ||
		(m.NoTypes == nil && len(m.ConfigTypes) == 0)
}

// exactMatch returns true if the validator passes and config types match AST types exactly in count.
func (m aSTNodeSliceMatch[T, C]) exactMatch() bool {
	return m.Validator(m.ConfigTypes, m.ASTNodeTypes) && m.Exact &&
		len(m.ConfigTypes) == len(m.ASTNodeTypes)
}

// validate returns true if any of the match conditions (noTypes, non-exact, exact) pass.
func (m aSTNodeSliceMatch[T, C]) validate() bool {
	return m.noTypesMatch() || m.nonExactMatch() || m.exactMatch()
}

// returnMatch validates that all return types exist in the default type map for the given nodeTypes.
func returnMatch(returns []string, nodeTypes []string) bool {
	validation := validation.TypeValidation{TypeMap: pkgutils.DefaultTypeMap(nodeTypes)}
	for _, returnType := range returns {
		validation.SetParamType(returnType)
		if !validation.TypeExclExists() ||
			validation.HasExhausted(returnType) {
			return false
		}
	}
	return true
}

// namedTypesMatch checks whether all named types in config match with those from the AST node types.
func namedTypesMatch(configTypes []NamedType, nodeTypes []NamedType) bool {
	validation := validation.TypeValidation{
		TypeMap:       namedTypesMapOfTypes(nodeTypes),
		NamedTypesMap: namedTypesMap(nodeTypes),
	}
	var parameterType string

	for _, parameter := range configTypes {
		validation.SetParamInfo(parameter.Name, parameter.Type)
		validation.SetNameInParams(parameter.Name)

		if !validation.CurNameInParams && parameter.Name != "" ||
			!validation.TypeExists() {
			return false
		}

		if validation.CurNameInParams {
			parameterType = validation.GetParamType(parameter.Name)
		} else {
			parameterType = parameter.Type
		}
		if validation.HasExhausted(parameterType) {
			return false
		}
	}
	return true
}

// accessedMatch returns true if all config fields are present in the AST node fields.
func accessedMatch(fields []string, nodeFields []string) bool {
	nodeMap := pkgutils.DefaultTypeNilMap(nodeFields)
	for _, field := range fields {
		if _, ok := nodeMap[field]; !ok {
			return false
		}
	}
	return true
}

// namedTypesMapOfTypes creates a map from type string to int from a list of NamedTypes.
func namedTypesMapOfTypes(namedTypes []NamedType) map[string]int {
	var parameterTypes []string
	for _, parameter := range namedTypes {
		parameterTypes = append(parameterTypes, parameter.Type)
	}
	return pkgutils.DefaultTypeMap(parameterTypes)
}

// namedTypesMap returns a map of parameter names to their types.
func namedTypesMap(namedTypes []NamedType) map[string]string {
	parameters := make(map[string]string)
	for _, parameter := range namedTypes {
		parameters[parameter.Name] = parameter.Type
	}
	return parameters
}
