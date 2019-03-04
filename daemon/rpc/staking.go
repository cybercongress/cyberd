package rpc

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdk "github.com/cosmos/cosmos-sdk/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func StakingValidators() ([]sdk.Validator, error) {

	respQuery := cyberdApp.Query(abci.RequestQuery{
		Path:  "custom/staking/validators",
		Prove: false,
	})

	validators := make([]sdk.Validator, 0)
	err := codec.UnmarshalJSON(respQuery.Value, &validators)
	if err != nil {
		return nil, err
	}

	return validators, nil
}

func StakingPool() (staking.Pool, error) {

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
