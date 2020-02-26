package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	cbd "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/bank/internal/types"
)

type Keeper interface {
	sdkbank.Keeper

	AddHook(types.CoinsTransferHook)

	SetStakingKeeper(types.StakingKeeper)
	SetSupplyKeeper(types.SupplyKeeper)

	GetAccountUnboundedStake(sdk.Context, sdk.AccAddress) int64
	GetAccountBoundedStake(sdk.Context, sdk.AccAddress) int64
	GetAccountTotalStake(sdk.Context, sdk.AccAddress) int64
	GetAccStakePercentage(sdk.Context, sdk.AccAddress) float64

	GetTotalSupply(sdk.Context) int64
}

type IndexedKeeper interface {
	Keeper

	Load(sdk.Context, sdk.Context)

	FixUserStake(ctx sdk.Context) bool
	UpdateStake(cbd.AccNumber, int64)
	GetTotalStakes() map[cbd.AccNumber]uint64
	EndBlocker(sdk.Context)
}
