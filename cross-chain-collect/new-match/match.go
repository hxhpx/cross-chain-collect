package newmatch

import (
	"app/dao"
	"app/match"
	"app/model"
	"bufio"
	"context"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Newmatch struct {
}

/*
func (a *Newmatch) MatchTxs(d *dao.Dao, projectName string) error {
	var flag = &map[int64]bool{}
	var chain_ids = match.ChainIds
	var err error
	batch_size := 10000
	match_nonce := uint64(0)
	latest_time, err := d.LatestTime()
	i := latest_time
	last_i := ""
	j := latest_time
	last_j := ""

	chain_id_list := []int{1, 10, 56, 137, 250, 42161, 43114}

	for _, id := range chain_id_list {
		in := d.SelectBatchProject(id, "In", i, batch_size)
		if last_i == i {
			println(i)
			println("out of pg range")
			i = latest_time
			last_i = ""
			continue
		}

		for k := 0; k < len(in); k++ {
			e := in[k]
			out := d.SelectBatchProject(id, "Out", j, 2*batch_size)
			if (last_j == j) || !(chain_ids[e.WrapFromChainId]) {
				j = latest_time
				last_j = ""

				f, err := os.OpenFile("anyswapToken.txt", os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					fmt.Println("failed to open file: ", err)
				}
				defer f.Close()
				w := bufio.NewWriter(f)
				w.WriteString(fmt.Sprint(e.Id) + "\n")
				w.Flush()

				continue
			}

			matched := false
			for _, ee := range out {
				if (e.SrcTxHash == ee.Hash) && (e.ToAddress == ee.ToAddress) {

					if (*flag)[ee.Id] == true {
						me, _ := d.GetOne_(-1, "", ee.Id, "")
						fmt.Println("Multiple matched!\ne:", e.Id, e.Chain, e.Hash)
						fmt.Println("ee:", ee.Id, ee.Chain, ee.Hash)
						fmt.Println("ee-preMatched:", me.Id)

					} else {
						(*flag)[ee.Id] = true
						err = d.UpdateMatchId(e.Id, ee.Id)
						err = d.UpdateMatchId(ee.Id, e.Id)

						matched = true
						match_nonce++
						if match_nonce%5000 == 0 {
							println("already matched: ", match_nonce)
						}

						break
					}

					if err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
			if !matched {
				last_j = j
				j = out[len(out)-1].Timestamp
				println(j)
				k--
			} else {
				j = latest_time
			}
		}
		last_i = i
		i = in[len(in)-1].Timestamp
	}
	return err
}
*/

func (a *Newmatch) MatchTxsIdx(d *dao.Dao, id int, start_time string) error {
	var flag = &map[int64]bool{}
	var chain_ids = match.ChainIds
	var err error
	batch_size := 10000
	match_nonce := uint64(0)
	i := start_time

	//chain_id_list := []int{43114, 42161, 250, 137, 56, 10, 1}
	//id := chain_id_list[id_idx]

	in := d.SelectBatchProject(id, "In", i, batch_size)
	for k := 0; k < len(in); k++ {
		e := in[k]
		//out := d.SelectBatchProject(id, "Out", j, 2*batch_size)
		if !(chain_ids[e.WrapFromChainId]) || (e.MatchId.Valid) {
			continue
		}

		ee, _ := d.GetOne_(-1, e.SrcTxHash, -1, "")

		// if unmatched
		if (len(ee.Hash)) == 0 || (ee.ToChainId != e.ToChainId) || (ee.ToAddress != e.ToAddress) {
			f, err := os.OpenFile("./new-match/"+fmt.Sprint(id)+".txt", os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("failed to open file: ", err)
			}
			defer f.Close()

			w := bufio.NewWriter(f)
			w.WriteString(fmt.Sprint(e.Id) + "\n")
			w.Flush()

			continue
		}

		//if multi-matched
		if (*flag)[ee.Id] == true {
			me, _ := d.GetOne_(-1, "", ee.Id, "")
			fmt.Println("Multiple matched!\ne:", e.Id, e.Chain, e.Hash)
			fmt.Println("ee:", ee.Id, ee.Chain, ee.Hash)
			fmt.Println("ee-preMatched:", me.Id)

		} else {
			(*flag)[ee.Id] = true
			err = d.UpdateMatchId(e.Id, ee.Id)
			err = d.UpdateMatchId(ee.Id, e.Id)

			match_nonce++

			if match_nonce%500 == 0 {
				println("already matched: ", match_nonce)
			}
		}

		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	println("already matched: ", id, in[len(in)-1].Timestamp)
	return err
}

func (a *Newmatch) FastMatch(d *dao.Dao) error {
	res, err := d.FastMatch()
	if err != nil {
		return err
	}
	fmt.Sprint((*res)[0])
	tx, err := d.DB().BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	for _, match_id := range *res {
		_, err := tx.Exec("UPDATE anyswap SET match_id = $2 WHERE id = $1", match_id.SrcID, match_id.DstID)
		_, err = tx.Exec("UPDATE anyswap SET match_id = $2 WHERE id = $1", match_id.DstID, match_id.SrcID)
		if err != nil {
			fmt.Println("error to update ", match_id)
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("error to commit")
		return err
	}
	return err
}

func (a *Newmatch) MatchTxsIdxTime(d *dao.Dao, id int, start_time string) error {
	var flag = &map[int64]bool{}
	var chain_ids = match.ChainIds
	var err error

	in_size := 30000
	out_size := 2000
	match_nonce := uint64(0)
	latest_time := start_time
	i := latest_time
	last_i := ""
	j := latest_time
	last_j := ""

	//chain_id_list := []int{43114, 42161, 250, 137, 56, 10, 1}
	//id := chain_id_list[id_idx]

	in := d.SelectBatchProject(id, "In", i, in_size)
	if last_i == i {
		println(i)
		println("out of pg range")
		i = latest_time
		last_i = ""
		return nil
	}

	for k := 0; k < len(in); k++ {
		e := in[k]
		if e.MatchId.Valid || !(chain_ids[e.WrapFromChainId]) {
			continue
		}

		if last_j == j {
			j = latest_time
			last_j = ""

			f, err := os.OpenFile("./new-match/"+fmt.Sprint(id)+".txt", os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				fmt.Println("failed to open file: ", err)
			}
			defer f.Close()
			w := bufio.NewWriter(f)
			w.WriteString(fmt.Sprint(e.Id) + "\n")
			w.Flush()

			continue
		}

		j = e.Timestamp // out_tx的时间一定在in_tx之前
		out := d.SelectBatchProject(id, "Out", j, out_size)
		matched := false

		for _, ee := range out {
			if (e.SrcTxHash == ee.Hash) && (e.ToAddress == ee.ToAddress) {

				if (*flag)[ee.Id] == true {
					me, _ := d.GetOne_(-1, "", ee.Id, "")
					fmt.Println("Multiple matched!\ne:", e.Id, e.Chain, e.Hash)
					fmt.Println("ee:", ee.Id, ee.Chain, ee.Hash)
					fmt.Println("ee-preMatched:", me.Id)

				} else {
					(*flag)[ee.Id] = true
					err = d.UpdateMatchId(e.Id, ee.Id)
					err = d.UpdateMatchId(ee.Id, e.Id)

					matched = true
					match_nonce++

					if match_nonce%1000 == 0 {
						println("already matched: ", match_nonce)
					}

					break
				}

				if err != nil {
					fmt.Println(err)
					return err
				}
			}
		}
		if !matched {
			last_j = j
			j = out[len(out)-1].Timestamp // 如果在一个batch中都没有match成功，就继续往前再找一个batch
			k--
			println("one more batch: ", id, fmt.Sprint(last_j))
		}
	}
	last_i = i
	i = in[len(in)-1].Timestamp

	return err
}

func (a *Newmatch) MatchWithTokenCount(d *dao.Dao, chain, direction string) error {
	token_count := d.GetTokenCount(chain, direction)
	var err error

	for _, t := range token_count {
		data := d.SelectWithToken(chain, direction, t.Token)
		for _, e := range data {
			if e.MatchId.Valid {
				continue
			}

			ee := &model.Data_{}

			if direction == "In" {
				ee, err = d.GetOne_(-1, e.SrcTxHash, -1, "")

			} else if direction == "Out" {
				ee, err = d.GetOne_(-1, "", -1, e.Hash)
			}

			if len(ee.Hash) != 0 {
				d.UpdateMatchId(e.Id, ee.Id)
				d.UpdateMatchId(ee.Id, e.Id)

			} else {
				f, err := os.OpenFile("anyswapToken.txt", os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					fmt.Println("failed to open file: ", err)
				}
				defer f.Close()
				w := bufio.NewWriter(f)
				w.WriteString(fmt.Sprint(e.Id) + "\n")
				w.Flush()
			}
		}
	}
	return err
}
