package wasm

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/cybercongress/go-cyber/x/rank/keeper"
)

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

var _ WasmQuerierInterface = WasmQuerier{}

type WasmQuerier struct {
	keeper *keeper.StateKeeper
}

func NewWasmQuerier(keeper *keeper.StateKeeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	RankValueByID     *QueryRankValueByIdParams  `json:"rank_value_by_id,omitempty"`
	RankValueByCid    *QueryRankValueByCidParams `json:"rank_value_by_cid,omitempty"`
}

type QueryRankValueByIdParams struct {
	CidNumber uint64 `json:"cid_number"`
}

type QueryRankValueByCidParams struct {
	Cid 	  string `json:"cid"`
}

type RankQueryResponse struct {
	Rank 	  uint64 `json:"rank_value"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.RankValueByID != nil {
		rank := querier.keeper.GetRankValueByNumber(query.RankValueByID.CidNumber)
		bz, err = json.Marshal(RankQueryResponse{Rank: rank})
	} else if query.RankValueByCid != nil {
		rank := querier.keeper.GetRankValueByCid(ctx, query.RankValueByCid.Cid)
		bz, err = json.Marshal(RankQueryResponse{Rank: rank})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}