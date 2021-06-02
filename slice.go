package ethtypes

import (
	"fmt"
	"math"
	"reflect"
)

type BasicSlice struct {
	arr   Array
	name  string
	state *ContractState
	len   StateVariable
}

const (
	slicePrefix       = "slice_"
	sliceLengthPrefix = "slice_length_"
)

func NewBasicSlice(state *ContractState, name string, len, cap int, typ reflect.Type) (*BasicSlice, error) {
	array, err := NewBasicArray(state, name, cap, typ)
	if err != nil {
		return nil, err
	}

	lenStr := sliceLengthPrefix + name
	lenV, err := NewBasicStateVariable(state, lenStr, len)
	if err != nil {
		return nil, err
	}

	slice := &BasicSlice{
		arr:   array,
		name:  name,
		len:   lenV,
		state: state,
	}

	return slice, nil
}

func GetBasicSlice(state *ContractState, name string, typ reflect.Type) (*BasicSlice, error) {
	array, err := GetBasicArray(state, name, typ)
	if err != nil {
		return nil, err
	}

	lenStr := sliceLengthPrefix + name
	lenV, err := GetBasicStateVariable(state, lenStr, IntType)
	if err != nil {
		return nil, err
	}

	slice := &BasicSlice{
		arr:   array,
		name:  name,
		len:   lenV,
		state: state,
	}

	return slice, nil
}

func (s *BasicSlice) Len() int {
	var length int
	s.len.Get(&length)

	return length
}

func (s *BasicSlice) Name() string {
	return s.name
}

func (s *BasicSlice) Cap() int {
	return s.arr.Len()
}

func (s *BasicSlice) ElemType() reflect.Type {
	return s.arr.ElemType()
}

// func (s *BasicSlice) getRaw(index int) []byte { return s.arr.getRaw(index) }

// func (s *BasicSlice) setRaw(index int, val []byte) { s.arr.setRaw(index, val) }

func (s *BasicSlice) Get(index int, val interface{}) {
	if s.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}

	s.arr.Get(index, val)
}
func (s *BasicSlice) Set(index int, val interface{}) {
	if s.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}

	s.arr.Set(index, val)
}

func (s *BasicSlice) Del(index int) {
	if s.isOutOfRange(index) {
		panic(ErrIndexOutOfRange)
	}
	s.arr.Del(index)
	s.len.Set(s.Len() - 1)
}

func (s *BasicSlice) CopyFrom(src Array, dstFrom, srcFrom, srcTo int) {
	count := srcTo - srcFrom
	if count > s.Len()-dstFrom {
		panic(ErrIndexOutOfRange)
	}

	s.arr.CopyFrom(src, dstFrom, srcFrom, srcTo)
}

func (s *BasicSlice) Append(vals ...interface{}) {
	length := s.Len()

	if cap := s.Cap(); cap <= len(vals)+length {
		// extend cap as max (2 * old cap, the length of elements waited appending)
		// TODO: add max function in utils package in future,
		// and instead of math.Max()
		newCap := int(math.Max(float64(2*cap), float64(len(vals))))
		newArray, err := NewBasicArray(s.state, s.name, newCap, s.ElemType())
		if err != nil {
			panic(fmt.Sprintf("extend slice cap err: %s", err))
		}
		newArray.CopyFrom(s.arr, 0, 0, s.arr.Len())
		// TODO: delete old array
		s.arr = newArray
	}

	for _, v := range vals {
		length++
		index := length - 1
		s.arr.Set(index, v)
	}

	s.len.Set(length)
}

func (s *BasicSlice) Pop(val interface{}) {
	s.Get(s.Len()-1, val)
	// s.setRaw(s.Len()-1, nil)
	s.len.Set(s.Len() - 1)
}

func (s *BasicSlice) isOutOfRange(index int) bool {
	return s.Len() <= index || index < 0
}
