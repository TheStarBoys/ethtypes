package ethtypes

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// Token
	Wei        *big.Int = big.NewInt(1)
	Kwei       *big.Int = new(big.Int).Mul(Wei, big.NewInt(1000))
	Mwei       *big.Int = new(big.Int).Mul(Kwei, big.NewInt(1000))
	Gwei       *big.Int = new(big.Int).Mul(Mwei, big.NewInt(1000))
	Microether *big.Int = new(big.Int).Mul(Gwei, big.NewInt(1000))
	Milliether *big.Int = new(big.Int).Mul(Microether, big.NewInt(1000))
	Ether      *big.Int = new(big.Int).Mul(Milliether, big.NewInt(1000))
	Kether     *big.Int = new(big.Int).Mul(Ether, big.NewInt(1000))
	Mether     *big.Int = new(big.Int).Mul(Kether, big.NewInt(1000))
	Gether     *big.Int = new(big.Int).Mul(Mether, big.NewInt(1000))
)

// FromWei convert wei as Kwei, Mwei, ..., Ether..
func FromWei(wei, unit *big.Int) (*big.Int, error) {
	switch {
	case unit.Cmp(Ether) == 0:
		return new(big.Int).Div(wei, unit), nil
	default:
		return nil, fmt.Errorf("unknown unit")
	}
}

var (
	// Basic Type
	StringType  = reflect.TypeOf("")
	IntType     = reflect.TypeOf(int(0))
	Int8Type    = reflect.TypeOf(int8(0))
	Int16Type   = reflect.TypeOf(int16(0))
	Int32Type   = reflect.TypeOf(int32(0))
	Int64Type   = reflect.TypeOf(int64(0))
	UintType    = reflect.TypeOf(uint(0))
	Uint8Type   = reflect.TypeOf(uint8(0))
	Uint16Type  = reflect.TypeOf(uint16(0))
	Uint32Type  = reflect.TypeOf(uint32(0))
	Uint64Type  = reflect.TypeOf(uint64(0))
	Float32Type = reflect.TypeOf(float32(0))
	Float64Type = reflect.TypeOf(float64(0))
	BoolType    = reflect.TypeOf(false)
	BytesType   = reflect.TypeOf([]byte{})
)

var (
	// Complex Type
	AddressType = reflect.TypeOf(common.Address{})
	BigIntType  = reflect.TypeOf(big.Int{})
)
