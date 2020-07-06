package rpc

import (
	"encoding/hex"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type Signature struct {
	PubKey        secp256k1.PubKeySecp256k1 `json:"pub_key"` // optional
	Signature     []byte                    `json:"signature"`
	AccountNumber uint64                    `json:"account_number"`
	Sequence      uint64                    `json:"sequence"`
}

type TxRequest interface {
	GetFee() auth.StdFee
	GetSignatures() []Signature
	GetMemo() string
	GetMsgs() []sdk.Msg
}

type UnmarshalTxRequestFn func([]byte) (TxRequest, error)

// takes signed messages, construct tx, invoke local client.
// there is no way to assemble right amino tx object on js, so we need to do it on go side.
func SignedMsgHandler(unmarshal UnmarshalTxRequestFn) func(ctx *rpctypes.Context, dataHex string) (*ctypes.ResultBroadcastTx, error) {
	return func(ctx *rpctypes.Context, dataHex string) (*ctypes.ResultBroadcastTx, error) {
		data, err := hex.DecodeString(dataHex)
		if err != nil {
			return nil, err
		}

		txReq, err := unmarshal(data)
		if err != nil {
			return nil, err
		}
		// BUILDING COSMOS SDK TX
		signatures := make([]auth.StdSignature, 0, len(txReq.GetSignatures()))
		for _, sig := range txReq.GetSignatures() {
			stdSig := auth.StdSignature{PubKey: sig.PubKey, Signature: sig.Signature}
			signatures = append(signatures, stdSig)
		}

		stdTx := auth.StdTx{Msgs: txReq.GetMsgs(), Fee: txReq.GetFee(), Signatures: signatures, Memo: txReq.GetMemo()}

		stdTxBytes, err := codec.MarshalBinaryLengthPrefixed(stdTx)
		if err != nil {
			return nil, err
		}

		return core.BroadcastTxSync(&rpctypes.Context{}, stdTxBytes)
	}
}
