package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryParams 	 = "params"
	QueryJob         = "job"
	QueryJobStats    = "job_stats"
	QueryJobs   	 = "jobs"
)

type QueryJobParams struct {
	Contract 		sdk.AccAddress
	Creator         sdk.AccAddress
	Label           string
}

func NewQueryJobParams(contract, creator sdk.AccAddress, label string) QueryJobParams {
	return QueryJobParams{
		Contract: contract,
		Creator:  creator,
		Label: 	  label,
	}
}