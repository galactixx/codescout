package validation

type TypeValidation struct {
	TypeMap      map[string]int
	ParameterMap map[string]string

	CurParamName    string
	CurParamType    string
	CurNameInParams bool
}

func (v TypeValidation) GetParamType(name string) string {
	return v.ParameterMap[name]
}

func (v *TypeValidation) HasExhausted(paramType string) bool {
	v.TypeMap[paramType] -= 1
	return v.TypeMap[paramType] < 0
}

func (v *TypeValidation) SetParamName(name string) {
	v.CurParamName = name
}

func (v *TypeValidation) SetParamType(paramType string) {
	v.CurParamType = paramType
}

func (v *TypeValidation) SetParamInfo(name string, paramType string) {
	v.SetParamName(name)
	v.SetParamType(paramType)
}

func (v *TypeValidation) SetNameInParams(name string) {
	_, isInParams := v.ParameterMap[name]
	v.CurNameInParams = isInParams
}

func (v TypeValidation) TypeExclExists() bool {
	_, paramTypeExists := v.TypeMap[v.CurParamType]
	return paramTypeExists
}

func (v TypeValidation) TypeExists() bool {
	if v.CurNameInParams && v.CurParamType != "" {
		return v.CurParamType == v.ParameterMap[v.CurParamName]
	}
	if v.CurParamName == "" {
		return v.TypeExclExists()
	}
	return v.CurParamType == ""
}
