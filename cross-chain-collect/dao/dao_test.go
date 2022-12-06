package dao

import (
	"app/cross_chain/hop"
	"app/model"
	"encoding/json"
	"fmt"
	"testing"
)

func TestDao(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	var f *uint64

	err := d.db.Get(&f, "select match_id from common_cross_chain where hash = '0x400b912ee6f55c80facf3e0f14347a1ad994fc241cd888dd00e31ec8db327915'")
	fmt.Println(err)
	fmt.Println(*f)
}

func TestGet(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")

	stmt := "SELECT * FROM common_cross_chain WHERE id = $1"

	res := model.Data{}
	_ = d.db.Get(&res, stmt, 1156329)
	fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)

	_ = d.db.Get(&res, stmt, 1131020)
	fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)
	_ = d.db.Get(&res, stmt, 1131019)
	fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)
	_ = d.db.Get(&res, stmt, 720835)
	fmt.Println(res.Id, res.Chain, res.Hash, res.MatchId, res.MatchTag)

}

func TestUpdate(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	stmt := "UPDATE common_cross_chain SET match_id = $2 WHERE id = $1"
	_, _ = d.db.Exec(stmt, 85608, 1001689)
	_, _ = d.db.Exec(stmt, 85623, 1001701)
	_, _ = d.db.Exec(stmt, 1001701, 85623)

}

func TestDelCol(t *testing.T) {
	d := NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")

	res := []model.Data{}
	//r := &model.Data{}

	stmt := "SELECT * FROM common_cross_chain WHERE project = $1"
	err := d.db.Select(&res, stmt, "Hop")

	for _, e := range res {
		if len(e.MatchTag) != 66 {
			var buf hop.Detail
			err := json.Unmarshal(e.Detail, &buf)
			fmt.Println(err)
			str := buf.DDL.String() + e.ToAddress + buf.MinDy.String()
			st := "UPDATE common_cross_chain SET match_tag = $1 WHERE id = $2"
			_, _ = d.db.Exec(st, str, e.Id)

		}
	}
	fmt.Println(err)
}

//1135080
