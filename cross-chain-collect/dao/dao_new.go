package dao

import (
	"app/model"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

// @title	SelectBatchProject
// @description	To select data with a specific batch size
// Used for split tables
// @auth Xiaohui Hu
func (d *Dao) SelectBatchProject(to_chain_id int, direction string, start_time string, batch_size int) []model.Data_ {
	res := []model.Data_{}

	stmt := "SELECT * FROM anyswap WHERE timestamp <= $3 AND direction=$1 AND to_chain = $2 AND match_id IS NULL ORDER BY timestamp DESC LIMIT $4"
	err := d.db.Select(&res, stmt, direction, to_chain_id, start_time, batch_size)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res
}

// @title	GetOne
// @description	Allows to get one data from pg with its (id, hash, matchId) OR just one info
// @auth	Hu xiaohui
// @param	id, hash, match_id
// @return	complete infomation of the data
func (d *Dao) GetOne_(id int64, hash string, matchId int64, src_tx_hash string) (*model.Data_, error) {
	res := &model.Data_{}
	res_ := &model.Data_{}

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

		if len(src_tx_hash) == 66 {
			err = d.db.Get(&(*res_), "SELECT * FROM anyswap WHERE src_tx_hash = $1", src_tx_hash)
			if (len(res.Hash) != 0) && (res.Hash != res_.Hash) {
				log.Warn("paramters are not matched with each other")
				return nil, err
			}
			res = res_
		}
	}
	return res, err

}

type TokenCount struct {
	Count int    `db:"count"`
	Token string `db:"token"`
}

// @title GetTokenCount
// @description get token count asc
func (d *Dao) GetTokenCount(chain, direction string) []TokenCount {
	token_count := []TokenCount{}
	stmt := "select count(id), token from anyswap where chain = $1 and direction = $2 group by token order by count(id) asc "
	err := d.db.Select(&token_count, stmt, chain, direction)
	if err != nil {
		fmt.Println("select with token err: ", err)
	}
	return token_count
}

func (d *Dao) SelectWithToken(chain, direction, token string) []model.Data_ {
	res := []model.Data_{}
	stmt := "select * from anyswap where chain = $1 and direction = $2 and token = $3"
	err := d.db.Select(&res, stmt, chain, direction, token)
	if err != nil {
		fmt.Println(err)
	}
	return res
}

type MatchedId struct {
	SrcID uint64 `db:"src_id"`
	DstID uint64 `db:"dest_id"`
}

/*
@title	fastMatch
@description	use inner join to match
@auth Xiaohui Hu
*/
func (d *Dao) FastMatch() (*[]MatchedId, error) {
	res := &[]MatchedId{}
	stmt := "select src_table.id as src_id, dest_table.id as dest_id from anyswap src_table inner join anyswap dest_table on src_table.src_tx_hash = dest_table.tx_hash and src_table.direction = 'In' and dest_table.direction = 'Out'"
	err := d.db.Select(&(*res), stmt)
	fmt.Println(len(*res))

	return res, err
}
