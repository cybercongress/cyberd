package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
	"github.com/cybercongress/go-cyber/x/link"
)

func MsgBandwidthCosts(ctx sdk.Context, pk params.Keeper, msg sdk.Msg) int64 {
	subspace, ok := pk.GetSubspace(types.DefaultParamspace)
	if !ok {
		panic("bandwidth params subspace is not found")
	}
	var paramset Params
	subspace.GetParamSet(ctx, &paramset)

	switch msg.(type) {
	case link.Msg:
		linkMsg := msg.(link.Msg)
		return int64(len(linkMsg.Links)) * paramset.LinkMsgCost
	default:
		return paramset.NonLinkMsgCost
	}
}
