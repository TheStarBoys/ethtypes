package ethtypes

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
)

type ContractState struct {
	db   vm.StateDB
	addr common.Address
}

// TODO: To avoid same variable's name

func NewContractState(db vm.StateDB, addr common.Address) *ContractState {
	return &ContractState{
		db:   db,
		addr: addr,
	}
}

const (
	lengthSuffix = "length"
	indexSuffix  = "index"
)

//Exists - check data has been saved or not
func (s *ContractState) Exists(hash common.Hash) bool {
	lhash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash[:]...), lengthSuffix...)))
	state := s.db.GetState(s.addr, lhash)

	return state != common.Hash{}
}

func (s *ContractState) Delete(hash common.Hash) {
	lhash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash[:]...), lengthSuffix...)))
	s.db.SetState(s.addr, lhash, common.Hash{})
}

func (s *ContractState) Write(hash common.Hash, data []byte) {
	lhash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash[:]...), lengthSuffix...)))
	length := len(data)
	for offset := 0; offset < length; offset += 32 {
		ihash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash.Bytes()...),
			[]byte(fmt.Sprintf("%d_%s", offset/32, indexSuffix))...)))
		end := offset + 32
		if end > length {
			end = length
		}
		s.db.SetState(s.addr, ihash, common.BytesToHash(data[offset:end]))
	}
	s.db.SetState(s.addr, lhash, common.BigToHash(big.NewInt(int64(length))))
}

func (s *ContractState) Read(hash common.Hash) []byte {
	lhash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash[:]...), lengthSuffix...)))
	length := int(s.db.GetState(s.addr, lhash).Big().Int64())
	data := make([]byte, length)
	for offset := 0; offset < length; offset += 32 {
		ihash := common.Hash(sha256.Sum256(append(append([]byte(nil), hash.Bytes()...),
			[]byte(fmt.Sprintf("%d_%s", offset/32, indexSuffix))...)))
		val := s.db.GetState(s.addr, ihash)
		end := offset + 32
		if end > length {
			end = length
			copy(data[offset:end], val.Bytes()[32-end%32:])
		} else {
			copy(data[offset:end], val.Bytes())
		}
	}
	return data
}
