package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)

const (
	// DefaultParamspace default name for parameter store
	DefaultParamspace = types.ModuleName
)

// ParamTable for bandwidth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

func (b *BaseRankKeeper) SetParams(ctx sdk.Context, params types.Params) {
	b.paramSpace.SetParamSet(ctx, &params)
}

func (b *BaseRankKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	b.paramSpace.GetParamSet(ctx, &params)
	return params
}
