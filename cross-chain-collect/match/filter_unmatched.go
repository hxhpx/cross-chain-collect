package match

import (
	"app/dao"
	"app/model"
	"fmt"
	"math/big"
	"strconv"
)

var ChainIds = map[int]bool{
	1:     true,
	10:    true,
	56:    true,
	137:   true,
	250:   true,
	43114: true,
	42161: true,
}

// @title		filterUnMatched
// @auth		Xiaohui Hu
// @description	To filter unmatched data whose chainId is supported by this project
func (a *Match) filterUnmatched(d *dao.Dao, projectName, direction string) []model.Data {

	if direction != "in" && direction != "out" {
		return nil
	}

	res := []model.Data{}
	ret := d.SelectProject(projectName, direction, "true")

	for _, e := range ret {
		var flag = false
		var id string

		if direction == "in" {
			id = (*big.Int)(e.WrapFromChainId).String()
		} else {
			id = (*big.Int)(e.WrapToChainId).String()
		}

		id_, _ := strconv.Atoi(id)

		if e.Id == 22572 {
			fmt.Println(id_)
		}
		//println(e.Id)
		flag = ChainIds[id_]

		if flag == true {
			res = append(res, e)
		}
	}
	return res
}
