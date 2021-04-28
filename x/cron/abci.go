package cron

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/cron/keeper"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.EndBlocker(ctx)
}
