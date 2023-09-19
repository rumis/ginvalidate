package ginvalidate

type ParamsCollection map[string]interface{}

func NewParamsCollection() ParamsCollection {
	return make(ParamsCollection)
}

// Set 设置值
func (pc ParamsCollection) Set(k string, v []string) {
	ev, ok := pc[k]
	if !ok {
		if len(v) == 1 {
			pc[k] = v[0] // 如果参数为单值，则直接赋值一个字符串
			return
		}
		pc[k] = v // 如果参数为多值，则直接赋值为数组
		return
	}
	switch eVal := ev.(type) {
	case string:
		if len(v) == 1 {
			pc[k] = []string{eVal, v[0]}
			return
		}
		pc[k] = append([]string{eVal}, v...)
	case []string:
		pc[k] = append(eVal, v...)
	}
}

// To 返回map格式对象
func (pc ParamsCollection) To() map[string]interface{} {
	return pc
}
