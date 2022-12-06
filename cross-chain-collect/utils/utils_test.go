package utils

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	v := "0x"
	fmt.Println(ParseStrToUint64(v))
	v = "0x1"
	fmt.Println(ParseStrToUint64(v))
	v = "0x23"
	fmt.Println(ParseStrToUint64(v))
	v = "1564"
	fmt.Println(ParseStrToUint64(v))
	v = "1ebf"
	fmt.Println(ParseStrToUint64(v))
	v = ""
	fmt.Println(ParseStrToUint64(v))
	v = "pppp"
	fmt.Println(ParseStrToUint64(v))
}
