package plugins

import (
	"encoding/json"
	"fmt"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"
	wasm "github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank"
)

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type Querier struct {
	Queriers map[string]WasmQuerierInterface
}

func NewQuerier() Querier {
	return Querier{
		Queriers: make(map[string]WasmQuerierInterface),
	}
}

type WasmCustomQuery struct {
	Route     string          `json:"route"`
	QueryData json.RawMessage `json:"query_data"`
}

const (
	WasmQueryRouteLink     = link.ModuleName
	WasmQueryRouteRank     = rank.ModuleName
)

func (q Querier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var customQuery WasmCustomQuery
	err := json.Unmarshal(data, &customQuery)
	fmt.Println("[!] Wasm query routed to module: ", customQuery.Route)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if querier, ok := q.Queriers[customQuery.Route]; ok {
		return querier.QueryCustom(ctx, customQuery.QueryData)
	} else {
		return nil, sdkerrors.Wrap(wasm.ErrQueryFailed, customQuery.Route)
	}
}
