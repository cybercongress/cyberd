package bandwidth

import (
	"github.com/cybercongress/go-cyber/x/bandwidth/exported"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/keeper"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute            = types.QuerierRoute
	QueryParameters         = types.QueryParameters
	QueryDesirableBandwidth = types.QueryDesirableBandwidth
	QueryMaxBlockBandwidth  = types.QueryMaxBlockBandwidth
	QueryRecoveryPeriod     = types.QueryRecoveryPeriod
	QueryAdjustPricePeriod  = types.QueryAdjustPricePeriod
	QueryBaseCreditPrice    = types.QueryBaseCreditPrice
	QueryTxCost             = types.QueryTxCost
	QueryLinkMsgCost        = types.QueryLinkMsgCost
	QueryNonLinkMsgCost     = types.QueryNonLinkMsgCost
)

type (
	Keeper                    = exported.Keeper
	BlockSpentBandwidthKeeper = exported.BlockSpentBandwidthKeeper

	Meter        = types.BandwidthMeter
	AcсBandwidth = types.AcсBandwidth
	Params       = types.Params
	GenesisState = types.GenesisState
)

var (
	ModuleCdc           = types.ModuleCdc
	RegisterCodec       = types.RegisterCodec
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
	NewDefaultParams    = types.NewDefaultParams
	NewQuerier          = keeper.NewQuerier

	NewAccBandwidthKeeper        = keeper.NewAccBandwidthKeeper
	NewBlockSpentBandwidthKeeper = keeper.NewBlockSpentBandwidthKeeper
	ParamKeyTable                = keeper.ParamKeyTable
)
