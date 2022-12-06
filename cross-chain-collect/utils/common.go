package utils

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/exp/constraints"
)

func PrintPretty(data interface{}) {
	res, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(res))
}

func StrSliceToLower(s []string) []string {
	ret := make([]string, 0)
	for _, r := range s {
		ret = append(ret, strings.ToLower(r))
	}
	return ret
}

func HexSum(hexes ...string) *big.Int {
	ret := new(big.Int).SetUint64(0)
	for _, hex := range hexes {
		t, _ := new(big.Int).SetString(strings.TrimPrefix(hex, "0x"), 16)
		ret.Add(ret, t)
	}
	return ret
}

func Contains[T constraints.Ordered](target T, slice []T) bool {
	for _, e := range slice {
		if target == e {
			return true
		}
	}
	return false
}
