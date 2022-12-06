package match

import (
	"app/dao"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestFilterUnmatched(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	data1 := m.filterUnmatched(d, "Anyswap", "in")
	data2 := m.filterUnmatched(d, "Anyswap", "out")
	println("in:", len(data1))
	println("out:", len(data2))
	fmt.Println(data1[0])
	fmt.Println(data2[0])
}

// since multimatched data are too many in Across
// TestGetMultiMatched_Across gets those multi-matched ones
func TestGetMultiMatched_Across(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	file, err := os.OpenFile("./t.json", os.O_RDWR, 0666)
	if err != nil {
		println("failed to open file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	n := 0
	for {
		str, err := reader.ReadString('\n')
		n++
		if err == io.EOF {
			break
		}
		if n%4 != 3 {
			continue
		}
		oriId, _ := strconv.ParseInt(strings.Fields(strings.TrimSpace(str))[1], 10, 64)
		println(oriId)
		data, err := d.GetOne(-1, "", oriId)
		fmt.Println(data.Id, data.Chain, data.Hash, data.MatchId, data.MatchTag)
	}
}

// Followings are functions to complete basic matches
func TestMatchAnyswap(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "Anyswap")
	fmt.Println(err)
}

func TestMatchHop(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "Hop")
	fmt.Println(err)
}

func TestMatchSynapse(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "Synapse")
	fmt.Println(err)
}

func TestMatchAcross(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "Across")
	fmt.Println(err)
}

func TestMatchCBridge(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "CBridge")
	fmt.Println(err)
}

func TestMatchStargate(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "Stargate")
	fmt.Println(err)
}

func TestMatchrenBridge(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Match{}
	err := m.MatchTx(d, "RenBridge")
	fmt.Println(err)
}
