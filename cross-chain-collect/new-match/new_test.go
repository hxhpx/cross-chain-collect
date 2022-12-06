package newmatch

import (
	"app/dao"
	"app/model"
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestNonMatched(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	f, err := os.OpenFile("anyswapToken.txt", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("failed to open file: ", err)
	}
	defer f.Close()

	res := []*model.Data_{}
	r := bufio.NewReader(f)
	for {
		buf, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		buf = strings.Split(buf, "\n")[0]
		id_, _ := strconv.Atoi(buf)
		data, _ := d.GetOne_(int64(id_), "", -1, "")
		res = append(res, data)
	}

	println(len(res))

	f, err = os.OpenFile("unmatched.json", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("failed to open file: ", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range res {
		w.WriteString(fmt.Sprint(e))
	}
	w.Flush()
}

func TestTokenMatch(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Newmatch{}
	m.MatchWithTokenCount(d, "ethereum", "In")
}

func TestFastMatch(t *testing.T) {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &Newmatch{}
	m.FastMatch(d)
}
