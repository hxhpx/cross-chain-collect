package dao

import (
	"app/model"
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Dao struct {
	db *sqlx.DB
}

func NewDao(host string) *Dao {
	db, err := sqlx.Connect("postgres", host)
	if err != nil {
		panic(err)
	}
	return &Dao{
		db: db,
	}
}

func (d *Dao) Save(results model.Results) error {
	if len(results) == 0 {
		return nil
	}
	for _, r := range results {
		if len(r.Detail) == 0 {
			r.Detail = []byte(`{}`)
		}
		r.WrapFromChainId = (*model.BigInt)(r.FromChainId)
		r.WrapToChainId = (*model.BigInt)(r.ToChainId)
		r.WrapAmount = (*model.BigInt)(r.Amount)
	}
	stmt := "insert into common_cross_chain (chain, number, index, hash, action_id, project, contract, direction, from_chain_id, from_address, to_chain_id, to_address, token, amount, match_tag, detail) values (:chain, :number, :index, :hash, :action_id, :project, :contract, :direction, :from_chain_id, :from_address, :to_chain_id, :to_address, :token, :amount, :match_tag, :detail)"
	_, err := d.db.NamedExec(stmt, results)
	return err
}

func (d *Dao) UpdateMatchId(id, match_id int64) error {
	stmt := "UPDATE anyswap SET match_id = $2 WHERE id = $1"
	_, err := d.db.Exec(stmt, id, match_id)
	return err
}

func (d *Dao) LatestTime() (string, error) {
	data := model.Data_{}
	stmt := "SELECT * from anyswap ORDER BY timestamp DESC LIMIT 1"
	err := d.db.Get(&data, stmt)
	time := data.Timestamp
	return time, err
}

func (d *Dao) LastUpdate(chain, project string) (uint64, error) {
	var last uint64
	stmt := "select number from common_cross_chain where chain = $1 and project = $2 order by number desc limit 1"
	err := d.db.Get(&last, stmt, chain, project)
	if err == sql.ErrNoRows {
		err = nil
	}
	return last, err
}

// @title	SelectProject
// @description	Get data of one project with the limitation of direction and whether the match_id is empty
// @auth	Hu xiaohui
// @param	projectName	string	the project name
// @param
// @return	complete infomation of the data
func (d *Dao) SelectProject(projectName, direct, empty string) []model.Data {
	if len(projectName) == 0 || len(direct) == 0 {
		return nil
	}

	res := []model.Data{}

	var err error
	if empty == "false" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 AND match_id IS NOT NULL order by number desc"
		err = d.db.Select(&res, stmt, projectName, direct)
		//fmt.Println(*res[0])
	} else if empty == "true" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 AND match_id IS NULL order by number desc"
		err = d.db.Select(&res, stmt, projectName, direct)
	} else if empty == "all" {
		stmt := "SELECT * FROM common_cross_chain WHERE project = $1 AND direction = $2 "
		err = d.db.Select(&res, stmt, projectName, direct)
	} else {
		log.Warn("unavilable 'empty'")
		return nil
	}

	//res = append(res, r)

	if err != nil {
		fmt.Println("err:", err)
		return nil
	}
	return res
}

// @title	GetOne
// @description	Allows to get one data from pg with its (id, hash, matchId) OR just one info
// @auth	Hu xiaohui
// @param	id, hash, match_id
// @return	complete infomation of the data
func (d *Dao) GetOne(id int64, hash string, matchId int64) (*model.Data, error) {
	res := &model.Data{}
	res_ := &model.Data{}
	_res := &model.Data{}

	var err error
	if id >= 0 {
		err = d.db.Get(res, "SELECT * from common_cross_chain WHERE id = $1", id)
	}

	if len(hash) == 66 {
		err = d.db.Get(res_, "SELECT * FROM common_cross_chain WHERE hash = $1", hash)
		if (len(res.Hash) != 0) && (res.Hash != res_.Hash) {
			log.Warn("paramters are not matched with each other")
			return nil, err
		}
		res = res_
	}

	if matchId >= 0 {
		err = d.db.Get(_res, "SELECT * from common_cross_chain WHERE match_id = $1", matchId)
		if (len(res.Hash) != 0) && (res.Hash != _res.Hash) {
			log.Warn("paramters are not matched with each other")
			return nil, err
		}
		res = _res
	}

	return res, err
}

func (d *Dao) DB() *sqlx.DB { return d.db }
