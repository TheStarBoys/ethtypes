package ethtypes

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type BasicMap struct {
	state   *ContractState
	name    string
	keyType reflect.Type
	valType reflect.Type
}

const (
	mapPrefix    = "map_"
	mapKeyPrefix = "map_key_"
)

func NewBasicMap(state *ContractState, name string, keyType, valType reflect.Type) (*BasicMap, error) {
	for count := 3; count > 0 && keyType.Kind() == reflect.Ptr; count-- {
		valType = keyType.Elem()
	}

	for count := 3; count > 0 && valType.Kind() == reflect.Ptr; count-- {
		valType = valType.Elem()
	}

	if keyType.Kind() == reflect.Ptr || valType.Kind() == reflect.Ptr {
		panic("cannot set pointer key or val")
	}

	m := &BasicMap{
		state:   state,
		name:    name,
		keyType: keyType,
		valType: valType,
	}

	return m, nil
}

func GetBasicMap(state *ContractState, name string, keyType, valType reflect.Type) (*BasicMap, error) {
	for count := 3; count > 0 && keyType.Kind() == reflect.Ptr; count-- {
		valType = keyType.Elem()
	}

	for count := 3; count > 0 && valType.Kind() == reflect.Ptr; count-- {
		valType = valType.Elem()
	}

	if keyType.Kind() == reflect.Ptr || valType.Kind() == reflect.Ptr {
		panic("cannot set pointer key or val")
	}

	m := &BasicMap{
		state:   state,
		name:    name,
		keyType: keyType,
		valType: valType,
	}

	return m, nil
}

func (m *BasicMap) Name() string {
	return m.name
}

func (m *BasicMap) Get(key interface{}, val interface{}) bool {
	if actual, expect := reflect.TypeOf(key).Kind(), m.keyType.Kind(); actual != expect {
		panic(fmt.Sprintf("key not match, actual: %v, expect: %v", actual, expect))
	}

	if !m.Contains(key) {
		return false
	}

	m.getElem(key).Get(val)
	return true
}

func (m *BasicMap) Set(key, val interface{}) {
	if actual, expect := reflect.TypeOf(key).Kind(), m.keyType.Kind(); actual != expect {
		panic(fmt.Sprintf("key not match, actual: %v, expect: %v", actual, expect))
	}

	m.getElem(key).Set(val)
}

func (m *BasicMap) Contains(key interface{}) bool {
	v := m.getElem(key)

	return v.IsAssigned()
}

func (m *BasicMap) Del(key interface{}) {
	m.getElem(key).Del()
}

func (m *BasicMap) GetKVType() (key, val reflect.Type) {
	return m.keyType, m.valType
}

func (m *BasicMap) getElem(key interface{}) StateVariable {
	bts, err := json.Marshal(key)
	if err != nil {
		panic("key cannot marshal")
	}

	keyStr := mapPrefix + m.name + string(bts)
	elem, err := GetBasicStateVariable(m.state, keyStr, m.valType)
	if err != nil {
		panic(fmt.Sprintf("getElem err: %v", err))
	}

	return elem
}
