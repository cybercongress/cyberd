package app

import (
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmwasm/wasmd/x/wasm"
)

type cyberdAppDbKeys struct {
	main           *sdk.KVStoreKey
	auth           *sdk.KVStoreKey
	stake          *sdk.KVStoreKey
	supply         *sdk.KVStoreKey
	distr          *sdk.KVStoreKey
	gov            *sdk.KVStoreKey
	params         *sdk.KVStoreKey
	slashing       *sdk.KVStoreKey
	mint           *sdk.KVStoreKey
	upgrade  	   *sdk.KVStoreKey
	evidence	   *sdk.KVStoreKey

	cidNum         *sdk.KVStoreKey
	cidNumReverse  *sdk.KVStoreKey
	links          *sdk.KVStoreKey
	rank           *sdk.KVStoreKey
	accBandwidth   *sdk.KVStoreKey
	blockBandwidth *sdk.KVStoreKey

	wasm 		   *sdk.KVStoreKey

	tStake         *sdk.TransientStoreKey
	tParams        *sdk.TransientStoreKey
}

func NewCyberdAppDbKeys() cyberdAppDbKeys {
	return cyberdAppDbKeys{
		main:     sdk.NewKVStoreKey(bam.MainStoreKey),
		auth:      sdk.NewKVStoreKey(auth.StoreKey),
		stake:    sdk.NewKVStoreKey(staking.StoreKey),
		supply:   sdk.NewKVStoreKey(supply.StoreKey),
		distr:    sdk.NewKVStoreKey(distr.StoreKey),
		gov:      sdk.NewKVStoreKey(gov.StoreKey),
		slashing: sdk.NewKVStoreKey(slashing.StoreKey),
		params:   sdk.NewKVStoreKey(params.StoreKey),
		mint:     sdk.NewKVStoreKey(mint.StoreKey),
		upgrade:  sdk.NewKVStoreKey(upgrade.StoreKey),
		evidence: sdk.NewKVStoreKey(evidence.StoreKey),

		cidNum:         sdk.NewKVStoreKey("cid_index"), // TODO
		cidNumReverse:  sdk.NewKVStoreKey("cid_index_reverse"),
		links:          sdk.NewKVStoreKey("cyberlinks"),
		rank:           sdk.NewKVStoreKey("rank"),
		accBandwidth:   sdk.NewKVStoreKey("acc_bandwidth"),
		blockBandwidth: sdk.NewKVStoreKey("block_spent_bandwidth"),

		wasm:			sdk.NewKVStoreKey(wasm.StoreKey),

		tStake:   sdk.NewTransientStoreKey(staking.TStoreKey),
		tParams:  sdk.NewTransientStoreKey(params.TStoreKey),
	}
}

func (k cyberdAppDbKeys) GetStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		k.main, k.auth, k.cidNum, k.cidNumReverse, k.links, k.rank, k.stake, k.supply, k.gov,
		k.slashing, k.params, k.distr, k.accBandwidth, k.blockBandwidth, k.mint, k.upgrade, k.evidence, k.wasm,
	}
}

func (k cyberdAppDbKeys) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		k.tStake, k.tParams,
	}
}
