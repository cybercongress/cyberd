package cron

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/cron/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	//k.SetJobs(ctx, data.Jobs)
}

func ExportGenesis(ctx sdk.Context, k Keeper) (data GenesisState) {
	params := k.GetParams(ctx)
	jobs := k.GetAllJobs(ctx)

	return types.NewGenesisState(params, jobs)
}
