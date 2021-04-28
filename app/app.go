package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/abci/version"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	"github.com/CosmWasm/wasmd/x/wasm"

	"github.com/cybercongress/go-cyber/util"
	bandwidth "github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/cron"
	cyberbank "github.com/cybercongress/go-cyber/x/cyberbank"
	"github.com/cybercongress/go-cyber/x/energy"
	link "github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank"
)

const (
	appName = "cyber"
)

// default home directories for expected binaries
var (
	CLIDir       = ".cyberdcli"
	NodeDir      = ".cyberd"

	DefaultCLIHome = os.ExpandEnv("$HOME/") + CLIDir
	DefaultNodeHome = os.ExpandEnv("$HOME/") + NodeDir

	ProposalsEnabled = "false"
	EnableSpecificProposals = ""
)

func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

var (

	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler, upgradeclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},

		cyberbank.AppModuleBasic{},
		link.AppModuleBasic{},
		bandwidth.AppModuleBasic{},
		rank.AppModuleBasic{},
		energy.AppModuleBasic{},
		cron.AppModuleBasic{},
		wasm.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		energy.EnergyPoolName:     nil,
	}
)

type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txDecoder      sdk.TxDecoder
	invCheckPeriod uint

	dbKeys cyberdAppDbKeys

	subspaces map[string]params.Subspace

	accountKeeper      auth.AccountKeeper
	bankKeeper		   bank.Keeper
	supplyKeeper       supply.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	mintKeeper         mint.Keeper
	distrKeeper        distr.Keeper
	govKeeper          gov.Keeper
	paramsKeeper       params.Keeper
	crisisKeeper       crisis.Keeper
	upgradeKeeper      upgrade.Keeper
	evidenceKeeper     *evidence.Keeper

	cyberbank          *cyberbank.Keeper
	bandwidthMeter     *bandwidth.BandwidthMeter
	graphKeeper 	   link.GraphKeeper
	indexKeeper        *link.IndexKeeper
	rankKeeper    	   *rank.StateKeeper
	energyKeeper       energy.Keeper
	cronKeeper 		   cron.Keeper
	wasmKeeper     	   wasm.Keeper

	latestBlockHeight int64

	mm *module.Manager
}

type WasmWrapper struct {
	Wasm wasm.Config `mapstructure:"wasm"`
}

func NewCyberdApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, enabledProposals []wasm.ProposalType, skipUpgradeHeights map[int64]bool,
	computeUnit rank.ComputeUnit, allowSearch bool,
	baseAppOptions ...func(*baseapp.BaseApp),
) *CyberdApp {

	cdc := MakeCodec()
	txDecoder := auth.DefaultTxDecoder(cdc)
	baseApp := baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	dbKeys := NewCyberdAppDbKeys()
	baseApp.SetAppVersion(version.Version)

	var app = &CyberdApp{
		BaseApp:        baseApp,
		cdc:            cdc,
		txDecoder:      txDecoder,
		invCheckPeriod: invCheckPeriod,
		dbKeys:         dbKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams)
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	app.subspaces[slashing.ModuleName] = app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	app.subspaces[mint.ModuleName] = app.paramsKeeper.Subspace(mint.DefaultParamspace)
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)
	app.subspaces[crisis.ModuleName] = app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	app.subspaces[evidence.ModuleName] = app.paramsKeeper.Subspace(evidence.DefaultParamspace)
	app.subspaces[gov.ModuleName] = app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	app.subspaces[bandwidth.ModuleName] = app.paramsKeeper.Subspace(bandwidth.DefaultParamspace)
	app.subspaces[rank.ModuleName] = app.paramsKeeper.Subspace(rank.DefaultParamspace)
	app.subspaces[energy.ModuleName] = app.paramsKeeper.Subspace(energy.DefaultParamspace)
	app.subspaces[cron.ModuleName] = app.paramsKeeper.Subspace(cron.DefaultParamspace)
	app.subspaces[wasm.ModuleName] = app.paramsKeeper.Subspace(wasm.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc, dbKeys.auth, app.subspaces[auth.ModuleName], auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper, app.subspaces[bank.ModuleName], app.ModuleAccountAddrs(),
	)
	app.supplyKeeper = supply.NewKeeper(
		app.cdc, dbKeys.supply, app.accountKeeper, app.bankKeeper, maccPerms,
	)
	stakingKeeper := staking.NewKeeper(
		app.cdc, dbKeys.stake, app.supplyKeeper, app.subspaces[staking.ModuleName],
	)
	app.mintKeeper = mint.NewKeeper(
		app.cdc, dbKeys.mint, app.subspaces[mint.ModuleName], &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName,
	)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.distr, app.subspaces[distr.ModuleName], &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing, &stakingKeeper, app.subspaces[slashing.ModuleName],
	)
	app.crisisKeeper = crisis.NewKeeper(
		app.subspaces[crisis.ModuleName], invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName,
	)
	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, dbKeys.upgrade, app.cdc)

	//app.upgradeKeeper.SetUpgradeHandler("_", func(ctx sdk.Context, plan upgrade.Plan) {})

	evidenceKeeper := evidence.NewKeeper(
		app.cdc, dbKeys.evidence, app.subspaces[evidence.ModuleName], &stakingKeeper, app.slashingKeeper,
	)
	evidenceRouter := evidence.NewRouter()
	evidenceKeeper.SetRouter(evidenceRouter)
	app.evidenceKeeper = evidenceKeeper

	govRouter := gov.NewRouter().AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper)).
		AddRoute(upgrade.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.upgradeKeeper),
	)

	var wasmRouter = baseApp.Router()
	homeDir := viper.GetString(cli.HomeFlag)
	wasmDir := filepath.Join(homeDir, "wasm")

	wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
	err := viper.Unmarshal(&wasmWrap)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	wasmConfig := wasmWrap.Wasm

	supportedFeatures := "staking"
	app.wasmKeeper = wasm.NewKeeper(
		app.cdc,
		dbKeys.wasm,
		app.subspaces[wasm.ModuleName],
		app.accountKeeper,
		app.bankKeeper, // TODO test and put cyberbank
		app.stakingKeeper,
		app.distrKeeper,
		wasmRouter,
		wasmDir,
		wasmConfig,
		supportedFeatures,
		nil,
		nil,
	)

	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.wasmKeeper, enabledProposals))
	}

	app.govKeeper = gov.NewKeeper(
		app.cdc, dbKeys.gov, app.subspaces[gov.ModuleName],
		app.supplyKeeper, &stakingKeeper, govRouter,
	)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)
	app.cyberbank = cyberbank.NewIndexedKeeper(
		cyberbank.NewWrap(&app.bankKeeper, &stakingKeeper, app.supplyKeeper), app.accountKeeper,
	)
	app.energyKeeper = energy.NewKeeper(cdc, dbKeys.energy, app.supplyKeeper, app.cyberbank, app.accountKeeper, app.subspaces[energy.ModuleName])
	app.cyberbank.SetEnergyKeeper(&app.energyKeeper)
	app.cyberbank.Proxy.AddHook(bandwidth.CollectAddressesWithStakeChange())

	app.bandwidthMeter = bandwidth.NewBandwidthMeter(
		dbKeys.bandwidth, app.cyberbank, bandwidth.MsgBandwidthCosts, app.subspaces[bandwidth.ModuleName],
	)
	app.graphKeeper = link.NewKeeper(dbKeys.links)
	app.indexKeeper = link.NewIndexKeeper(app.graphKeeper)
	app.rankKeeper = rank.NewKeeper(
		dbKeys.rank, app.subspaces[rank.ModuleName],
		allowSearch, app.cyberbank, app.indexKeeper, app.graphKeeper, computeUnit,
	)
	app.cronKeeper = cron.NewKeeper(cdc, dbKeys.cron, app.wasmKeeper, app.supplyKeeper, app.cyberbank, app.accountKeeper, app.subspaces[cron.ModuleName])

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.cyberbank.Proxy, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		evidence.NewAppModule(*app.evidenceKeeper),

		bandwidth.NewAppModule(app.accountKeeper, app.bandwidthMeter),
		cyberbank.NewAppModule(app.cyberbank),
		link.NewAppModule(app.graphKeeper, app.indexKeeper, app.accountKeeper),
		rank.NewAppModule(app.rankKeeper),
		energy.NewAppModule(app.energyKeeper),
		cron.NewAppModule(app.cronKeeper),
		wasm.NewAppModule(app.wasmKeeper),
	)

	app.mm.SetOrderBeginBlockers(
		upgrade.ModuleName, mint.ModuleName, distr.ModuleName, slashing.ModuleName,
		evidence.ModuleName, staking.ModuleName,
	)
	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName,
		cyberbank.ModuleName, bandwidth.ModuleName, rank.ModuleName, cron.ModuleName)

	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, gov.ModuleName, mint.ModuleName, supply.ModuleName, crisis.ModuleName,
		bandwidth.ModuleName, genutil.ModuleName, cyberbank.ModuleName, rank.ModuleName,
		evidence.ModuleName, wasm.ModuleName, energy.ModuleName, cron.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.MountStores(
		dbKeys.main, dbKeys.auth, dbKeys.links,
		dbKeys.rank, dbKeys.stake, dbKeys.slashing, dbKeys.gov, dbKeys.params,
		dbKeys.distr, dbKeys.bandwidth, dbKeys.tParams,
		dbKeys.tStake, dbKeys.mint, dbKeys.supply, dbKeys.upgrade, dbKeys.evidence,
		dbKeys.wasm, dbKeys.energy, dbKeys.cron,
	)

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetAnteHandler(
		NewAnteHandler(
			app.accountKeeper, app.bandwidthMeter, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer,
		),
	)

	if loadLatest {
		if err := app.LoadLatestVersion(dbKeys.main); err != nil {
			tmos.Exit(err.Error())
		}
	}

	// ------------------------------

	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.latestBlockHeight = int64(app.rankKeeper.GetLatestBlockNumber(ctx))
	ctx = ctx.WithBlockHeight(app.latestBlockHeight)

	bandwidthParams := bandwidth.DefaultParams()
	rankParams := rank.DefaultParams()
	if app.latestBlockHeight >= 1 {
		bandwidthParams = app.bandwidthMeter.GetParams(ctx)
		rankParams = app.rankKeeper.GetParams(ctx)
	}
	app.bandwidthMeter.SetParams(ctx, bandwidthParams)
	app.rankKeeper.SetParams(ctx, rankParams)

	calculationPeriod := app.rankKeeper.GetParams(ctx).CalculationPeriod
	rankRoundBlockNumber := (app.latestBlockHeight / calculationPeriod) * calculationPeriod
	if rankRoundBlockNumber == 0 && app.latestBlockHeight >= 1 {
		rankRoundBlockNumber = 1
	}
	rankCtx, err := util.NewContextWithMSVersion(db, rankRoundBlockNumber, dbKeys.GetStoreKeys()...)
	if err != nil {
		tmos.Exit(err.Error())
	}

	// ------------------------------

	start := time.Now()
	app.BaseApp.Logger().Info("Loading the in-memory state")
	app.cyberbank.Load(rankCtx, ctx)
	app.indexKeeper.Load(rankCtx, ctx)
	app.bandwidthMeter.Load(ctx)
	app.rankKeeper.Load(ctx, app.Logger())
	app.BaseApp.Logger().Info("Cyber is loaded, Takeoff!", "time", time.Since(start))

	app.Seal()

	return app
}

func (app *CyberdApp) Name() string { return app.BaseApp.Name() }

func (app *CyberdApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *CyberdApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *CyberdApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	resp := app.mm.InitGenesis(ctx, genesisState)
	app.mm.Modules[link.ModuleName].InitGenesis(ctx, nil) // cause link skipped through module manager
	app.mm.Modules[cyberbank.ModuleName].InitGenesis(ctx, nil) // cause cyberrank skipped through module manager
	return resp
}

func (app *CyberdApp) Commit() (res abci.ResponseCommit) {
	app.BaseApp.Commit()
	return abci.ResponseCommit{Data: app.appHash()}
}

func (app *CyberdApp) Info(req abci.RequestInfo) abci.ResponseInfo {

	return abci.ResponseInfo{
		Data:             app.BaseApp.Name(),
		AppVersion:       1,
		LastBlockHeight:  app.LastBlockHeight(),
		LastBlockAppHash: app.appHash(),
	}
}

func (app *CyberdApp) appHash() []byte {

	if app.LastBlockHeight() == 0 {
		return make([]byte, 0)
	}

	return app.rankKeeper.GetNetworkRankHash()
}

func (app *CyberdApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.dbKeys.main)
}

func (app *CyberdApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}