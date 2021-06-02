package ethtypes

import (
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

// StateVariable represents global variable and stores in chain db
type StateVariable interface {
	// Del delete this state variable
	Del()
	// Set set state variable as val
	Set(val interface{})
	// Get get state variable, val must be pointer.
	// The value of state variable will store into val.
	// Return true if the initial value has been assigned.
	Get(val interface{}) bool
	// CopyFrom copy the value of src into
	// this state variable
	CopyFrom(src StateVariable)
	// Type returns type of this state variable
	Type() reflect.Type
	// IsAssigned returns true if the initial value
	// has been assigned.
	IsAssigned() bool
	// Addr returns the location of this state variable
	Addr() common.Hash
	// Name returns the name of state varialbe
	Name() string
}

// Array represents a simple array.
type Array interface {
	// Get get val from index, val must
	// be pointer
	Get(index int, val interface{})
	Set(index int, val interface{})
	// Del deletes element in index, and
	// move all elements after index move forward by one position
	Del(index int)
	Len() int
	CopyFrom(src Array, dstFrom, srcFrom, srcTo int)
	// ElemType returns element type of the array
	ElemType() reflect.Type
	// Range(fn func(index int, val interface{}) bool)
	// Name returns the name of array
	Name() string
}

// Slice represents a slice like slice in golang
type Slice interface {
	Array
	// Cap returns capacity of the Slice
	Cap() int
	// Append append elements in the tail of the slice
	Append(vals ...interface{})
	// Pop remove elements in the tail of the slice
	Pop(val interface{})
}

// Map represents key-value pair mapping
type Map interface {
	// Get val from key, val must be pointer,
	// returns false if element not exsit.
	Get(key interface{}, val interface{}) (ok bool)
	Set(key, val interface{})
	Contains(key interface{}) bool
	Del(key interface{})
	// GetKVType returns key-value pair type
	GetKVType() (key, val reflect.Type)
	// Name returns the name of map
	Name() string
}

// IterableMap represents a iterable mapping
type IterableMap interface {
	Map
	// Len returns the number of elems
	Len() int
	// Index get key and val by index i,
	// key and val must be pointer
	Index(i int, key, val interface{})
	// Range iterate all key-value pair, it
	// will stop if fn returns false
	Range(fn func(key, val interface{}) bool)
}
