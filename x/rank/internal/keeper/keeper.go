package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/go-cyber/x/rank/exported"
)

var _ exported.Keeper = &BaseRankKeeper{}

type BaseRankKeeper struct {
	paramSpace *params.Subspace
}

func NewBaseRankKeeper(paramSpace *params.Subspace) *BaseRankKeeper {
	return &BaseRankKeeper{
		paramSpace: paramSpace,
	}
}
