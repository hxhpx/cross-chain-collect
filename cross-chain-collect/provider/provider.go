package provider

import (
	"app/config"
	"app/model"
	"app/provider/chainbase"
	"app/provider/etherscan"
	"app/provider/geth"
	"math/big"
)

type Provider struct {
	geth      *geth.GethProvider
	scan      *etherscan.EtherscanProvider
	chainbase *chainbase.Provider
}

func (p *Provider) LatestNumber() (uint64, error) {
	return p.geth.LatestNumber()
}

func (p *Provider) Call(from, to, input string, value *big.Int, number *big.Int) ([]byte, error) {
	return p.geth.Call(from, to, input, value, number)
}

func (p *Provider) GetLogs(addresses []string, topics0 []string, from, to uint64) (model.Events, error) {
	return p.geth.GetLogs(addresses, topics0, from, to)
}

func (p *Provider) GetContractFirstInvocation(address string) (uint64, error) {
	return p.scan.GetContractFirstInvocation(address)
}

func (p *Provider) GetCalls(addresses, selectors []string, from, to uint64) ([]*model.Call, error) {
	if p.chainbase == nil {
		return nil, nil
	}
	return p.chainbase.GetCalls(addresses, selectors, from, to)
}

type Providers struct {
	providers map[string]*Provider
}

func NewProviders(cfg *config.Config) *Providers {
	providers := make(map[string]*Provider)
	for chainName, providerCfg := range cfg.ChainProviders {
		gethP := geth.NewGethProvider(providerCfg.Node)
		scanP := etherscan.NewEtherScanProvider(providerCfg.ScanUrl, providerCfg.ApiKeys, cfg.Proxy)
		providers[chainName] = &Provider{
			geth: gethP,
			scan: scanP,
		}
		if providerCfg.ChainbaseTable != "" {
			providers[chainName].chainbase = chainbase.NewProvider(providerCfg.ChainbaseTable, cfg.ChainbaseApiKey, providerCfg.EnableTraceCall)
		}
	}
	return &Providers{providers: providers}
}

func (p *Providers) Get(chain string) *Provider {
	if val, ok := p.providers[chain]; ok {
		return val
	}
	return nil
}
