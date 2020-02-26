package rpc

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/tendermint/tendermint/rpc/lib/types"
)

func IsLinkExist(ctx *rpctypes.Context, from string, to string, address string) (bool, error) {

	if len(address) == 0 {
		return cyberdApp.IsLinkExist(link.Cid(from), link.Cid(to), nil), nil
	}

	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return false, err
	}

	return cyberdApp.IsLinkExist(link.Cid(from), link.Cid(to), accAddress), nil
}

type LinkRequest struct {
	Fee        auth.StdFee `json:"fee"`
	Msgs       []link.Msg  `json:"msgs"`
	Signatures []Signature `json:"signatures"`
	Memo       string      `json:"memo"`
}

func (r LinkRequest) GetFee() auth.StdFee        { return r.Fee }
func (r LinkRequest) GetSignatures() []Signature { return r.Signatures }
func (r LinkRequest) GetMemo() string            { return r.Memo }
func (r LinkRequest) GetMsgs() []types.Msg {
	msgs := make([]types.Msg, 0, len(r.Msgs))
	for _, msg := range r.Msgs {
		msgs = append(msgs, msg)
	}
	return msgs
}

func UnmarshalLinkRequestFn(reqBytes []byte) (TxRequest, error) {
	var req LinkRequest
	err := json.Unmarshal(reqBytes, &req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
