package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/link"
)

type CyberdClient interface {
	// Cyberd Client Specification

	// returns current connected node chain id
	GetChainId() string

	// returns, if given link already exists
	IsLinkExist(from link.Cid, to link.Cid, addr sdk.AccAddress) (result bool, err error)

	// returns, if given link already exists
	IsAnyLinkExist(from link.Cid, to link.Cid) (result bool, err error)

	// get current bandwidth credits price
	// price 1 is price for situation, when all users use all their bandwidth (all blocks are filled for 100%)
	// if price < 1, that means blocks filled partially, thus allow more active users to do more transactions
	// if price > 1, that means network is under high load.
	GetCurrentBandwidthCreditPrice() (float64, error)

	// returns account bandwidth information for given account
	GetAccountBandwidth() (bandwidth.AcсBandwidth, error)

	// links two cids for given user
	// this method also should check, either cids are correct cids and given user is msg signer
	// do not wait till tx will be mined, just returns results from tx mempool check
	SubmitLinkSync(link link.Link) error

	// see `SubmitLinkAsync`. Links will be submitted as single tx with multiple msges.
	// `submitOnlyNew` - if true, only links not exists before(no one do the same link) will be submitted.
	SubmitLinksSync(links []link.Link, submitOnlyNew bool) error
}
