package main

import (
	"app/aggregator"
	"app/config"
	"app/dao"
	newmatch "app/new-match"
	"app/svc"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/log"
)

var logLvl = flag.String("log_level", "info", "set log level")

func main_() {
	flag.Parse()
	lvl, err := log.LvlFromString(*logLvl)
	if err != nil {
		panic(err)
	}
	fmt.Println("log level:", lvl.String())
	// log.Root().SetHandler(log.LvlFilterHandler(
	// 	lvl, log.StreamHandler(os.Stderr, log.TerminalFormat(false)),
	// ))
	log.Root().SetHandler(log.MultiHandler(
		log.LvlFilterHandler(log.LvlError, log.Must.FileHandler("./error.log", log.TerminalFormat(false))),
		log.LvlFilterHandler(lvl, log.StreamHandler(os.Stderr, log.TerminalFormat(true))),
	))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		<-sig
		cancel()
	}()
	var cfg config.Config

	config.LoadCfg(&cfg, "./config.yaml")
	srvCtx := svc.NewServiceContext(ctx, &cfg)
	for name := range srvCtx.Config.ChainProviders {
		agg := aggregator.NewAggregator(srvCtx, name)
		go agg.Start()
	}
	<-ctx.Done()
	srvCtx.Wg.Wait()
	fmt.Println("exit")
}

// =========================================================================================

func main() {
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")
	m := &newmatch.Newmatch{}

	go m.MatchTxsIdxTime(d, 10, "2022-06-24")
	go m.MatchTxsIdxTime(d, 56, "2022-07-29")
	go m.MatchTxsIdxTime(d, 137, "2022-07-22")
	go m.MatchTxsIdxTime(d, 250, "2022-07-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-07-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-07-19")
	go m.MatchTxsIdxTime(d, 1, "2022-07-01")

	go m.MatchTxsIdxTime(d, 10, "2022-05-24")
	go m.MatchTxsIdxTime(d, 56, "2022-06-29")
	go m.MatchTxsIdxTime(d, 137, "2022-06-22")
	go m.MatchTxsIdxTime(d, 250, "2022-06-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-06-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-06-19")
	go m.MatchTxsIdxTime(d, 1, "2022-06-02")

	go m.MatchTxsIdxTime(d, 10, "2022-04-24")
	go m.MatchTxsIdxTime(d, 56, "2022-05-29")
	go m.MatchTxsIdxTime(d, 137, "2022-05-22")
	go m.MatchTxsIdxTime(d, 250, "2022-05-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-05-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-05-19")
	go m.MatchTxsIdxTime(d, 1, "2022-05-02")

	go m.MatchTxsIdxTime(d, 10, "2022-03-24")
	go m.MatchTxsIdxTime(d, 56, "2022-04-29")
	go m.MatchTxsIdxTime(d, 137, "2022-04-22")
	go m.MatchTxsIdxTime(d, 250, "2022-04-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-04-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-04-19")
	go m.MatchTxsIdxTime(d, 1, "2022-04-02")

	//go m.MatchTxsIdxTime(d, 10, "2022-02-24")
	go m.MatchTxsIdxTime(d, 56, "2022-03-29")
	go m.MatchTxsIdxTime(d, 137, "2022-03-22")
	go m.MatchTxsIdxTime(d, 250, "2022-03-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-03-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-03-19")
	go m.MatchTxsIdxTime(d, 1, "2022-03-02")

	//go m.MatchTxsIdxTime(d, 10, "2022-01-24")
	go m.MatchTxsIdxTime(d, 56, "2022-02-28")
	//go m.MatchTxsIdxTime(d, 137, "2022-02-22")
	go m.MatchTxsIdxTime(d, 250, "2022-02-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-02-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-02-19")
	go m.MatchTxsIdxTime(d, 1, "2022-02-02")

	//go m.MatchTxsIdxTime(d, 10, "2021-12-24")
	go m.MatchTxsIdxTime(d, 56, "2022-01-29")
	go m.MatchTxsIdxTime(d, 137, "2022-01-22")
	go m.MatchTxsIdxTime(d, 250, "2022-01-24")
	go m.MatchTxsIdxTime(d, 42161, "2022-01-02")
	go m.MatchTxsIdxTime(d, 43114, "2022-01-19")
	m.MatchTxsIdxTime(d, 1, "2022-01-02")

}
