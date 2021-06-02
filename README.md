# ethtypes

## Description
Encapsulates commonly used data structures for EVM(Ethereum Virtual Machine) state variables.

## Prerequisites
golang

## Quick Start
```go
package main

import (
	"fmt"

	"github.com/TheStarBoys/ethtypes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
)

func main() {
	// New stateDB in memory
	db := rawdb.NewMemoryDatabase()
	state, _ := state.New(common.Hash{}, state.NewDatabase(db), nil)

	// Your precompiled contract is at this address
	contractAddr := common.HexToAddress("123")
	tf, err := ethtypes.NewTypeFactory(state, contractAddr)
	if err != nil {
		panic(err)
	}

	arrayName := "Array"
	array := tf.NewArray(arrayName, 3, ethtypes.Int64Type)

	for i := 0; i < array.Len(); i++ {
		array.Set(i, int64(i+1))
	}

	for i := 0; i < array.Len(); i++ {
		var val int64
		array.Get(i, &val)
		fmt.Printf("index %d val: %d\n", i, val)
	}
	// index 0 val: 1
	// index 1 val: 2
	// index 2 val: 3

	fmt.Println("----------------")
	array.Del(0)
	for i := 0; i < array.Len(); i++ {
		var val int64
		array.Get(i, &val)
		fmt.Printf("index %d val: %d\n", i, val)
	}
}
```
## License
The ethtypes library is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included in our repository in the COPYING.LESSER file.