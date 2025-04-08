package codescout

import (
	"github.com/galactixx/codescout/internal/pkgutils"
	"github.com/galactixx/codescout/internal/validation"
)

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

type nodeMatchTypes interface{ []string | []NamedType }

type aSTNodeSliceMatch[T, C nodeMatchTypes] struct {
	ConfigTypes  T
	ASTNodeTypes C
	Exact        bool
	NoTypes      *bool
	Validator    func(configTypes T, nodeTypse C) bool
}

func (m aSTNodeSliceMatch[T, C]) nonExactMatch() bool {
	return m.Validator(m.ConfigTypes, m.ASTNodeTypes) && !m.Exact
}

func (m aSTNodeSliceMatch[T, C]) noTypesMatch() bool {
	return (m.NoTypes != nil && *m.NoTypes == (len(m.ASTNodeTypes) == 0)) ||
		(m.NoTypes == nil && len(m.ConfigTypes) == 0)
}

func (m aSTNodeSliceMatch[T, C]) exactMatch() bool {
	return m.Validator(m.ConfigTypes, m.ASTNodeTypes) && m.Exact &&
		len(m.ConfigTypes) == len(m.ASTNodeTypes)
}

func (m aSTNodeSliceMatch[T, C]) validate() bool {
	return m.noTypesMatch() || m.nonExactMatch() || m.exactMatch()
}

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

func accessedMatch(fields []string, nodeFields []string) bool {
	nodeMap := pkgutils.DefaultTypeNilMap(nodeFields)
	for _, field := range fields {
		if _, ok := nodeMap[field]; !ok {
			return false
		}
	}
	return true
}

func namedTypesMapOfTypes(namedTypes []NamedType) map[string]int {
	var parameterTypes []string
	for _, parameter := range namedTypes {
		parameterTypes = append(parameterTypes, parameter.Type)
	}
	return pkgutils.DefaultTypeMap(parameterTypes)
}

func namedTypesMap(namedTypes []NamedType) map[string]string {
	parameters := make(map[string]string)
	for _, parameter := range namedTypes {
		parameters[parameter.Name] = parameter.Type
	}
	return parameters
}
