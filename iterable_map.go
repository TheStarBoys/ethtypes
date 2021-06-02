package ethtypes

import (
	"fmt"
	"reflect"
)

type BasicIterableMap struct {
	data Map
	// keys[index] = key
	keys Slice
}

const (
	iterableMapKeysInitialSize = 10
	iterableMapKeysPrefix      = "iterable_map_keys_"
)

func NewBasicIterableMap(state *ContractState, name string, keyType, valType reflect.Type) (*BasicIterableMap, error) {
	m, err := NewBasicMap(state, name, keyType, valType)
	if err != nil {
		return nil, err
	}

	keysName := iterableMapKeysPrefix + name
	slice, err := NewBasicSlice(state, keysName, 0, iterableMapKeysInitialSize, keyType)
	if err != nil {
		return nil, err
	}

	return &BasicIterableMap{
		data: m,
		keys: slice,
	}, nil
}

func GetBasicIterableMap(state *ContractState, name string, keyType, valType reflect.Type) (*BasicIterableMap, error) {
	m, err := GetBasicMap(state, name, keyType, valType)
	if err != nil {
		return nil, err
	}

	keysName := iterableMapKeysPrefix + name
	slice, err := GetBasicSlice(state, keysName, keyType)
	if err != nil {
		return nil, err
	}

	return &BasicIterableMap{
		data: m,
		keys: slice,
	}, nil
}

func (im *BasicIterableMap) Name() string {
	return im.data.Name()
}

func (im *BasicIterableMap) Get(key interface{}, val interface{}) (ok bool) {
	return im.data.Get(key, val)
}

func (im *BasicIterableMap) Set(key, val interface{}) {
	if !im.data.Contains(key) {
		im.keys.Append(key)
	}
	im.data.Set(key, val)
}

func (im *BasicIterableMap) Contains(key interface{}) bool {
	return im.data.Contains(key)
}

func (im *BasicIterableMap) Del(key interface{}) {
	if !im.Contains(key) {
		return
	}

	im.data.Del(key)
	for i := 0; i < im.keys.Len(); i++ {
		keyType := im.keys.ElemType()
		k := reflect.New(keyType).Interface()
		im.keys.Get(i, k)
		if reflect.DeepEqual(reflect.ValueOf(k).Elem().Interface(), key) {
			im.keys.Del(i)
			return
		}
	}

	panic(fmt.Sprintf("cannot find key %v in keys", key))
}

func (im *BasicIterableMap) Len() int {
	return im.keys.Len()
}

func (im *BasicIterableMap) GetKVType() (key, val reflect.Type) { return im.data.GetKVType() }

func (im *BasicIterableMap) Index(i int, key, val interface{}) {
	im.keys.Get(i, key)

	v := reflect.ValueOf(key)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	key = v.Interface()
	im.data.Get(key, val)
}

func (im *BasicIterableMap) Range(fn func(key, val interface{}) bool) {
	for i := 0; i < im.Len(); i++ {
		keyType, valType := im.data.GetKVType()
		key := reflect.New(keyType).Interface()
		val := reflect.New(valType).Interface()
		im.Index(i, key, val)

		key = reflect.ValueOf(key).Elem().Interface()
		if fn(key, reflect.ValueOf(val).Elem().Interface()) == false {
			break
		}
	}
}
