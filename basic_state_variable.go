package ethtypes

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

type BasicStateVariable struct {
	state *ContractState
	name  string
	loc   common.Hash
	typ   reflect.Type
}

const (
	stateVariablePrefix = "variable_"
)

var _ StateVariable = (*BasicStateVariable)(nil)

func NewBasicStateVariable(state *ContractState, variableName string, initialVal interface{}) (*BasicStateVariable, error) {
	typ := reflect.TypeOf(initialVal)

	v, err := GetBasicStateVariable(state, variableName, typ)
	if err != nil {
		return nil, err
	}

	v.Set(initialVal)

	return v, nil
}

func GetBasicStateVariable(state *ContractState, variableName string, typ reflect.Type) (*BasicStateVariable, error) {
	for count := 3; count > 0 && typ.Kind() == reflect.Ptr; count-- {
		// panic("cannot set pointer variable")
		typ = typ.Elem()
	}

	sv := &BasicStateVariable{
		state: state,
		name:  variableName,
		typ:   typ,
	}

	sv.loc = sha256.Sum256([]byte(stateVariablePrefix + variableName))

	return sv, nil
}

func (sv *BasicStateVariable) IsAssigned() bool {
	return sv.state.Exists(sv.loc)
}

func (sv *BasicStateVariable) Addr() common.Hash {
	return sv.loc
}

func (sv *BasicStateVariable) Type() reflect.Type {
	return sv.typ
}

func (sv *BasicStateVariable) Name() string {
	return sv.name
}

func (sv *BasicStateVariable) Del() {
	sv.state.Delete(sv.loc)
}

func (sv *BasicStateVariable) Set(val interface{}) {
	v := reflect.ValueOf(val)
	for count := 3; count > 0 && v.Kind() == reflect.Ptr; count-- {
		// panic("cannot set pointer variable")
		v = v.Elem()
	}
	if v.Kind() != sv.typ.Kind() {
		panic(fmt.Sprintf("expect kind: %v, actual kind: %v", sv.typ.Kind(), v.Kind()))
	}

	if v.Kind() == reflect.Struct && v.Type().Name() != sv.typ.Name() {
		panic(fmt.Sprintf("expect type: %v, actual type: %v", sv.typ.Name(), v.Type().Name()))
	}

	byts, err := json.Marshal(val)
	if err != nil {
		panic("cannot marshal val")
	}

	sv.state.Write(sv.loc, byts)
}

func (sv *BasicStateVariable) Get(val interface{}) bool {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Ptr {
		panic("val must be pointer")
	}

	elem := v.Elem()
	if elem.Kind() != sv.typ.Kind() {
		panic(fmt.Sprintf("expect kind: %v, actual kind: %v", sv.typ.Kind(), elem.Kind()))
	}
	bts := sv.state.Read(sv.loc)

	// returns zero value if err != nil
	err := json.Unmarshal(bts, val)
	if err != nil {
		zeroVal := reflect.Zero(elem.Type())
		elem.Set(zeroVal)
	}

	return sv.IsAssigned()
}

func (sv *BasicStateVariable) CopyFrom(src StateVariable) {
	if src.Type().Kind() != sv.Type().Kind() {
		panic("kind not match")
	}
	val := reflect.New(src.Type()).Interface()
	src.Get(val)
	sv.Set(val)
}
