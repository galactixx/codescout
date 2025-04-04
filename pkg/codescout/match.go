package codescout

import (
	"github.com/galactixx/codescout/internal/validation"
)

func fullReturnMatch(types []string, boolVar *bool, ops CallableOps) bool {
	return boolMatch(boolVar, ops.ReturnTypes()) &&
		returnMatch(types, ops)
}

func fullParamsMatch(types []NamedType, boolVar *bool, ops CallableOps) bool {
	return boolMatch(boolVar, ops.Parameters()) &&
		parameterMatch(types, ops)
}

func returnMatch(returns []string, ops CallableOps) bool {
	validation := validation.TypeValidation{TypeMap: ops.returnTypesMap()}
	for _, returnType := range returns {
		validation.SetParamType(returnType)
		if !validation.TypeExclExists() ||
			validation.HasExhausted(returnType) {
			return false
		}
	}
	return true
}

func parameterMatch(params []NamedType, ops CallableOps) bool {
	validation := validation.TypeValidation{
		TypeMap:      ops.parameterTypeMap(),
		ParameterMap: ops.ParametersMap(),
	}
	var parameterType string

	for _, parameter := range params {
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

func accessedMatch(fields []string, nodeMap map[string]*int) bool {
	for _, field := range fields {
		if _, ok := nodeMap[field]; !ok {
			return false
		}
	}
	return true
}

func fullAccessedMatch(fields []string, boolVar *bool, node MethodNode) bool {
	return boolMatch(boolVar, node.FieldsAccessed()) &&
		accessedMatch(fields, node.fieldsAccessed)
}

func fullCalledMatch(methods []string, boolVar *bool, node MethodNode) bool {
	return boolMatch(boolVar, node.MethodsCalled()) &&
		accessedMatch(methods, node.methodsCalled)
}

type CallableTypes interface {
	~[]string | ~[]NamedType
}

func boolMatch[T CallableTypes](boolVar *bool, nodeTypes T) bool {
	return boolVar == nil || (*boolVar == (len(nodeTypes) == 0))
}
