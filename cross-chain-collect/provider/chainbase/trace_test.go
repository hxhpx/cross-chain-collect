package chainbase

import (
	"app/utils"
	"fmt"
	"os"
	"sort"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

func TestExec1(t *testing.T) {
	ret, err := Exec[*Trace]("select * from ethereum.trace_calls where block_number >= 16000000 and block_number < 16000003", "2FtLTBTxc9h7CX3YwBeEkrMlnhc")
	fmt.Println(err)
	utils.PrintPretty(ret)
}

func TestExec2(t *testing.T) {
	ret, err := Exec[*Trace]("select * from ethereum.transactions where block_number >= 16000000 and block_number < 16000015", "2FtLTBTxc9h7CX3YwBeEkrMlnhc")
	fmt.Println(err)
	utils.PrintPretty(ret)
}

func TestGetCalls(t *testing.T) {
	log.Root().SetHandler(log.LvlFilterHandler(
		log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	))
	p := NewProvider("ethereum", "2FtLTBTxc9h7CX3YwBeEkrMlnhc", true)
	ret, err := p.GetCalls([]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"}, []string{"0x2e1a7d4d", "0xa9059cbb"}, 16068492, 16068492)
	fmt.Println(err)
	utils.PrintPretty(ret)
}

type S struct {
	Id int
}
type SS []*S

func (s SS) Len() int           { return len(s) }
func (s SS) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SS) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func TestSli(t *testing.T) {
	a := []*S{{2}, {1}, {3}}
	sort.Stable(SS(a))
	utils.PrintPretty(a)
}
