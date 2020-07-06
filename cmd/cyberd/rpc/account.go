package rpc

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	bw "github.com/cybercongress/go-cyber/x/bandwidth"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
)

type ResultAccount struct {
	Account auth.BaseAccount `json:"account"`
}

func Account(ctx *rpctypes.Context, address string) (*ResultAccount, error) {
	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	account := cyberdApp.Account(accAddress)

	if account == nil {
		return nil, errors.New("account not found")
	}

	return &ResultAccount{
		auth.BaseAccount{
			Address:       account.GetAddress(),
			Coins:         account.GetCoins(),
			PubKey:        account.GetPubKey(),
			AccountNumber: account.GetAccountNumber(),
			Sequence:      account.GetSequence()},
	}, nil
}

func AccountBandwidth(ctx *rpctypes.Context, address string) (*bw.AccountBandwidth, error) {

	accAddress, err := types.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	accBdwth := cyberdApp.AccountBandwidth(accAddress)
	return &accBdwth, nil
}
