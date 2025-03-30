package codescout

type typeValidation struct {
	TypeMap      map[string]int
	ParameterMap map[string]string
	NameInParams bool
}

func (v *typeValidation) typeNotInParams(name string, paramType string) bool {
	if v.NameInParams && paramType != "" {
		if paramType != v.ParameterMap[name] {
			return true
		}
	}
	return false
}

func (v *typeValidation) hasExhaustedTypes(paramType string) bool {
	v.TypeMap[paramType] -= 1
	return v.TypeMap[paramType] < 0
}

func (v *typeValidation) setNameInParams(name string) {
	v.NameInParams = name != "" && v.isInParams(name)
}

func (v *typeValidation) isInParams(name string) bool {
	_, isInMap := v.ParameterMap[name]
	return isInMap
}

func (v *typeValidation) typeExclusiveNotInParams(paramType string) bool {
	if v.NameInParams && paramType != "" {
		if _, ok := v.TypeMap[paramType]; !ok {
			return true
		}
	}
	return false
}
