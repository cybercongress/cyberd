package wasm

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	wasmTypes "github.com/CosmWasm/go-cosmwasm/types"

	"github.com/cybercongress/go-cyber/x/link/keeper"
	"github.com/cybercongress/go-cyber/x/link/types"
)

var _ WasmQuerierInterface = WasmQuerier{}
var _ WasmMsgParserInterface = WasmMsgParser{}

//--------------------------------------------------

type WasmMsgParserInterface interface {
	Parse(contractAddr sdk.AccAddress, msg wasmTypes.CosmosMsg) ([]sdk.Msg, error)
	ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error)
}

type WasmMsgParser struct{}

func NewWasmMsgParser() WasmMsgParser {
	return WasmMsgParser{}
}

func (WasmMsgParser) Parse(_ sdk.AccAddress, _ wasmTypes.CosmosMsg) ([]sdk.Msg, error) { return nil, nil }

type CosmosMsg struct {
	Cyberlink *types.MsgCyberlink `json:"cyberlink,omitempty"`
}

func (WasmMsgParser) ParseCustom(contractAddr sdk.AccAddress, data json.RawMessage) ([]sdk.Msg, error) {
	var sdkMsg CosmosMsg
	err := json.Unmarshal(data, &sdkMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to parse link custom msg")
	}

	if sdkMsg.Cyberlink != nil {
		return []sdk.Msg{*sdkMsg.Cyberlink}, sdkMsg.Cyberlink.ValidateBasic()
	}

	return nil, sdkerrors.Wrap(wasm.ErrInvalidMsg, "Unknown variant of Link")
}

//--------------------------------------------------

type WasmQuerierInterface interface {
	Query(ctx sdk.Context, request wasmTypes.QueryRequest) ([]byte, error)
	QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error)
}

type WasmQuerier struct {
	keeper.GraphKeeper
}

func NewWasmQuerier(keeper keeper.GraphKeeper) WasmQuerier {
	return WasmQuerier{keeper}
}

func (WasmQuerier) Query(_ sdk.Context, _ wasmTypes.QueryRequest) ([]byte, error) { return nil, nil }

type CosmosQuery struct {
	CidsCount     *struct{} `json:"cids_count,omitempty"`
	LinksCount    *struct{} `json:"links_count,omitempty"`
}

type CidsCountQueryResponse struct {
	CidsCount 	  uint64 `json:"cids_count"`
}

type LinksCountQueryResponse struct {
	LinksCount 	  uint64 `json:"links_count"`
}

func (querier WasmQuerier) QueryCustom(ctx sdk.Context, data json.RawMessage) ([]byte, error) {
	var query CosmosQuery
	err := json.Unmarshal(data, &query)

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	var bz []byte

	if query.CidsCount != nil {
		count := querier.GraphKeeper.GetCidsCount(ctx)
		bz, err = json.Marshal(CidsCountQueryResponse{CidsCount: count})
	} else if query.LinksCount != nil {
		count := querier.GraphKeeper.GetLinksCount(ctx)
		bz, err = json.Marshal(LinksCountQueryResponse{LinksCount: count})
	} else {
		return nil, sdkerrors.ErrInvalidRequest
	}

	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}