package bandwidth

import (
	"encoding/json"

	"github.com/cybercongress/go-cyber/x/bandwidth/client/rest"
	"github.com/cybercongress/go-cyber/x/bandwidth/client/cli"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// app module Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return ModuleName }

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { RegisterCodec(cdc) }

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// Once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// get the root tx command of this module
func (AppModuleBasic) GetTxCmd(_ *codec.Codec) *cobra.Command { return nil }

// get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

type AppModule struct {
	AppModuleBasic
	AccBandwidthKeeper        Keeper
	BlockSpentBandwidthKeeper BlockSpentBandwidthKeeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(
	accBandwidthKeeper Keeper,
	blockSpentBandwidthKeeper BlockSpentBandwidthKeeper) AppModule {
	return AppModule{
		AppModuleBasic:            AppModuleBasic{},
		AccBandwidthKeeper:        accBandwidthKeeper,
		BlockSpentBandwidthKeeper: blockSpentBandwidthKeeper,
	}
}

func (AppModule) Name() string { return ModuleName }

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string           { return RouterKey }
func (am AppModule) NewHandler() sdk.Handler { return nil }

func (am AppModule) QuerierRoute() string { return QuerierRoute }

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.AccBandwidthKeeper)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}
func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.AccBandwidthKeeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
