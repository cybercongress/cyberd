package rpc

import (
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultNetworkLoad struct {
	Load float64 `amino:"unsafe" json:"load"`
}

func CurrentNetworkLoad(ctx *rpctypes.Context) (*ResultNetworkLoad, error) {
	return &ResultNetworkLoad{cyberdApp.CurrentNetworkLoad()}, nil
}
