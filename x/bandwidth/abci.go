package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

var accsToUpdate = make([]sdk.AccAddress, 0)

// user to recover and update bandwidth for accs with changed stake
func updateAccMaxBandwidth(ctx sdk.Context, meter types.BandwidthMeter) {
	for _, addr := range accsToUpdate {
		meter.UpdateAccMaxBandwidth(ctx, addr)
	}
	accsToUpdate = make([]sdk.AccAddress, 0)
}

// collect all addresses with updated stake
func CollectAddressesWithStakeChange() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	return func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
		if ctx.IsCheckTx() {
			return
		}
		if from != nil {
			accsToUpdate = append(accsToUpdate, from)
		}
		if to != nil {
			accsToUpdate = append(accsToUpdate, to)
		}
	}
}

// Used for 2 points:
// 1. Adjust credit price each `AdjustPricePeriod` blocks
// 2. For accs with updated on current block stake adjust max bandwidth. Why not update on `onCoinsTransfer`?
//  Cuz for some bound related operations, coins already added/reduced from acc, but not added to
//  validator\delegator pool.
func EndBlocker(ctx sdk.Context, pk params.Keeper, meter types.BandwidthMeter) {
	subspace, ok := pk.GetSubspace(types.DefaultParamspace)
	if !ok {
		panic("bandwidth params subspace is not found")
	}
	var paramset Params
	subspace.GetParamSet(ctx, &paramset)

	if ctx.BlockHeight() != 0 && ctx.BlockHeight()%paramset.AdjustPricePeriod == 0 {
		meter.AdjustPrice(ctx)
	}
	meter.CommitBlockBandwidth(ctx)
	updateAccMaxBandwidth(ctx, meter)
}
