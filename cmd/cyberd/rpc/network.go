package rpc

import (
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultIndexStats struct {
	Height     uint64 `json:"height"`
	Objects    uint64 `json:"cidsCount"`
	Cyberlinks uint64 `json:"linksCount"`
	Subjects   uint64 `json:"accountsCount"`
	Karma      uint64 `json:"totalKarma"`
}

func IndexStats(ctx *rpctypes.Context) (*ResultIndexStats, error) {
	stats := &ResultIndexStats{}
	stats.Height     = uint64(cyberdApp.LastBlockHeight())
	stats.Objects    = cyberdApp.CidsCount()
	stats.Cyberlinks = cyberdApp.LinksCount()
	stats.Subjects   = cyberdApp.AccsCount()
	stats.Karma      = cyberdApp.CurrentTotalKarma()
	return stats, nil
}
