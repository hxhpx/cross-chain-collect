package geth

import (
	"app/model"
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

type GethProvider struct {
	client *ethclient.Client
}

func NewGethProvider(addr string) *GethProvider {
	c, err := ethclient.Dial(addr)
	if err != nil {
		panic(err)
	}
	return &GethProvider{
		client: c,
	}
}

func (p *GethProvider) Call(from, to, input string, value *big.Int, number *big.Int) ([]byte, error) {
	var toAddr *common.Address
	if to != "" {
		tmp := common.HexToAddress(to)
		toAddr = &tmp
	}
	msg := ethereum.CallMsg{
		From:  common.HexToAddress(from),
		To:    toAddr,
		Value: value,
		Data:  common.FromHex(input),
	}
	return p.client.CallContract(context.Background(), msg, number)
}

func (p *GethProvider) LatestNumber() (uint64, error) {
	return p.client.BlockNumber(context.Background())
}

func (p *GethProvider) GetLogs(addresses []string, topics0 []string, from, to uint64) (model.Events, error) {
	ret := make(model.Events, 0)
	addrs := make([]common.Address, 0)
	for _, t := range addresses {
		addrs = append(addrs, common.HexToAddress(t))
	}
	qry := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(to),
		Addresses: addrs,
		Topics:    make([][]common.Hash, 0),
	}
	topic0 := make([]common.Hash, 0)
	for _, t := range topics0 {
		if len(t) != 0 {
			topic0 = append(topic0, common.HexToHash(t))
		}
	}
	if len(topic0) != 0 {
		qry.Topics = append(qry.Topics, topic0)
	}

	rawLogs, err := p.client.FilterLogs(context.Background(), qry)
	if err != nil {
		return nil, err
	}

	for _, rawLog := range rawLogs {
		if rawLog.Removed {
			continue
		}
		topics := make([]string, 0)
		for _, t := range rawLog.Topics {
			topics = append(topics, hexutil.Encode(t[:]))
		}
		ret = append(ret, &model.Event{
			Number:  rawLog.BlockNumber,
			Index:   uint64(rawLog.TxIndex),
			Hash:    hexutil.Encode(rawLog.TxHash[:]),
			Id:      uint64(rawLog.Index),
			Address: strings.ToLower(rawLog.Address.Hex()),
			Topics:  topics,
			Data:    hexutil.Encode(rawLog.Data),
		})
	}
	return ret, nil
}
