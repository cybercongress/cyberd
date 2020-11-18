package rank

import (
	"github.com/cybercongress/go-cyber/x/rank/keeper"
	"github.com/cybercongress/go-cyber/x/rank/types"
	"github.com/cybercongress/go-cyber/x/rank/wasm"
)

const (
	ModuleName 			   = types.ModuleName
	DefaultParamspace 	   = types.DefaultParamspace
	StoreKey   			   = types.StoreKey
	QuerierRoute           = types.QuerierRoute
	QueryParameters        = types.QueryParameters
	CPU        			   = types.CPU
	GPU        			   = types.GPU
)

var (
	NewKeeper 			= keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewWasmQuerier      = wasm.NewWasmQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	DefaultParams       = types.DefaultParams
	ModuleCdc			= types.ModuleCdc
)

type (
	StateKeeper  = keeper.StateKeeper
	GenesisState = types.GenesisState
	Params       = types.Params
	ComputeUnit  = types.ComputeUnit
)

