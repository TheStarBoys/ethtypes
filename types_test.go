package ethtypes

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/stretchr/testify/assert"
)

// structure for test
type Location struct {
	X, Y, Z int
}

type Person struct {
	Name string
	Age  uint8
	Loc  Location
}

func TestStateVariable(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	addr := common.HexToAddress("123")
	cs := NewContractState(state, addr)
	typeFactory := TypeFactory{state: cs}

	tests := []struct {
		variableName string
		wantVal      interface{}
		isGoodCase   bool
	}{
		{
			"good_int1",
			10,
			true,
		},
		{
			"good_int2",
			0,
			true,
		},
		{
			"bad_int",
			123,
			false,
		},
		{
			"good_uint8",
			uint8(12),
			true,
		},
		{
			"bad_uint8",
			uint8(117),
			false,
		},
		{
			"good_string1",
			"hello",
			true,
		},
		{
			"good_string2",
			"",
			true,
		},
		{
			"bad_string",
			"123",
			false,
		},
		{
			"good_person1",
			&Person{
				Name: "Bob",
				Age:  12,
				Loc:  Location{1, 2, 3},
			},
			true,
		},
		{
			"good_person2",
			Person{
				Name: "Ivan",
				Age:  34,
				Loc:  Location{6, 123, 311},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.variableName, func(t *testing.T) {
			t.Run("IsAssigned", func(t *testing.T) {
				v, err := GetBasicStateVariable(cs, tt.variableName, reflect.TypeOf(tt.wantVal))
				assert.Equal(t, err, nil)
				assert.Equal(t, v.IsAssigned(), false)
			})

			v := typeFactory.NewVariable(tt.variableName, tt.wantVal)

			t.Run("Get&Set", func(t *testing.T) {
				if tt.isGoodCase {
					actualVal := reflect.New(v.Type())
					v.Get(actualVal.Interface())
					assert.Equal(t, v.IsAssigned(), true)
					if typ := reflect.TypeOf(tt.wantVal); typ.Kind() == reflect.Ptr {
						assert.Equal(t, actualVal.Interface(), tt.wantVal)
					} else {
						assert.Equal(t, actualVal.Elem().Interface(), tt.wantVal)
					}
				} else {
					defer func() {
						if err := recover(); err == nil {
							t.Error("want panic")
						}
					}()
					// panic: val must be pointer
					actualVal := reflect.New(v.Type()).Elem()
					v.Get(actualVal.Interface())
				}
			})

			t.Run("Del", func(t *testing.T) {
				if tt.isGoodCase {
					actualVal := reflect.New(v.Type())
					v.Del()
					v.Get(actualVal.Interface())
					zeroVal := reflect.New(v.Type()).Elem().Interface()
					assert.Equal(t, v.IsAssigned(), false)
					assert.Equal(t, actualVal.Elem().Interface(), zeroVal)
				}
			})

			t.Run("CopyFrom", func(t *testing.T) {
				zeroVal := reflect.New(v.Type()).Elem().Interface()
				copiedV := typeFactory.NewVariable("copied"+tt.variableName, zeroVal)
				copiedV.CopyFrom(v)

				actualVal := reflect.New(v.Type())
				v.Get(actualVal.Interface())

				actualCopiedVal := reflect.New(v.Type())
				copiedV.Get(actualCopiedVal.Interface())
				assert.Equal(t, actualCopiedVal.Elem().Interface(), actualVal.Elem().Interface())
			})
		})

	}
}

func TestArray(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	addr := common.HexToAddress("123")
	cs := NewContractState(state, addr)
	typeFactory := TypeFactory{state: cs}
	tests := []struct {
		name   string
		length int
		typ    reflect.Type
	}{
		{
			"array1",
			0,
			StringType,
		},
		{
			"array2",
			3,
			IntType,
		},
		{},
	}

	for _, tt := range tests {
		// test array
		t.Run(tt.name, func(t *testing.T) {
			array := typeFactory.NewArray(tt.name, tt.length, tt.typ)
			assert.Equal(t, array.Len(), tt.length)
			assert.Equal(t, array.ElemType(), tt.typ)
			t.Run(tt.name+"Get&Set", func(t *testing.T) {
				if array.Len() == 0 {
					return
				}
				rand.Seed(time.Now().UnixNano())
				randomIndex := rand.Intn(array.Len())
				var (
					actualVal = reflect.New(tt.typ).Interface()
					expectVal interface{}
				)
				switch tt.typ.Kind() {
				case reflect.String:
					expectVal = fmt.Sprintf("%d", randomIndex)
				case reflect.Int:
					expectVal = randomIndex
				}
				array.Set(randomIndex, expectVal)
				array.Get(randomIndex, actualVal)
				assert.Equal(t, reflect.ValueOf(actualVal).Elem().Interface(), expectVal)
			})
			t.Run(tt.name+"-del", func(t *testing.T) {

			})
		})
	}

	array := typeFactory.NewStringArray("array", 3, []string{"1", "2", "3"})
	array.Del(2)
	fmt.Println(ArrayToStr(array))
}

func TestSlice(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	addr := common.HexToAddress("123")
	cs := NewContractState(state, addr)
	typs := TypeFactory{state: cs}

	slice1 := typs.NewStringSlice("ExampeStringSlice1", 4, 5, []string{"Hello", ",", "world", "!"})
	fmt.Println(ArrayToStr(slice1))
	// Output: ["Hello", ",", "world", "!"]
	slice2 := typs.NewStringSlice("ExampeStringSlice2", 2, 3, []string{"World", "!?"})
	slice1.CopyFrom(slice2, 2, 0, slice2.Len())
	fmt.Println(ArrayToStr(slice1))

	// Output: ["Hello", ",", "World", "!?"]
	slice1.Append("I'm", "the", "star")
	fmt.Println(ArrayToStr(slice1))

	slice1.Del(3)
	fmt.Println(ArrayToStr(slice1))

	slice3 := typs.NewStringSlice("ExampeStringSlice3", 0, 0, []string{})

	append := func(s Slice, val ...interface{}) {
		s.Append(val...)
		fmt.Println(ArrayToStr(slice3))
	}
	append(slice3, "test1", "test2")
	fmt.Println(ArrayToStr(slice3))

	// Output: slice3: ["test1", "test2"]
}

func TestMap(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	addr := common.HexToAddress("123")
	cs := NewContractState(state, addr)
	typeFactory := TypeFactory{state: cs}

	personMap1 := typeFactory.NewMap("personMap1", StringType, reflect.TypeOf(Person{}))
	testMap(t, personMap1)
}

func testMap(t *testing.T, m Map) {
	var (
		key       string = "person1"
		expectVal Person = Person{
			Name: "Bob",
			Age:  18,
			Loc: Location{
				X: 123,
				Y: 321,
				Z: 111,
			},
		}
	)
	var actualVal Person

	// test set get
	ok := m.Get(key, &actualVal)
	assert.Equal(t, ok, m.Contains(key))
	assert.Equal(t, ok, false)
	m.Set(key, expectVal)
	ok = m.Get(key, &actualVal)
	assert.Equal(t, ok, m.Contains(key))
	assert.Equal(t, ok, true)
	assert.Equal(t, actualVal, expectVal)

	// test del
	m.Del(key)
	assert.Equal(t, m.Contains(key), false)
	// empty actualVal
	actualVal = Person{}
	ok = m.Get(key, &actualVal)
	assert.Equal(t, ok, false)
	assert.Equal(t, actualVal, Person{})
}

func TestIterableMap(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)
	addr := common.HexToAddress("123")
	cs := NewContractState(state, addr)

	typeFactory := TypeFactory{state: cs}
	personMap1 := typeFactory.NewIterableMap("personMap1", StringType, reflect.TypeOf(Person{}))
	testMap(t, personMap1)

	personMap1.Set("1", Person{Name: "Bob", Age: 13})
	personMap1.Set("2", Person{Name: "Alice", Age: 16})
	fmt.Println("map: ", IterableMapToStr(personMap1))
	personMap1.Del("1")
	fmt.Println("map: ", IterableMapToStr(personMap1))
}
