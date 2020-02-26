package app

import (
	"errors"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/go-cyber/merkle"
	cbd "github.com/cybercongress/go-cyber/types"
	bw "github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/link"
	abci "github.com/tendermint/tendermint/abci/types"
)

type RankedCid struct {
	Cid  link.Cid `json:"cid"`
	Rank float64  `amino:"unsafe" json:"rank"`
}

func (app *CyberdApp) RpcContext() sdk.Context {
	return app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})
}

func (app *CyberdApp) Search(cid string, page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()
	cidNumber, exists := app.cidNumKeeper.GetCidNumber(ctx, link.Cid(cid))
	if !exists || cidNumber > app.rankStateKeeper.GetLastCidNum() {
		return nil, 0, errors.New("no such cid found")
	}

	rankedCidNumbers, size, err := app.rankStateKeeper.Search(cidNumber, page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.cidNumKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Top(page, perPage int) ([]RankedCid, int, error) {

	ctx := app.RpcContext()

	rankedCidNumbers, size, err := app.rankStateKeeper.Top(page, perPage)

	if err != nil {
		return nil, size, err
	}

	result := make([]RankedCid, 0, len(rankedCidNumbers))
	for _, c := range rankedCidNumbers {
		result = append(result, RankedCid{Cid: app.cidNumKeeper.GetCid(ctx, c.GetNumber()), Rank: c.GetRank()})
	}

	return result, size, nil
}

func (app *CyberdApp) Rank(cid string, proof bool) (float64, []merkle.Proof, error) {

	cidNumber, exists := app.cidNumKeeper.GetCidNumber(app.RpcContext(), link.Cid(cid))
	if !exists || cidNumber > app.rankStateKeeper.GetLastCidNum() {
		return 0, nil, errors.New("no such cid found")
	}

	rankValue := app.rankStateKeeper.GetRankValue(cidNumber)

	if proof {
		proofs := make([]merkle.Proof, 0, int64(math.Log2(float64(app.CidsCount()))))
		proofs = app.rankStateKeeper.GetMerkleTree().GetIndexProofs(int(cidNumber))
		return rankValue, proofs, nil
	}
	return rankValue, nil, nil
}

func (app *CyberdApp) Account(address sdk.AccAddress) auth.Account {
	return app.accountKeeper.GetAccount(app.RpcContext(), address)
}

func (app *CyberdApp) AccountBandwidth(address sdk.AccAddress) bw.AcсBandwidth {
	return app.bandwidthMeter.GetCurrentAccBandwidth(app.RpcContext(), address)
}

func (app *CyberdApp) IsLinkExist(from link.Cid, to link.Cid, address sdk.AccAddress) bool {

	ctx := app.RpcContext()
	fromNumber, fromExist := app.cidNumKeeper.GetCidNumber(ctx, from)
	toNumber, toExists := app.cidNumKeeper.GetCidNumber(ctx, to)

	if fromExist && toExists {
		if address != nil {
			acc := app.accountKeeper.GetAccount(ctx, address)
			if acc != nil {
				accNumber := cbd.AccNumber(acc.GetAccountNumber())
				return app.linkIndexedKeeper.IsLinkExist(link.NewLink(fromNumber, toNumber, accNumber))
			}
		} else {
			return app.linkIndexedKeeper.IsAnyLinkExist(fromNumber, toNumber)
		}
	}
	return false
}

func (app *CyberdApp) CurrentBandwidthPrice() float64 {
	return app.bandwidthMeter.GetCurrentCreditPrice()
}

func (app *CyberdApp) CidsCount() uint64 {
	return app.mainKeeper.GetCidsCount(app.RpcContext())
}

func (app *CyberdApp) LinksCount() uint64 {
	return app.mainKeeper.GetLinksCount(app.RpcContext())
}

func (app *CyberdApp) AccsCount() uint64 {
	return app.accountKeeper.GetNextAccountNumber(app.RpcContext())
}

func (app *CyberdApp) CurrentLinkedBandwidth() uint64 {
	return app.bandwidthMeter.GetCurrentBandwidthLinked()
}
