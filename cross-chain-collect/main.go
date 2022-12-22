package main

import (
	"app/aggregator"
	"app/config"
	"app/dao"
	"app/match"
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

type T struct {
	Id    string `db:"id"`
	Chain string `db:"chain"`
}

func gogo(d *dao.Dao) *[]match.MatchedId {
	t := &[]match.MatchedId{}
	stmt := "select a.id as src_id, a.match_id as dest_id from across a inner join common_cross_chain b on a.match_id = b.id and b.match_id is null"

	//t := &[]T{}
	//stmt := "select id, chain from synapse_v2 where direction = 'in' and to_chain = '-1'"
	err := d.DB().Select(&(*t), stmt)
	println(len(*t))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return t
}

func f(d *dao.Dao, t *[]match.MatchedId, start, batch_size int) {
	for i := 0; i < batch_size; i++ {
		_, err := d.DB().Exec("UPDATE common_cross_chain SET match_id = $1 WHERE id = $2", (*t)[i+start].SrcID, (*t)[i+start].DstID)
		if err != nil {
			fmt.Println(err)
			return
		}
		if i%500 == 0 {
			print("done: ", i+start)
		}
	}
	println("already: ", start+batch_size)
}

func main() {
	//d := dao.NewDao("postgres://yufeng:yufengblockSec888@192.168.3.155:8888/cross_chain?sslmode=disable")
	d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")

	/*go f(d, t, batch_size, batch_size)
	go f(d, t, 1*batch_size, batch_size)
	go f(d, t, 2*batch_size, batch_size)
	go f(d, t, 3*batch_size, batch_size)
	go f(d, t, 4*batch_size, batch_size)
	go f(d, t, 5*batch_size, batch_size)
	go f(d, t, 6*batch_size, batch_size)
	go f(d, t, 7*batch_size, batch_size)
	go f(d, t, 8*batch_size, batch_size)
	go f(d, t, 9*batch_size, batch_size)
	go f(d, t, 10*batch_size, batch_size)
	go f(d, t, 11*batch_size, batch_size)
	go f(d, t, 12*batch_size, batch_size)
	go f(d, t, 13*batch_size, batch_size)
	go f(d, t, 14*batch_size, batch_size)
	go f(d, t, 15*batch_size, batch_size)
	go f(d, t, 16*batch_size, batch_size)
	go f(d, t, 17*batch_size, batch_size)
	go f(d, t, 18*batch_size, batch_size)
	go f(d, t, 19*batch_size, batch_size)
	go f(d, t, 20*batch_size, batch_size)
	go f(d, t, 21*batch_size, batch_size)
	go f(d, t, 22*batch_size, batch_size)
	go f(d, t, 23*batch_size, batch_size)
	go f(d, t, 24*batch_size, batch_size)
	go f(d, t, 25*batch_size, batch_size)
	go f(d, t, 26*batch_size, batch_size)
	go f(d, t, 27*batch_size, batch_size)
	go f(d, t, 28*batch_size, batch_size)
	go f(d, t, 29*batch_size, batch_size)
	go f(d, t, 30*batch_size, batch_size)
	go f(d, t, 31*batch_size, batch_size)
	go f(d, t, 32*batch_size, batch_size)
	go f(d, t, 33*batch_size, batch_size)
	go f(d, t, 34*batch_size, batch_size)
	go f(d, t, 35*batch_size, batch_size)
	go f(d, t, 36*batch_size, batch_size)
	go f(d, t, 37*batch_size, batch_size)
	go f(d, t, 38*batch_size, batch_size)
	go f(d, t, 39*batch_size, batch_size)
	go f(d, t, 40*batch_size, batch_size)
	go f(d, t, 40*batch_size, batch_size)
	go f(d, t, 41*batch_size, batch_size)
	go f(d, t, 42*batch_size, batch_size)
	go f(d, t, 43*batch_size, batch_size)
	go f(d, t, 44*batch_size, batch_size)
	go f(d, t, 45*batch_size, batch_size)
	go f(d, t, 46*batch_size, batch_size)
	go f(d, t, 47*batch_size, batch_size)
	go f(d, t, 48*batch_size, batch_size)
	go f(d, t, 49*batch_size, batch_size)
	go f(d, t, 50*batch_size, batch_size)
	go f(d, t, 51*batch_size, batch_size)
	go f(d, t, 52*batch_size, batch_size)
	f(d, t, 53*batch_size, batch_size)*/
	m := &match.Match{}
	matched_pair, _ := m.GetMatchedPair(d, "Anyswap")

	/*d := dao.NewDao("postgres://cross_chain:cross_chain_blocksec666@192.168.3.155:8888/cross_chain?sslmode=disable")

	//m := &match.Match{}

	matched_pair, err := m.GetMatchedPair(d, "Anyswap")
	if err != nil {
		fmt.Println(err)
		return
	}*/

	batch_size := len(*matched_pair)

	m.BatchMatch(d, matched_pair, 0, batch_size)
	/*go m.BatchMatch(d, matched_pair, batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 2*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 3*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 4*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 5*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 6*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 7*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 8*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 9*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 10*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 11*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 12*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 13*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 14*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 15*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 16*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 17*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 18*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 19*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 20*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 21*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 22*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 23*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 24*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 25*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 26*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 27*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 28*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 29*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 30*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 31*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 32*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 33*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 34*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 35*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 36*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 37*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 38*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 39*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 40*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 41*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 42*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 43*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 44*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 45*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 46*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 47*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 48*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 49*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 50*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 51*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 52*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 53*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 54*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 55*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 56*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 57*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 58*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 59*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 60*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 61*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 62*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 63*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 64*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 65*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 66*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 67*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 68*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 69*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 70*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 71*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 72*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 73*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 74*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 75*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 76*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 77*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 78*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 79*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 80*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 81*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 82*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 83*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 84*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 85*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 86*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 87*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 88*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 89*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 90*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 91*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 92*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 93*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 94*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 95*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 96*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 97*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 98*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 99*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 100*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 101*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 102*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 103*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 104*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 105*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 106*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 107*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 108*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 109*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 110*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 111*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 112*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 113*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 114*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 115*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 116*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 117*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 118*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 119*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 120*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 121*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 122*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 123*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 124*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 125*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 126*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 127*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 128*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 129*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 130*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 131*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 132*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 133*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 134*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 135*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 136*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 137*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 138*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 139*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 140*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 141*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 142*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 143*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 144*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 145*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 146*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 147*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 148*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 149*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 150*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 151*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 152*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 153*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 154*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 155*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 156*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 157*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 158*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 159*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 160*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 161*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 162*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 163*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 164*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 165*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 166*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 167*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 168*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 169*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 170*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 171*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 172*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 173*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 174*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 175*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 176*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 177*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 178*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 178*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 179*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 180*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 181*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 182*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 183*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 184*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 185*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 186*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 187*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 188*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 189*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 190*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 191*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 192*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 193*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 194*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 195*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 196*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 197*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 198*batch_size, batch_size)
	go m.BatchMatch(d, matched_pair, 199*batch_size, batch_size)
	m.BatchMatch(d, matched_pair, 200*batch_size, batch_size)*/
}
