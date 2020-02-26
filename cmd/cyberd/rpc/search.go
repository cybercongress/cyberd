package rpc

import (
	"github.com/cybercongress/go-cyber/app"
	"github.com/tendermint/tendermint/rpc/lib/types"
)

type ResultSearch struct {
	Cids       []app.RankedCid `json:"cids"`
	TotalCount int             `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
}

func Search(ctx *rpctypes.Context, cid string, page, perPage int) (*ResultSearch, error) {
	if perPage == 0 {
		perPage = 100
	}
	links, totalSize, err := cyberdApp.Search(cid, page, perPage)
	return &ResultSearch{links, totalSize, page, perPage}, err
}
