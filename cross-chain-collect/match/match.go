package match

import (
	"app/dao"
	"database/sql"
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/lib/pq"
)

type Match struct {
}

func (a *Match) MatchTx(d *dao.Dao, projectName string) error {
	flag := &map[int64]bool{}
	var err error

	for {
		in := d.SelectProject(projectName, "in", "true")
		out := d.SelectProject(projectName, "out", "true")

		if len(in) == 0 || len(out) == 0 {
			log.Warn("got nothing from db")
			err = nil
			break
		}
		println(len(in))
		println(len(out))

		for _, e := range in {
			for _, ee := range out {
				id1, _ := e.WrapToChainId.Value()
				id2, _ := ee.WrapToChainId.Value()
				id3, _ := e.WrapFromChainId.Value()
				id4, _ := ee.WrapFromChainId.Value()

				if (projectName == "Stargate") && (id3 != id4) {
					continue
				} else if (projectName != "Stargate") && (e.ToAddress != ee.ToAddress) {
					continue
				}

				if (e.MatchTag == ee.MatchTag) && (id1 == id2) {
					if (*flag)[ee.Id] == true {
						fmt.Println("Multiple matched!\ne:", e.Id, e.Chain, e.Hash)
						fmt.Println("ee:", ee.Id, ee.Chain, ee.Hash)

						fmt.Println("ee-preMatched:", ee.MatchId)
						continue
					}
					e.MatchId = sql.NullInt64{ee.Id, true}
					ee.MatchId = sql.NullInt64{e.Id, true}

					(*flag)[ee.Id] = true

					err = d.UpdateMatchId(e.Id, e.MatchId.Int64)
					if err != nil {
						fmt.Println(err)
						return err
					}
					err = d.UpdateMatchId(ee.Id, ee.MatchId.Int64)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
		}
	}
	return err
}
