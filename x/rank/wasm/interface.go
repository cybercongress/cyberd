package wasm

import (
	"encoding/json"
	"strconv"

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
	RankValue *QueryRankValueParams `json:"rank_value,omitempty"`
}

// string used only to start work from

type QueryRankValueParams struct {
	CidNumber string `json:"cid_number"`
}

type RankQueryResponse struct {
	Rank string `json:"rank_value"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.RankValue != nil {
		number, _ := strconv.ParseUint(query.RankValue.CidNumber, 10, 64)
		rank := querier.keeper.GetRankValues(number)
		bz, err = json.Marshal(RankQueryResponse{Rank: strconv.FormatUint(rank, 10)})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}