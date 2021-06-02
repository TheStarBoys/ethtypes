package ethtypes

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type TypeFactory struct {
	state *ContractState
}

func NewTypeFactory(db vm.StateDB, contractAddr common.Address) (*TypeFactory, error) {
	state := NewContractState(db, contractAddr)
	return &TypeFactory{
		state: state,
	}, nil
}

func (t *TypeFactory) NewVariable(name string, initialVal interface{}) StateVariable {
	v, err := NewBasicStateVariable(t.state, name, initialVal)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) GetVariable(name string, typ reflect.Type) StateVariable {
	v, err := GetBasicStateVariable(t.state, name, typ)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) NewString(name, initialVal string) StateVariable {
	v, err := NewBasicStateVariable(t.state, name, initialVal)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) NewInt(name string, initialVal int) StateVariable {
	v, err := NewBasicStateVariable(t.state, name, initialVal)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) NewUint(name string, initialVal uint) StateVariable {
	v, err := NewBasicStateVariable(t.state, name, initialVal)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) NewFloat64(name string, initialVal float64) StateVariable {
	v, err := NewBasicStateVariable(t.state, name, initialVal)
	if err != nil {
		panic(err)
	}

	return v
}

func (t *TypeFactory) NewArray(name string, length int, typ reflect.Type) Array {
	arr, err := NewBasicArray(t.state, name, length, typ)
	if err != nil {
		panic(err)
	}

	return arr
}

func (t *TypeFactory) GetArray(name string, length int, typ reflect.Type) Array {
	arr, err := GetBasicArray(t.state, name, typ)
	if err != nil {
		panic(err)
	}

	return arr
}

func (t *TypeFactory) NewStringArray(name string, length int, initialData []string) Array {
	if len(initialData) > length {
		panic("initialData's length more than length")
	}

	arr, err := NewBasicArray(t.state, name, length, reflect.TypeOf(""))
	if err != nil {
		panic(err)
	}

	for i, v := range initialData {
		arr.Set(i, v)
	}
	return arr
}

func (t *TypeFactory) NewSlice(name string, length, cap int, typ reflect.Type) Slice {
	slice, err := NewBasicSlice(t.state, name, length, cap, typ)
	if err != nil {
		panic(err)
	}

	return slice
}

func (t *TypeFactory) GetSlice(name string, length, cap int, typ reflect.Type) Slice {
	slice, err := GetBasicSlice(t.state, name, typ)
	if err != nil {
		panic(err)
	}

	return slice
}

func (t *TypeFactory) NewStringSlice(name string, length, cap int, initialData []string) Slice {
	if len(initialData) > length {
		panic("initialData's length more than length")
	}

	slice, err := NewBasicSlice(t.state, name, length, cap, reflect.TypeOf(""))
	if err != nil {
		panic(err)
	}

	for i, v := range initialData {
		slice.Set(i, v)
	}
	return slice
}

func (t *TypeFactory) NewMap(name string, keyType, valType reflect.Type) Map {
	m, err := NewBasicMap(t.state, name, keyType, valType)
	if err != nil {
		panic(err)
	}

	return m
}

func (t *TypeFactory) GetMap(name string, keyType, valType reflect.Type) Map {
	m, err := GetBasicMap(t.state, name, keyType, valType)
	if err != nil {
		panic(err)
	}

	return m
}

func (t *TypeFactory) NewIterableMap(name string, keyType, valType reflect.Type) IterableMap {
	m, err := NewBasicIterableMap(t.state, name, keyType, valType)
	if err != nil {
		panic(err)
	}

	return m
}

func (t *TypeFactory) GetIterableMap(name string, keyType, valType reflect.Type) IterableMap {
	m, err := GetBasicIterableMap(t.state, name, keyType, valType)
	if err != nil {
		panic(err)
	}

	return m
}
