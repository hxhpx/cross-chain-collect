package etherscan

import (
	"fmt"
	"testing"
)

func TestGetLogs(t *testing.T) {
	p := NewEtherScanProvider("https://api.etherscan.io/", []string{"Y5CIXMXJ23Y6H8JSRAUQ5T8SMT2VV9W17Z", "4RYCK1WU1W2NBCGDNVEV36GHSZTF6CGW2M"}, "http://192.168.3.59:10809")
	ret, err := p.GetLogs([]string{"0xbd3531da5cf5857e7cfaa92426877b022e612cf8"}, nil, 12878196, 12878196)
	fmt.Println(len(ret), err)
}

func TestFirstCalled(t *testing.T) {
	p := NewEtherScanProvider("https://api.etherscan.io/", []string{"Y5CIXMXJ23Y6H8JSRAUQ5T8SMT2VV9W17Z", "4RYCK1WU1W2NBCGDNVEV36GHSZTF6CGW2M"}, "http://192.168.3.59:10809")
	fmt.Println(p.GetContractFirstInvocation("0x7782046601e7b9b05ca55a3899780ce6ee6b8b2b"))
}

func TestOptiscanFirstCalled(t *testing.T) {
	p := NewEtherScanProvider("https://api-optimistic.etherscan.io/", []string{"TX5FYFU9QWEMCQ9UGP865H74VTBWEBVW8X", "TV4RKAHHUXRVKYDJJ1ZPDXR75QJGG75WRB"}, "http://192.168.3.59:10809")
	fmt.Println(p.GetContractFirstInvocation("0x94b008aa00579c1307b0ef2c499ad98a8ce58e58"))
}
