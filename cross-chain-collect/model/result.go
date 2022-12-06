package model

import (
	"database/sql"
	"math/big"
)

type Result struct {
	Chain           string   `db:"chain"`
	Number          uint64   `db:"number"`
	Index           uint64   `db:"index"`
	Hash            string   `db:"hash"`
	ActionId        uint64   `db:"action_id"`
	Project         string   `db:"project"`
	Contract        string   `db:"contract"`
	Direction       string   `db:"direction"`
	FromChainId     *big.Int `db:"-"`
	WrapFromChainId *BigInt  `db:"from_chain_id"`
	FromAddress     string   `db:"from_address"`
	ToChainId       *big.Int `db:"-"`
	WrapToChainId   *BigInt  `db:"to_chain_id"`
	ToAddress       string   `db:"to_address"`
	Token           string   `db:"token"`
	Amount          *big.Int `db:"-"`
	WrapAmount      *BigInt  `db:"amount"`
	MatchTag        string   `db:"match_tag"`
	Detail          []byte   `db:"detail"`
}

type Results []*Result

type Data struct {
	Id              int64         `db:"id"`
	MatchId         sql.NullInt64 `db:"match_id"`
	Chain           string        `db:"chain"`
	Number          uint64        `db:"number"`
	Index           uint64        `db:"index"`
	Hash            string        `db:"hash"`
	ActionId        uint64        `db:"action_id"`
	Project         string        `db:"project"`
	Contract        string        `db:"contract"`
	Direction       string        `db:"direction"`
	FromChainId     *big.Int      `db:"-"`
	WrapFromChainId *BigInt       `db:"from_chain_id"`
	FromAddress     string        `db:"from_address"`
	ToChainId       *big.Int      `db:"-"`
	WrapToChainId   *BigInt       `db:"to_chain_id"`
	ToAddress       string        `db:"to_address"`
	Token           string        `db:"token"`
	Amount          *big.Int      `db:"-"`
	WrapAmount      *BigInt       `db:"amount"`
	MatchTag        string        `db:"match_tag"`
	Detail          []byte        `db:"detail"`
}

type Data_ struct {
	Id              int64         `db:"id"`
	MatchId         sql.NullInt64 `db:"match_id"`
	Chain           string        `db:"chain"`
	Number          uint64        `db:"block_number"`
	Index           uint64        `db:"tx_index"`
	Hash            string        `db:"tx_hash"`
	ActionId        uint64        `db:"action_id"`
	Contract        string        `db:"contract"`
	Direction       string        `db:"direction"`
	Timestamp       string        `db:"timestamp"`
	FromChainId     *big.Int      `db:"-"`
	WrapFromChainId int           `db:"from_chain"`
	FromAddress     string        `db:"from_address"`
	ToChainId       *big.Int      `db:"-"`
	WrapToChainId   int           `db:"to_chain"`
	ToAddress       string        `db:"to_address"`
	Token           string        `db:"token"`
	Amount          *big.Int      `db:"-"`
	WrapAmount      *BigInt       `db:"amount"`
	SrcTxHash       string        `db:"src_tx_hash"`
}
