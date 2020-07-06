package rpc

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

func StakingValidators(ctx *rpctypes.Context, page, limit int, status string) ([]sdk.Validator, error) {

	queryValsParams := staking.NewQueryValidatorsParams(page, limit, status)
	bz, err := codec.MarshalJSON(queryValsParams)
	if err != nil {
		return nil, err
	}

	respQuery := cyberdApp.Query(abci.RequestQuery{
		Path:  "custom/staking/validators",
		Data: bz,
	})

	validators := make([]sdk.Validator, 0)
	err = codec.UnmarshalJSON(respQuery.Value, &validators)
	if err != nil {
		return nil, err
	}

	return validators, nil
}

func StakingPool(ctx *rpctypes.Context) (staking.Pool, error) {

	respQuery := cyberdApp.Query(abci.RequestQuery{
		Path:  "custom/staking/pool",
		Prove: false,
	})

	pool := staking.Pool{}
	err := codec.UnmarshalJSON(respQuery.Value, &pool)
	if err != nil {
		return pool, err
	}

	return pool, nil
}
