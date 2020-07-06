package rpc

import (
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultBandwidthPrice struct {
	Price float64 `amino:"unsafe" json:"price"`
}

func CurrentBandwidthPrice(ctx *rpctypes.Context) (*ResultBandwidthPrice, error) {
	return &ResultBandwidthPrice{cyberdApp.CurrentBandwidthPrice()}, nil
}
