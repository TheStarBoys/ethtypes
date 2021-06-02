package ethtypes

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func VariableToStr(v StateVariable) string {
	val := reflect.New(v.Type())
	v.Get(val.Interface())

	return fmt.Sprintf("%s 0x%x => %v", v.Name(), v.Addr(), val.Elem().Interface())
}

func GetArrayElems(a Array) []interface{} {
	var res []interface{}
	for i := 0; i < a.Len(); i++ {
		val := reflect.New(a.ElemType()).Interface()
		a.Get(i, val)
		res = append(res, reflect.ValueOf(val).Elem().Interface())
	}

	return res
}

func ArrayToStr(a Array) string {
	return fmt.Sprintf("%s len: %d %v", a.Name(), a.Len(), GetArrayElems(a))
}

func SliceToStr(s Slice) string {
	var res string
	res += fmt.Sprintf("%s len: %d, cap: %d ", s.Name(), s.Len(), s.Cap())
	res += ArrayToStr(s)

	return res
}

func GetIterableMapElems(m IterableMap) map[string]interface{} {
	res := make(map[string]interface{})
	m.Range(func(key, val interface{}) bool {
		js, _ := json.Marshal(key)
		res[string(js)] = val

		return true
	})

	return res
}

func IterableMapToStr(m IterableMap) string {
	return fmt.Sprintf("%s len: %d %v", m.Name(), m.Len(), GetIterableMapElems(m))
}
