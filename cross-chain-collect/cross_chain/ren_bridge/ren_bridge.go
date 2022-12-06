package renbridge

import (
	"app/model"
	"app/utils"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
)

type RenBridge struct {
}

var _ model.MsgCollector = &RenBridge{}

func NewRenbridgeCollector() *RenBridge {
	return &RenBridge{}
}

func (r *RenBridge) Name() string {
	return "RenBridge"
}

func (r *RenBridge) Contracts(chain string) []string {
	if _, ok := contracts[chain]; !ok {
		return nil
	}
	addrs := make([]string, 0)
	for addr := range contracts[chain] {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (r *RenBridge) Selectors(chain string) []string {
	return []string{Burn, Mint}
}

func (r *RenBridge) Extract(chain string, msgs []*model.Call) model.Results {
	if _, ok := contracts[chain]; !ok {
		return nil
	}
	var ok bool
	ret := make(model.Results, 0)
	for _, msg := range msgs {
		if _, ok := contracts[chain][msg.To]; !ok {
			continue
		}
		if len(msg.Input) <= 10 {
			continue
		}
		sig, rawParam := msg.Input[:10], msg.Input[10:]
		params, err := Decode(sig, rawParam)
		if err != nil {
			log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash, "err", err)
			continue
		}
		res := &model.Result{
			Chain:    chain,
			Number:   msg.Number,
			Index:    msg.Index,
			Hash:     msg.Hash,
			ActionId: msg.Id,
			Project:  r.Name(),
			Contract: msg.To,
			// non common
			Token: contracts[chain][msg.To].Token,
		}
		switch sig {
		case Burn:
			if len(params) < 2 {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Direction = model.OutDirection
			res.FromChainId = utils.GetChainId(chain)
			res.FromAddress = msg.From
			res.ToChainId = new(big.Int).Set(contracts[chain][msg.To].ChainId)
			to, ok := params[0].([]byte)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.ToAddress = hexutil.Encode(to)
			res.Amount, ok = params[1].(*big.Int)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
		case Mint:
			if len(params) < 4 {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
			res.Direction = model.InDirection
			res.FromChainId = new(big.Int).Set(contracts[chain][msg.To].ChainId)
			res.ToChainId = utils.GetChainId(chain)
			res.ToAddress = msg.From
			res.Amount, ok = params[1].(*big.Int)
			if !ok {
				log.Debug("decode ren bridge failed", "chain", chain, "hash", msg.Hash)
				continue
			}
		default:
			continue
		}
		ret = append(ret, res)
	}
	return ret
}
