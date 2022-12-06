package chainbase

import (
	"app/model"
	"app/utils"
	"context"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

const (
	ChainbaseUrl = "https://api.chainbase.online/v1/dw/query"
)

type Provider struct {
	table           string
	apiKey          string
	enableTraceCall bool
}

func NewProvider(table, apiKey string, enableTraceCall bool) *Provider {
	return &Provider{
		table:           table,
		apiKey:          apiKey,
		enableTraceCall: enableTraceCall,
	}
}

func (p *Provider) GetCalls(addresses, selectors []string, from, to uint64) ([]*model.Call, error) {
	if len(addresses) == 0 {
		return nil, nil
	}
	res := make([]*model.Call, 0)
	stmt := fmt.Sprintf("select * from %v.transactions where block_number >= %v and block_number <= %v and status = 1 and %v", p.table, from, to, formatOrCondition("to_address", addresses))
	if p.enableTraceCall {
		stmt = fmt.Sprintf("select * from %v.trace_calls where block_number >= %v and block_number <= %v and error = '' and %v and call_type = 'call'", p.table, from, to, formatOrCondition("to_address", addresses))
		if len(selectors) != 0 {
			stmt += fmt.Sprintf(" and %v", formatOrCondition("method_id", trimPrefix(selectors, "0x")))
		}
	}
	log.Debug(stmt)
	ret, err := Exec[*Trace](stmt, p.apiKey)
	if err != nil {
		return nil, err
	}
	sort.Stable(Traces(ret))
	id := uint64(0)
	prevHash := ""
	for _, t := range ret {
		if prevHash == "" || prevHash != t.Hash {
			// next is internal tx
			if len(t.TraceAddress) != 0 {
				id = 1
			} else {
				//next is external tx
				id = 0
			}
		} else {
			id += 1
		}
		prevHash = t.Hash
		if len(selectors) != 0 && !isTargetCall(t.Input, selectors) {
			continue
		}
		bigVal, _ := new(big.Int).SetString(t.Value, 10)
		res = append(res, &model.Call{
			Number: utils.ParseStrToUint64(t.Number),
			Index:  t.Index,
			Hash:   t.Hash,
			Id:     id,
			From:   t.From,
			To:     t.To,
			Input:  t.Input,
			Value:  bigVal,
		})
	}
	return res, nil
}

func Exec[T any](stmt, apiKey string) ([]T, error) {
	res := make([]T, 0)
	taskId := ""
	page := uint(0)
	for {
		ret, err := exec[T](stmt, taskId, apiKey, page)
		if err != nil {
			return nil, err
		}
		if ret.Message != "ok" {
			return nil, fmt.Errorf("chainbase error: %v", ret.Message)
		}
		if ret.Data.ErrMsg != "" {
			return nil, fmt.Errorf("chainbase error: %v", ret.Data.ErrMsg)
		}
		res = append(res, ret.Data.Result...)
		if ret.Data.NextPage != 0 {
			taskId = ret.Data.TaskId
			page = ret.Data.NextPage
		} else {
			break
		}
	}
	return res, nil
}

func exec[T any](stmt, taskId, apiKey string, page uint) (ret *Result[T], err error) {
	ret = &Result[T]{}
	u, err := url.Parse(ChainbaseUrl)
	if err != nil {
		return nil, err
	}
	reqBody := map[string]any{
		"query": stmt,
	}
	if taskId != "" && page != 0 {
		reqBody["task_id"] = taskId
		reqBody["page"] = page
	}
	opt := utils.HttpOption{
		Method: http.MethodPost,
		Url:    u,
		Header: map[string]string{
			"x-api-key": apiKey,
		},
		RequestBody: reqBody,
		Response:    ret,
	}
	err = opt.Send(context.Background())
	return
}

func trimPrefix(ss []string, prefix string) []string {
	ret := make([]string, 0, len(ss))
	for _, s := range ss {
		ret = append(ret, strings.TrimPrefix(s, prefix))
	}
	return ret
}

func formatOrCondition(field string, args []string) string {
	if len(args) == 0 {
		return ""
	}
	cond := "("
	for idx, arg := range args {
		cond += fmt.Sprintf("%v = '%v'", field, arg)
		if idx < len(args)-1 {
			cond += " or "
		}
	}
	cond += ")"
	return cond
}

func isTargetCall(input string, selectors []string) bool {
	if len(selectors) == 0 {
		return true
	}
	for _, s := range selectors {
		if len(s) != 10 {
			continue
		}
		if strings.HasPrefix(input, s) {
			return true
		}
	}
	return false
}
