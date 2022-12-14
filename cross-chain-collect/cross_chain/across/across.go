package across

import (
	"app/model"
	"encoding/json"
	"math/big"
)

var _ model.EventCollector = &Across{}

type Across struct {
}

func NewAcrossCollector() *Across {
	return &Across{}
}

func (a *Across) Name() string {
	return "Across"
}

func (a *Across) Contracts(chain string) []string {
	if _, ok := AcrossContracts[chain]; !ok {
		return nil
	}
	return AcrossContracts[chain]
}

func (a *Across) Topics0(chain string) []string {
	return []string{FundsDeposited, FilledRelay}
}

func (a *Across) Extract(chain string, events model.Events) model.Results {
	ret := make(model.Results, 0)
	len_out := 386
	len_in := 834

	for _, e := range events {
		res := &model.Result{
			Chain:    chain,
			Number:   e.Number,
			Index:    e.Index,
			Hash:     e.Hash,
			ActionId: e.Id,
			Project:  a.Name(),
			Contract: e.Address,
		}

		switch e.Topics[0] {
		case FundsDeposited:
			if len(e.Topics) < 4 || len(e.Data) < len_out {
				continue
			}
			res.Direction = model.OutDirection
			res.FromChainId, _ = new(big.Int).SetString(e.Data[2+64:2+128], 16)
			res.ToChainId, _ = new(big.Int).SetString(e.Data[2+128:2+192], 16)
			res.ToAddress = "0x" + e.Data[len_out-64+24:]
			res.FromAddress = "0x" + e.Topics[3][26:]
			res.Token = "0x" + e.Topics[2][26:]
			res.Amount, _ = new(big.Int).SetString(e.Data[2:2+64], 16)

			depositId := e.Topics[1]
			d := &Detail{
				DepositId: depositId,
			}
			detail, err := json.Marshal(d)
			if err == nil {
				res.Detail = detail
			}
			res.MatchTag = d.DepositId

		case FilledRelay:
			if len(e.Topics) < 3 || len(e.Data) < len_in {
				continue
			}
			res.Direction = model.InDirection
			relayer := "0x" + e.Topics[1][26:]
			res.FromChainId, _ = new(big.Int).SetString(e.Data[2+64*4:2+64*5], 16)
			res.FromAddress = "0x" + e.Topics[2][26:]

			res.ToChainId, _ = new(big.Int).SetString(e.Data[2+64*5:2+64*6], 16)
			depositId := "0x" + e.Data[len_in-64*4:len_in-64*3]
			res.Token = "0x" + e.Data[len_in-64*3+24:len_in-128]
			res.ToAddress = "0x" + e.Data[len_in-64*2+24:len_in-64]
			res.Amount, _ = new(big.Int).SetString(e.Data[2:2+64], 16)
			d := &Detail{
				DepositId: depositId,
				Relayer:   relayer,
			}
			detail, err := json.Marshal(d)
			if err == nil {
				res.Detail = detail
			}
			res.MatchTag = d.DepositId
		}
		ret = append(ret, res)
	}
	return ret
}
