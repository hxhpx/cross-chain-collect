package match

import (
	"app/dao"
	"app/model"
	"fmt"
	_ "github.com/lib/pq"
)

type Match struct {
}

type MatchedId struct {
	SrcID uint64 `db:"src_id"`
	DstID uint64 `db:"dest_id"`
}

const (
	MATCHED        = 0
	NO_MATCHED     = -1
	MULTI_MATCHED  = 2
	INFO_NOT_MATCH = 3 //to_addr或direction不匹配
)

// 输入的是一笔cross-in tx
func (a *Match) OneMatch(d *dao.Dao, data model.Data) (int, []MatchedId) {
	var res []MatchedId
	match_list, err := d.GetOneMatched(data)

	if len(match_list) == 0 || err != nil { //匹配不成功
		fmt.Println(err)
		return NO_MATCHED, res
	}

	//先将所有信息录入
	for _, e := range match_list {
		var ret = MatchedId{
			e.Id, data.Id,
		}
		res = append(res, ret)
	}

	//验证筛选出的条目
	if len(match_list) == 1 { // 只有一条数据，该数据的match_id一定为空，只需要匹配信息即可
		e := match_list[0]
		if e.ToAddress == data.ToAddress && e.Direction == model.OutDirection {
			d.UpdateMatchId(e.Id, data.Id)
			d.UpdateMatchId(data.Id, e.Id)
			return MATCHED, res

		} else {
			return INFO_NOT_MATCH, res
		}

	} else { //如果有多条数据
		return MULTI_MATCHED, res
	}
}

/*
@title	fastMatch
@description	use inner join to match
@auth Xiaohui Hu
*/
func (a *Match) BatchMatch_(d *dao.Dao, matched_pair *[]MatchedId, start, batch_size int) error {
	var err error
	for i := 0; i < batch_size; i++ {
		_, err = d.DB().Exec("UPDATE "+d.DBname()+" SET match_id = $2 WHERE id = $1", (*matched_pair)[i+start].SrcID, (*matched_pair)[i+start].DstID)
		if err != nil {
			fmt.Println("error to update ", (*matched_pair)[i+start].SrcID)
			return err
		}

		_, err = d.DB().Exec("UPDATE "+d.DBname()+" SET match_id = $2 WHERE id = $1", (*matched_pair)[i+start].DstID, (*matched_pair)[i+start].SrcID)
		if err != nil {
			fmt.Println("error to update ", (*matched_pair)[i+start].DstID)
			return err
		}

		if i%500 == 0 {
			fmt.Println("done: ", start+i)
		}
	}

	fmt.Println("update done: ", start+batch_size)
	return err
}

// Anyswap, Synapse, CBridge, WormHole
func (a *Match) GetMatchedPair(d *dao.Dao, project_name string) (*[]MatchedId, error) {
	res := &[]MatchedId{}

	stmt := "with t as (select * from " + d.DBname() + " where project = $1 and match_id is null)" +
		" select t1.id as dest_id, t2.id as src_id from t t1 inner join t t2 " +
		"on t1.match_tag = t2.match_tag and t1.to_address = t2.to_address and t1.direction='in' and t2.direction='out'"
	err := d.DB().Select(&(*res), stmt, project_name)

	fmt.Println(len(*res))
	return res, err
}

func (a *Match) BatchMatch(d *dao.Dao, matched_pair *[]MatchedId, start, batch_size int) error {
	var err error
	for i := 0; i < batch_size; i++ {
		_, err = d.DB().Exec("UPDATE common_cross_chain SET match_id = $2 WHERE id = $1", (*matched_pair)[i+start].SrcID, (*matched_pair)[i+start].DstID)
		if err != nil {
			fmt.Println("error to update ", (*matched_pair)[i+start].SrcID)
			fmt.Println(err)
			return err
		}

		_, err = d.DB().Exec("UPDATE common_cross_chain SET match_id = $2 WHERE id = $1", (*matched_pair)[i+start].DstID, (*matched_pair)[i+start].SrcID)
		if err != nil {
			fmt.Println("error to update ", (*matched_pair)[i+start].DstID)
			fmt.Println(err)
			return err
		}

		if i%500 == 0 {
			fmt.Println("done: ", start+i)
		}
	}

	fmt.Println("update done: ", start+batch_size)
	return err
}

func (a *Match) GetMatchedStargate(d *dao.Dao) (*[]MatchedId, error) {
	res := &[]MatchedId{}

	stmt := "with t as (select * from " + d.DBname() + " where project = 'Stargate' and match_id is null)" +
		" select t1.id as dest_id, t2.id as src_id from t t1 inner join t t2 " +
		"on t1.match_tag = t2.match_tag and t1.to_chain_id = t2.to_chain_id " +
		"and t1.from_chain_id = t2.from_chain_id and t1.direction='in' and t2.direction='out'"
	err := d.DB().Select(&(*res), stmt)

	fmt.Println(len(*res))
	return res, err
}
