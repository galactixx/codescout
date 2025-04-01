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
