package ethtypes

import (
	"fmt"
	"reflect"
)

type BasicArray struct {
	name  string
	len   StateVariable
	typ   reflect.Type
	state *ContractState
}

const (
	arrayPrefix       = "array_"
	arrayLengthPrefix = "array_length_"
)

var _ Array = (*BasicArray)(nil)

func NewBasicArray(state *ContractState, name string, len int, typ reflect.Type) (*BasicArray, error) {
	lenStr := arrayLengthPrefix + name
	lenV, err := NewBasicStateVariable(state, lenStr, len)
	if err != nil {
		return nil, err
	}

	array := &BasicArray{
		name:  name,
		len:   lenV,
		typ:   typ,
		state: state,
	}

	return array, nil
}

func GetBasicArray(state *ContractState, name string, typ reflect.Type) (*BasicArray, error) {
	lenStr := arrayLengthPrefix + name
	lenV, err := GetBasicStateVariable(state, lenStr, IntType)
	if err != nil {
		return nil, err
	}

	array := &BasicArray{
		name:  name,
		len:   lenV,
		typ:   typ,
		state: state,
	}

	return array, nil
}

func (a *BasicArray) Len() int {
	var length int
	a.len.Get(&length)

	return length
}

func (a *BasicArray) Name() string {
	return a.name
}

func (a *BasicArray) Set(index int, val interface{}) {
	if a.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}

	a.getElem(index).Set(val)
}

func (a *BasicArray) Get(index int, val interface{}) {
	if a.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}

	a.getElem(index).Get(val)
}

func (a *BasicArray) ElemType() reflect.Type {
	return a.typ
}

func (a *BasicArray) CopyFrom(src Array, dstFrom, srcFrom, srcTo int) {
	count := srcTo - srcFrom
	if count > a.Len()-dstFrom {
		panic(ErrIndexOutOfRange)
	}

	if src.ElemType().Kind() != a.ElemType().Kind() {
		panic("kind not match")
	}

	for i := 0; i < count && srcFrom < srcTo; i++ {
		val := reflect.New(src.ElemType()).Interface()
		src.Get(srcFrom, val)
		// byts := src.getRaw(srcFrom)
		srcFrom++
		a.Set(dstFrom, val)
		// a.setRaw(dstFrom, byts)
		dstFrom++
	}
}

func (a *BasicArray) Del(index int) {
	if a.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}

	switch index {
	case a.Len() - 1:
		a.getElem(index).Del()
	default:
		a.CopyFrom(a, index, index+1, a.Len())
		a.getElem(a.Len() - 1).Del()
	}
}

func (a *BasicArray) isOutOfRange(index int) bool {
	return a.Len() <= index || index < 0
}

func (a *BasicArray) getElem(index int) StateVariable {
	// prefix + name + index
	indexStr := fmt.Sprintf("%s_%d", arrayPrefix+a.name, index)
	v, err := GetBasicStateVariable(a.state, indexStr, a.typ)
	if err != nil {
		panic(fmt.Sprintf("getElem err: %v", err))
	}

	return v
}
