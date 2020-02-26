package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)

type Keeper interface {
	SetParams(sdk.Context, types.Params)
	GetParams(sdk.Context) types.Params
}

type StateKeeper interface {
	Keeper

	Load(sdk.Context, log.Logger)
	BuildSearchIndex(log.Logger) types.SearchIndex

	EndBlocker(sdk.Context, log.Logger)

	Search(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error)
	Top(page, perPage int) ([]types.RankedCidNumber, int, error)

	GetRankValue(link.CidNumber) float64
	GetNetworkRankHash() []byte

	GetLastCidNum() link.CidNumber
	GetMerkleTree() *merkle.Tree
	GetIndexError() error
}
