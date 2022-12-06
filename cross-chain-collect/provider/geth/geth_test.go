package geth

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestGetLogs(t *testing.T) {
	p := NewGethProvider("https://rpc.ankr.com/bsc")
	r, err := p.GetLogs([]string{"0xe9e7cea3dedca5984780bafc599bd69add087d56"}, []string{"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, 23477015, 23477115)
	fmt.Println(len(r), err)
}

func TestCall(t *testing.T) {
	p := NewGethProvider("https://rpc.ankr.com/eth")
	r, err := p.Call("", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0x313ce567", nil, nil)
	dec := new(big.Int).SetBytes(r[:])
	fmt.Println(r, err, dec, common.BytesToAddress(r))
}
