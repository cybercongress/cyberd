package rank

import (
	"github.com/cybercongress/go-cyber/x/rank/exported"
	"github.com/cybercongress/go-cyber/x/rank/internal/keeper"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)

const (
	ModuleName = types.ModuleName
	DefaultParamspace = keeper.DefaultParamspace
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
	QuerierRoute           = types.QuerierRoute
	QueryParameters        = types.QueryParameters
	QueryCalculationWindow = types.QueryCalculationWindow
	QueryDampingFactor     = types.QueryDampingFactor
	QueryTolerance         = types.QueryTolerance
	CPU        = types.CPU
	GPU        = types.GPU
)

var (
	// keeper
	NewStateKeeper = keeper.NewStateKeeper
	ParamKeyTable  = keeper.ParamKeyTable
	NewQuerier          = keeper.NewQuerier
	// types
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewDefaultParams    = types.NewDefaultParams
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)

type (
	// exported
	Keeper      = exported.Keeper
	StateKeeper = exported.StateKeeper

	// types
	Params       = types.Params
	GenesisState = types.GenesisState
	ComputeUnit  = types.ComputeUnit
)
