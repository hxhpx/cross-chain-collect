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
	db   *sqlx.DB
	name string
}

func NewDao(host string) *Dao {
	db, err := sqlx.Connect("postgres", host)
	if err != nil {
		panic(err)
	}
	return &Dao{
		db:   db,
		name: "common_cross_chain",
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

func (d *Dao) UpdateMatchId(id, match_id uint64) error {
	stmt := "UPDATE anyswap SET match_id = $2 WHERE id = $1"
	_, err := d.db.Exec(stmt, id, match_id)
	return err
}

func (d *Dao) LatestTime() (string, error) {
	data := model.Data{}
	stmt := "SELECT * from anyswap ORDER BY timestamp DESC LIMIT 1"
	err := d.db.Get(&data, stmt)
	time := data.Time
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

/*
@title	SelectProject
@description	Get data of one project with the limitation of direction and whether the match_id is empty
@auth	Hu xiaohui
@param	projectName	string	the project name
@param
@return	complete infomation of the data
*/
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

	var err error
	if id >= 0 {
		err = d.db.Get(res, "SELECT * FROM anyswap WHERE id = $1", id)
	}

	if len(hash) == 66 {
		err = d.db.Get(&(*res_), "SELECT * FROM anyswap WHERE tx_hash = $1", hash)
		err = d.db.Get(res_, "SELECT * FROM anyswap WHERE hash = $1", hash)
		if (len(res.Hash) != 0) && (res.Hash != res_.Hash) {
			log.Warn("paramters are not matched with each other")
			return nil, err
		}
		res = res_
	}

	if matchId >= 0 {
		err = d.db.Get(&(*res_), "SELECT * from anyswap WHERE match_id = $1", matchId)
		if (len(res.Hash) != 0) && (res.Hash != res_.Hash) {
			err = d.db.Get(&res, "SELECT * from anyswap WHERE match_id = $1", matchId)
			if (len(res.Hash) != 0) && (res.Hash != res_.Hash) {
				log.Warn("paramters are not matched with each other")
				return nil, err
			}
			res = res_
		}
	}
	return res, err
}

/*
@title	GetOneMatched
@description	to match an unmatched tx
@auth	xiaohui Hu
@param	[]uint64
@return matched id
*/

func (d *Dao) GetOneMatched(data model.Data) ([]model.Data, error) {
	var res []model.Data
	stmt := "select * from " + d.DBname() + " where project = $1 and to_chain_id = $1 and match_tag = $3 and id != $4"
	err := d.db.Select(&res, stmt, data.Project, data.WrapToChainId, data.MatchTag, data.Id)
	return res, err
}

func (d *Dao) DB() *sqlx.DB { return d.db }

func (d *Dao) DBname() string { return d.name }

// @title	SelectBatchProject
// @description	To select data with a specific batch size
// Used for split tables
// @auth Xiaohui Hu
/*func (d *Dao) SelectBatchProject(to_chain_id int, direction string, start_time string, batch_size int) []model.Data {
	res := []model.Data{}

	stmt := "SELECT * FROM anyswap WHERE timestamp <= $3 AND direction=$1 AND to_chain = $2 AND match_id IS NULL ORDER BY timestamp DESC LIMIT $4"
	err := d.db.Select(&res, stmt, direction, to_chain_id, start_time, batch_size)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res
}*/

// @ title SelectWithToken
/*func (d *Dao) SelectWithToken(chain, direction, token string) []model.Data_ {
	res := []model.Data_{}
	stmt := "select * from anyswap where chain = $1 and direction = $2 and token = $3"
	err := d.db.Select(&res, stmt, chain, direction, token)
	if err != nil {
		fmt.Println(err)
	}
	return res
}*/
