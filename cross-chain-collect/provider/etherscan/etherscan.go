package etherscan

import (
	"app/model"
	"app/utils"
	"fmt"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

const (
	MaxBlockNumber  = 999999999
	DefaultPageSize = 1000

	normalTxApi      = "%s/api?module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s"
	internalTxApi    = "%s/api?module=account&action=txlistinternal&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s&apikey=%s"
	logWithTopicsApi = "%s/api?module=logs&action=getLogs&address=%s&fromBlock=%d&toBlock=%d&topic0=%s&page=%d&offset=%d&apikey=%s"
	logApi           = "%s/api?module=logs&action=getLogs&address=%s&fromBlock=%d&toBlock=%d&page=%d&offset=%d&apikey=%s"

	noTransactionsFound = "No transactions found"
	noRecordsFound      = "No records found"
)

type Option struct {
	Page       int
	PageSize   int
	Asc        bool
	StartBlock int64
	EndBlock   int64
}

type EtherscanProvider struct {
	baseUrl string
	proxy   string
	// key pool
	apiKeys []string
	keyIter uint
	l       sync.Mutex
}

func NewEtherScanProvider(baseUrl string, apiKeys []string, proxy string) *EtherscanProvider {
	return &EtherscanProvider{
		baseUrl: strings.TrimRight(baseUrl, "/"),
		proxy:   proxy,
		apiKeys: apiKeys,
	}
}

func (p *EtherscanProvider) GetLogs(addresses []string, topics0 []string, from, to uint64) (model.Events, error) {
	ret := make(model.Events, 0)
	for _, address := range addresses {
		for _, topic0 := range topics0 {
			page := 1
			for {
				rawLogs, err := p.getLogs(address, topic0, Option{
					Page:       page,
					PageSize:   DefaultPageSize,
					StartBlock: int64(from),
					EndBlock:   int64(to),
				})
				if err != nil {
					return nil, err
				}
				for _, l := range rawLogs {
					log := model.Event{
						Number:  utils.ParseStrToUint64(l.BlockNumber),
						Index:   utils.ParseStrToUint64(l.Index),
						Hash:    l.Hash,
						Id:      utils.ParseStrToUint64(l.LogIndex),
						Address: l.Address,
						Topics:  l.Topics,
						Data:    l.Data,
					}
					ret = append(ret, &log)
				}
				if len(rawLogs) < DefaultPageSize {
					break
				}
				page++
			}
		}
	}
	return ret, nil
}

func (p *EtherscanProvider) LatestNumber() (uint64, error) {
	return 0, nil
}

func (p *EtherscanProvider) nextKey() string {
	p.l.Lock()
	defer p.l.Unlock()
	key := p.apiKeys[p.keyIter]
	p.keyIter = (p.keyIter + 1) % uint(len(p.apiKeys))
	return key
}

func (p *EtherscanProvider) GetContractFirstInvocation(address string) (ret uint64, err error) {
	normal, err := p.getFirstTransaction(address)
	if err != nil {
		return
	}
	if normal != nil {
		ret = utils.ParseStrToUint64(normal.BlockNumber)
	}
	internal, err := p.getFirstInternalTransaction(address)
	if err != nil {
		return
	}
	if internal != nil {
		num := utils.ParseStrToUint64(internal.BlockNumber)
		if ret == 0 {
			ret = num
		} else if num != 0 && num < ret {
			ret = num
		}
	}
	return
}

func (p *EtherscanProvider) GetTransactions(address string, o Option) ([]*NormalTx, error) {
	url := fmt.Sprintf(normalTxApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, o.Page, o.PageSize, toSortStr(o.Asc), p.nextKey())
	return doFetchData[[]*NormalTx](url, p.proxy)
}

func (p *EtherscanProvider) GetInternalTransactions(address string, o Option) ([]*InternalTx, error) {
	url := fmt.Sprintf(internalTxApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, o.Page, o.PageSize, toSortStr(o.Asc), p.nextKey())
	return doFetchData[[]*InternalTx](url, p.proxy)
}

func (p *EtherscanProvider) getFirstTransaction(address string) (*NormalTx, error) {
	txs, err := p.GetTransactions(address, Option{
		Page:       1,
		PageSize:   1,
		Asc:        true,
		StartBlock: 0,
		EndBlock:   MaxBlockNumber,
	})
	if err != nil {
		return nil, err
	}

	if len(txs) > 0 {
		return txs[0], nil
	}

	return nil, nil
}

func (p *EtherscanProvider) getFirstInternalTransaction(address string) (*InternalTx, error) {
	txs, err := p.GetInternalTransactions(address, Option{
		Page:       1,
		PageSize:   1,
		Asc:        true,
		StartBlock: 0,
		EndBlock:   MaxBlockNumber,
	})
	if err != nil {
		return nil, err
	}
	if len(txs) > 0 {
		return txs[0], nil
	}

	return nil, nil
}

func (p *EtherscanProvider) getLogs(address, topics0 string, o Option) ([]*EtherscanEvent, error) {
	url := fmt.Sprintf(logWithTopicsApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, topics0, o.Page, o.PageSize, p.nextKey())
	if topics0 == "" {
		url = fmt.Sprintf(logApi, p.baseUrl, strings.ToLower(address), o.StartBlock, o.EndBlock, o.Page, o.PageSize, p.nextKey())
	}
	return doFetchData[[]*EtherscanEvent](url, p.proxy)
}

func doFetchData[T any](url, proxy string) (r T, err error) {
	log.Debug("invoke etherscan", "url", url)
	var resp EtherscanResponse[T]
	if err = utils.HttpGetObjectWithProxy(url, proxy, &resp); err != nil {
		log.Debug("http get failed", "err", err, "url", url)
		return
	}
	if resp.Status != "1" && (resp.Message != noTransactionsFound && resp.Message != noRecordsFound) {
		err = fmt.Errorf("%s", resp.Message)
		log.Debug("etherscan get result falied", "err", err, "url", url)
		return
	}

	return resp.Result, nil
}

func toSortStr(asc bool) string {
	if asc {
		return "asc"
	} else {
		return "desc"
	}
}
