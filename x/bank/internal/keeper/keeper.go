package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/go-cyber/types/coin"
	"github.com/cybercongress/go-cyber/x/bank/internal/types"
)

type Keeper struct {
	sdkbank.Keeper

	accountKeeper types.AccountKeeper
	stakingKeeper types.StakingKeeper
	supplyKeeper  types.SupplyKeeper

	coinsTransferHooks []types.CoinsTransferHook
}

func NewKeeper(accountKeeper types.AccountKeeper, subspace params.Subspace, blacklistedAddrs map[string]bool) *Keeper {

	return &Keeper{
		Keeper:        sdkbank.NewBaseKeeper(accountKeeper, subspace, blacklistedAddrs),
		accountKeeper: accountKeeper,
		coinsTransferHooks: make([]types.CoinsTransferHook, 0),
	}
}

func (k *Keeper) AddHook(hook types.CoinsTransferHook) {
	k.coinsTransferHooks = append(k.coinsTransferHooks, hook)
}

func (k *Keeper) SetStakingKeeper(sk types.StakingKeeper) {
	k.stakingKeeper = sk
}

func (k *Keeper) SetSupplyKeeper(sk types.SupplyKeeper) {
	k.supplyKeeper = sk
}

/* Override methods */
// sdk accountKeeper keeper is not interface yet
func (k Keeper) AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	coins, err := k.Keeper.AddCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return coins, err
}

func (k Keeper) SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error) {
	coins, err := k.Keeper.SubtractCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return coins, err
}

func (k Keeper) SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error {
	err := k.Keeper.SetCoins(ctx, addr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, addr)
	}
	return err
}

func (k Keeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {

	err := k.Keeper.SendCoins(ctx, fromAddr, toAddr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, fromAddr, toAddr)
	}
	return err
}

func (k Keeper) InputOutputCoins(
	ctx sdk.Context, inputs []sdkbank.Input, outputs []sdkbank.Output) error {
	err := k.Keeper.InputOutputCoins(ctx, inputs, outputs)
	if err == nil {
		for _, i := range inputs {
			k.onCoinsTransfer(ctx, i.Address, nil)
		}
		for _, j := range outputs {
			k.onCoinsTransfer(ctx, nil, j.Address)
		}
	}
	return err
}

func (k Keeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	err := k.Keeper.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, delegatorAddr, nil)
	}
	return err
}

func (k Keeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	err := k.Keeper.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt)
	if err == nil {
		k.onCoinsTransfer(ctx, nil, delegatorAddr)
	}
	return err
}

func (k Keeper) onCoinsTransfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress) {
	for _, hook := range k.coinsTransferHooks {
		hook(ctx, from, to)
	}
}

/* Query methods */

func (k Keeper) GetAccountUnboundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return 0
	}
	return acc.GetCoins().AmountOf(coin.CYB).Int64()
}

func (k Keeper) GetAccountBoundedStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	delegations := k.stakingKeeper.GetAllDelegatorDelegations(ctx, addr)
	boundedStake := int64(0)
	for _, del := range delegations {
		boundedStake += del.Shares.TruncateInt64()
	}
	return boundedStake
}

// Returns bounded plus unbounded account tokens
func (k Keeper) GetAccountTotalStake(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return k.GetAccountUnboundedStake(ctx, addr) + k.GetAccountBoundedStake(ctx, addr)
}

func (k Keeper) GetAccStakePercentage(ctx sdk.Context, addr sdk.AccAddress) float64 {
	a := k.GetAccountTotalStake(ctx, addr)
	aFloat := float64(a)

	b := k.GetTotalSupply(ctx)
	bFloat := float64(b)

	c := aFloat / bFloat
	return c
}

func (k Keeper) GetTotalSupply(ctx sdk.Context) int64 {
	keeperSupply := k.supplyKeeper.GetSupply(ctx)
	return keeperSupply.GetTotal().AmountOf(coin.CYB).Int64()
}
