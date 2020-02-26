package bandwidth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

// Genesis accounts should contains fully restored bandwidth on block 0
func InitGenesis(ctx sdk.Context, handler types.BandwidthMeter,
	keeper Keeper, addresses []sdk.AccAddress, data GenesisState) {

	keeper.SetParams(ctx, data.Params)
	for _, address := range addresses {
		accMaxBw := handler.GetAccMaxBandwidth(ctx, address)
		keeper.SetAccBandwidth(ctx, types.NewGenesisAccBandwidth(address, accMaxBw))
	}
}

func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetParams(ctx))
}
