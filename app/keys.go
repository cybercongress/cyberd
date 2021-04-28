package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"

	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/energy"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank"
	"github.com/cybercongress/go-cyber/x/cron"
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

	links          *sdk.KVStoreKey
	rank           *sdk.KVStoreKey
	bandwidth      *sdk.KVStoreKey
	energy         *sdk.KVStoreKey
	cron		   *sdk.KVStoreKey
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

		links:    sdk.NewKVStoreKey(link.StoreKey),
		rank:     sdk.NewKVStoreKey(rank.StoreKey),
		bandwidth:sdk.NewKVStoreKey(bandwidth.StoreKey),
		energy:	  sdk.NewKVStoreKey(energy.StoreKey),
		cron:	  sdk.NewKVStoreKey(cron.StoreKey),
		wasm:	  sdk.NewKVStoreKey(wasm.StoreKey),

		tStake:   sdk.NewTransientStoreKey(staking.TStoreKey),
		tParams:  sdk.NewTransientStoreKey(params.TStoreKey),
	}
}

func (k cyberdAppDbKeys) GetStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		k.main, k.auth, k.links, k.rank, k.stake, k.supply, k.gov,
		k.slashing, k.params, k.distr, k.bandwidth, k.mint, k.upgrade,
		k.evidence, k.wasm, k.energy, k.cron,
	}
}

func (k cyberdAppDbKeys) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		k.tStake, k.tParams,
	}
}
