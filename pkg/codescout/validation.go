package codescout

type typeValidation struct {
	TypeMap      map[string]int
	ParameterMap map[string]string

	CurParamName    string
	CurParamType    string
	CurNameInParams bool
}

func (v typeValidation) getParamType(name string) string {
	return v.ParameterMap[name]
}

func (v *typeValidation) hasExhausted(paramType string) bool {
	v.TypeMap[paramType] -= 1
	return v.TypeMap[paramType] < 0
}

func (v *typeValidation) setParamName(name string) {
	v.CurParamName = name
}

func (v *typeValidation) setParamType(paramType string) {
	v.CurParamType = paramType
}

func (v *typeValidation) setParamInfo(name string, paramType string) {
	v.setParamName(name)
	v.setParamType(paramType)
}

func (v *typeValidation) setNameInParams(name string) {
	_, isInParams := v.ParameterMap[name]
	v.CurNameInParams = isInParams
}

func (v typeValidation) typeExclExists() bool {
	_, paramTypeExists := v.TypeMap[v.CurParamType]
	return paramTypeExists
}

func (v typeValidation) typeExists() bool {
	if v.CurNameInParams && v.CurParamType != "" {
		return v.CurParamType == v.ParameterMap[v.CurParamName]
	}
	if v.CurParamName == "" {
		return v.typeExclExists()
	}
	return v.CurParamType == ""
}

func returnTypeValidation(returns []string, ops CallableNodeOps) bool {
	returnValidation := typeValidation{TypeMap: ops.returnTypesMap()}
	for _, returnType := range returns {
		returnValidation.setParamType(returnType)
		if !returnValidation.typeExclExists() ||
			returnValidation.hasExhausted(returnType) {
			return false
		}
	}
	return true
}

func parameterTypeValidation(params []Parameter, ops CallableNodeOps) bool {
	typesValidation := typeValidation{
		TypeMap:      ops.parameterTypeMap(),
		ParameterMap: ops.ParametersMap(),
	}
	var parameterType string

	for _, parameter := range params {
		typesValidation.setParamInfo(parameter.Name, parameter.Type)
		typesValidation.setNameInParams(parameter.Name)

		if !typesValidation.CurNameInParams && parameter.Name != "" ||
			!typesValidation.typeExists() {
			return false
		}

		if typesValidation.CurNameInParams {
			parameterType = typesValidation.getParamType(parameter.Name)
		} else {
			parameterType = parameter.Type
		}
		if typesValidation.hasExhausted(parameterType) {
			return false
		}
	}
	return true
}
