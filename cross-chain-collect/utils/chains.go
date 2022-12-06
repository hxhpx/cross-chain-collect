package utils

import "math/big"

// define unstandard chanid (btc solana ...)

var unstandardChains = map[string]*big.Int{
	"eth":       new(big.Int).SetUint64(1),
	"bsc":       new(big.Int).SetUint64(56),
	"polygon":   new(big.Int).SetUint64(137),
	"avalanche": new(big.Int).SetUint64(43114),
	"arbitrum":  new(big.Int).SetUint64(42161),
	"optimism":  new(big.Int).SetUint64(10),
	"cronos":    new(big.Int).SetUint64(25),
	"fantom":    new(big.Int).SetUint64(250),
	"moonbeam":  new(big.Int).SetUint64(1284),

	// non-standard chain id
	"btc":    new(big.Int).SetUint64(100000000),
	"solana": new(big.Int).SetUint64(100000001),
}

func GetChainId(name string) *big.Int {
	if val, ok := unstandardChains[name]; ok {
		return new(big.Int).Set(val)
	}
	return new(big.Int).SetUint64(0)
}
