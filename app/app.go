package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/abci/version"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/types"
	cbd "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/types/coin"
	"github.com/cybercongress/go-cyber/util"
	bw "github.com/cybercongress/go-cyber/x/bandwidth"
	cbdbank "github.com/cybercongress/go-cyber/x/bank"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank"
)

const (
	APP            = "cyberd"
	appName        = "CyberdApp"
	DefaultKeyPass = "12345678" //todo remove from here
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cyberdcli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cyberd")

	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		sdkbank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},

		link.AppModuleBasic{},
		bw.AppModuleBasic{},
		rank.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

// CyberdApp implements an extended ABCI application. It contains a BaseApp,
// a codec for serialization, KVStore dbKeys for multistore state management, and
// various mappers and keepers to manage getting, setting, and serializing the
// integral app types.
type CyberdApp struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	txDecoder      sdk.TxDecoder
	invCheckPeriod uint

	// bandwidth
	bandwidthMeter       bw.Meter
	blockBandwidthKeeper bw.BlockSpentBandwidthKeeper

	// keys to access the multistore
	dbKeys cyberdAppDbKeys

	// keepers: manage getting and setting app data
	mainKeeper         store.MainKeeper
	accountKeeper      auth.AccountKeeper
	bankKeeper         cbdbank.Keeper
	supplyKeeper       supply.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	mintKeeper         mint.Keeper
	distrKeeper        distr.Keeper
	govKeeper          gov.Keeper
	paramsKeeper       params.Keeper
	crisisKeeper       crisis.Keeper
	accBandwidthKeeper bw.Keeper

	// cyberd storage
	// todo: move all processes with this storages to another file
	linkIndexedKeeper  link.IndexedKeeper
	cidNumKeeper       link.CidNumberKeeper
	stakingIndexKeeper cbdbank.IndexedKeeper
	rankStateKeeper    rank.StateKeeper

	latestBlockHeight int64

	mm *module.Manager
}

// NewCyberdApp returns a reference to a new CyberdApp given a
// logger and database. Internally, a codec is created along with all the necessary dbKeys.
// In addition, all necessary mappers and keepers are created, routes
// registered, and finally the stores being mounted along with any necessary
// chain initialization.
func NewCyberdApp(logger log.Logger, db dbm.DB, traceStore io.Writer, height int64,
	invCheckPeriod uint, computeUnit rank.ComputeUnit, allowSearch bool, baseAppOptions ...func(*baseapp.BaseApp)) *CyberdApp {
	// create and register app-level codec for TXs and accounts
	cdc := MakeCodec()
	txDecoder := auth.DefaultTxDecoder(cdc)
	baseApp := baseapp.NewBaseApp(appName, logger, db, txDecoder, baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	dbKeys := NewCyberdAppDbKeys()
	mainKeeper := store.NewMainKeeper(dbKeys.main)
	baseApp.SetAppVersion(version.Version)

	// create your application type
	var app = &CyberdApp{
		cdc:            cdc,
		txDecoder:      txDecoder,
		invCheckPeriod: invCheckPeriod,
		BaseApp:        baseApp,
		dbKeys:         dbKeys,
		mainKeeper:     mainKeeper,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, dbKeys.params, dbKeys.tParams, params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(sdkbank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)
	bandwidthSubspace := app.paramsKeeper.Subspace(bw.DefaultParamspace).WithKeyTable(bw.ParamKeyTable())
	rankSubspace := app.paramsKeeper.Subspace(rank.DefaultParamspace).WithKeyTable(rank.ParamKeyTable())

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, dbKeys.acc, authSubspace, cbd.NewCyberdAccount)
	bankKeeper := cbdbank.NewKeeper(app.accountKeeper, bankSubspace, sdkbank.DefaultCodespace, app.ModuleAccountAddrs())
	app.bankKeeper = bankKeeper
	app.supplyKeeper = supply.NewKeeper(app.cdc, dbKeys.supply, app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, dbKeys.stake, dbKeys.tStake, app.supplyKeeper, stakingSubspace, staking.DefaultCodespace)
	app.bankKeeper.SetStakingKeeper(&stakingKeeper)
	app.bankKeeper.SetSupplyKeeper(app.supplyKeeper)
	app.mintKeeper = mint.NewKeeper(
		app.cdc, dbKeys.mint, mintSubspace, &stakingKeeper, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(
		app.cdc, dbKeys.distr, distrSubspace, &stakingKeeper, app.supplyKeeper,
		distr.DefaultCodespace, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, dbKeys.slashing, &stakingKeeper, slashingSubspace, slashing.DefaultCodespace)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	app.accBandwidthKeeper = bw.NewAccBandwidthKeeper(dbKeys.accBandwidth, &bandwidthSubspace)
	app.blockBandwidthKeeper = bw.NewBlockSpentBandwidthKeeper(dbKeys.blockBandwidth)

	// register the proposal types
	govRouter := gov.NewRouter().
		AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

	app.govKeeper = gov.NewKeeper(
		app.cdc, dbKeys.gov, app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
	)

	// cyberd keepers
	app.linkIndexedKeeper = link.NewIndexedKeeper(link.NewLinkKeeper(mainKeeper, dbKeys.links))
	app.cidNumKeeper = link.NewCidNumberKeeper(mainKeeper, dbKeys.cidNum, dbKeys.cidNumReverse)
	app.stakingIndexKeeper = cbdbank.NewIndexedKeeper(bankKeeper)
	app.rankStateKeeper = rank.NewStateKeeper(app.cdc, &rankSubspace,
		allowSearch, app.mainKeeper, app.stakingIndexKeeper,
		app.linkIndexedKeeper, app.cidNumKeeper, app.paramsKeeper, computeUnit,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above are passed by reference,
	// so that it can be modified like below:
	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakeHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)
	app.bankKeeper.AddHook(bw.CollectAddressesWithStakeChange())

	app.bandwidthMeter = bw.NewBaseMeter(
		app.mainKeeper, app.paramsKeeper, app.accountKeeper, app.accBandwidthKeeper,
		app.blockBandwidthKeeper, app.bankKeeper, bw.MsgBandwidthCosts,
	)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		auth.NewAppModule(app.accountKeeper),
		sdkbank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),

		bw.NewAppModule(app.accBandwidthKeeper, app.blockBandwidthKeeper),
		link.NewAppModule(app.cidNumKeeper, app.linkIndexedKeeper, app.accountKeeper),
		rank.NewAppModule(app.rankStateKeeper),
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// mount the multistore and load the latest state
	app.MountStores(
		dbKeys.main, dbKeys.acc, dbKeys.cidNum, dbKeys.cidNumReverse, dbKeys.links, dbKeys.rank, dbKeys.stake,
		dbKeys.slashing, dbKeys.gov, dbKeys.params, dbKeys.distr, dbKeys.fees, dbKeys.accBandwidth,
		dbKeys.blockBandwidth, dbKeys.tParams, dbKeys.tStake, dbKeys.mint, dbKeys.supply,
	)

	app.SetInitChainer(app.applyGenesis)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	//because genesis max_gas equals -1 there is NewInfiniteGasMeter
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))

	if height == -1 {
		err := app.LoadLatestVersion(dbKeys.main)
		if err != nil {
			cmn.Exit(err.Error())
		}
	} else {
		err := app.LoadVersion(height, dbKeys.main)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}
	ctx := app.BaseApp.NewContext(true, abci.Header{})
	app.latestBlockHeight = int64(mainKeeper.GetLatestBlockNumber(ctx))
	ctx = ctx.WithBlockHeight(app.latestBlockHeight)

	bandwidthParamset := bw.NewDefaultParams()
	bandwidthSubspace.SetParamSet(ctx, &bandwidthParamset)

	rankParamset := rank.NewDefaultParams()
	rankSubspace.SetParamSet(ctx, &rankParamset)

	// build context for current rank calculation round
	rankRoundBlockNumber := (app.latestBlockHeight / rankParamset.CalculationPeriod) * rankParamset.CalculationPeriod
	if rankRoundBlockNumber == 0 && app.latestBlockHeight >= 1 {
		rankRoundBlockNumber = 1 // special case cause tendermint blocks start from 1
	}
	rankCtx, err := util.NewContextWithMSVersion(db, rankRoundBlockNumber, dbKeys.GetStoreKeys()...)
	if err != nil {
		cmn.Exit(err.Error())
	}

	// load in-memory data
	start := time.Now()
	app.BaseApp.Logger().Info("Loading mem state")
	app.linkIndexedKeeper.Load(rankCtx, ctx)
	app.stakingIndexKeeper.Load(rankCtx, ctx)
	app.BaseApp.Logger().Info("App loaded", "time", time.Since(start))

	// BANDWIDTH LOAD
	app.bandwidthMeter.Load(ctx)

	// RANK PARAMS
	app.rankStateKeeper.Load(ctx, app.Logger())
	app.Seal()
	return app
}

func (app *CyberdApp) applyGenesis(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	start := time.Now()
	app.Logger().Info("Applying genesis")
	var genesisState GenesisState
	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	// load the accounts
	for _, gacc := range genesisState.Accounts {
		acc := gacc.ToAccount()
		app.accountKeeper.GetNextAccountNumber(ctx) //increment for future accs being started from right number
		app.accountKeeper.SetAccount(ctx, acc)
		app.stakingIndexKeeper.UpdateStake(
			types.AccNumber(acc.GetAccountNumber()),
			acc.GetCoins().AmountOf(coin.CYB).Int64())
	}

	cbdbank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, app.supplyKeeper, genesisState.DistrData)
	supply.InitGenesis(ctx, app.supplyKeeper, app.accountKeeper, genesisState.SupplyData)

	validators := staking.InitGenesis(
		ctx, app.stakingKeeper, app.accountKeeper, app.supplyKeeper, genesisState.StakingData)

	auth.InitGenesis(ctx, app.accountKeeper, genesisState.AuthData)
	slashing.InitGenesis(ctx, app.slashingKeeper, app.stakingKeeper, genesisState.SlashingData)
	gov.InitGenesis(ctx, app.govKeeper, app.supplyKeeper, genesisState.GovData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)
	bw.InitGenesis(ctx, app.bandwidthMeter, app.accBandwidthKeeper, genesisState.GetAddresses(),
		genesisState.BandwidthData)
	rank.InitGenesis(ctx, app.rankStateKeeper, genesisState.RankData)

	err = link.InitGenesis(ctx, app.cidNumKeeper, app.linkIndexedKeeper, app.Logger())
	if err != nil {
		panic(err)
	}

	crisis.InitGenesis(ctx, app.crisisKeeper, genesisState.Crisis)

	err = validateGenesisState(genesisState)
	if err != nil {
		panic(err)
	}

	// deliver genesis txes
	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := abci.RequestDeliverTx{Tx: app.cdc.MustMarshalBinaryLengthPrefixed(tx)}
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	// sanity check
	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(abci.ValidatorUpdates(req.Validators))
		sort.Sort(abci.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	app.Logger().Info("Genesis applied", "time", time.Since(start))
	return abci.ResponseInitChain{
		Validators: validators,
	}
}

func (app *CyberdApp) CheckTx(req abci.RequestCheckTx) (res abci.ResponseCheckTx) {

	ctx := app.NewContext(true, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, req.GetTx())

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(ctx, tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		curBlockSpentBandwidth := app.bandwidthMeter.GetCurBlockSpentBandwidth(ctx)
		maxBlockBandwidth := app.bandwidthMeter.GetMaxBlockBandwidth(ctx)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else if (uint64(txCost) + curBlockSpentBandwidth) > maxBlockBandwidth  {
			err = types.ErrExceededMaxBlockBandwidth()
		} else {
			resp := app.BaseApp.CheckTx(req)
			if resp.Code == 0 {
				app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)
			}
			return resp
		}
	}

	result := err.Result()
	return abci.ResponseCheckTx{
		Code:      uint32(result.Code),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *CyberdApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {

	ctx := app.NewContext(false, abci.Header{Height: app.latestBlockHeight})
	tx, acc, err := app.decodeTxAndAcc(ctx, req.GetTx())

	if err == nil {

		txCost := app.bandwidthMeter.GetPricedTxCost(ctx, tx)
		accBw := app.bandwidthMeter.GetCurrentAccBandwidth(ctx, acc)

		curBlockSpentBandwidth := app.bandwidthMeter.GetCurBlockSpentBandwidth(ctx)
		maxBlockBandwidth := app.bandwidthMeter.GetMaxBlockBandwidth(ctx)

		if !accBw.HasEnoughRemained(txCost) {
			err = types.ErrNotEnoughBandwidth()
		} else if (uint64(txCost) + curBlockSpentBandwidth) > maxBlockBandwidth  {
			err = types.ErrExceededMaxBlockBandwidth()
		} else {
			resp := app.BaseApp.DeliverTx(req)
			app.bandwidthMeter.ConsumeAccBandwidth(ctx, accBw, txCost)

			linkingCost := app.bandwidthMeter.GetPricedLinksCost(ctx, tx)
			if linkingCost != int64(0) {
				app.bandwidthMeter.UpdateLinkedBandwidth(ctx, accBw, linkingCost)
			}

			app.bandwidthMeter.AddToBlockBandwidth(app.bandwidthMeter.GetTxCost(ctx, tx))
			app.bandwidthMeter.AddToOverallKarma(linkingCost)

			return abci.ResponseDeliverTx{
				Code:      resp.GetCode(),
				Codespace: resp.GetCodespace(),
				Data:      resp.GetData(),
				Log:       resp.GetLog(),
				GasWanted: resp.GetGasWanted(),
				GasUsed:   resp.GetGasUsed(),
				Events:    append(resp.GetEvents(), getSignersEvents(tx).ToABCIEvents()...),
			}
		}
	}

	result := err.Result()
	return abci.ResponseDeliverTx{
		Code:      uint32(result.Code),
		Codespace: string(result.Codespace),
		Data:      result.Data,
		Log:       result.Log,
		GasWanted: int64(result.GasWanted),
		GasUsed:   int64(result.GasUsed),
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *CyberdApp) decodeTxAndAcc(ctx sdk.Context, txBytes []byte) (auth.StdTx, sdk.AccAddress, sdk.Error) {

	decoded, err := app.txDecoder(txBytes)
	if err != nil {
		return auth.StdTx{}, nil, err
	}

	tx := decoded.(auth.StdTx)
	if tx.GetMsgs() == nil || len(tx.GetMsgs()) == 0 {
		return tx, nil, sdk.ErrInternal("Tx.GetMsgs() must return at least one message in list")
	}

	if err := tx.ValidateBasic(); err != nil {
		return tx, nil, err
	}

	// signers acc [0] bandwidth will be consumed
	account := tx.GetSigners()[0]
	acc := app.accountKeeper.GetAccount(ctx, account)
	if acc == nil {
		return tx, nil, sdk.ErrUnknownAddress(account.String())
	}

	return tx, account, nil
}

func getSignersEvents(tx sdk.Tx) sdk.Events {
	events := sdk.EmptyEvents()

	for _, msg := range tx.GetMsgs() {
		for _, accAddress := range msg.GetSigners() {
			attribute := sdk.NewAttribute("address", accAddress.String())
			event := sdk.NewEvent("signer", attribute)
			events.AppendEvent(event)
		}
	}

	return events
}

func (app *CyberdApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {

	// mint new tokens for the previous block
	mint.BeginBlocker(ctx, app.mintKeeper)
	// distribute rewards for the previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// slash anyone who double signed.
	// NOTE: This should happen after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool,
	// so as to keep the CanWithdrawInvariant invariant.
	slashing.BeginBlocker(ctx, req, app.slashingKeeper)
	return abci.ResponseBeginBlock{Events: ctx.EventManager().ABCIEvents()}
}

// Calculates cyber.Rank for block N, and returns Hash of result as app state.
// Calculated app state will be included in N+1 block header, thus influence on block hash.
// App state is consensus driven state.
func (app *CyberdApp) EndBlocker(ctx sdk.Context, _ abci.RequestEndBlock) abci.ResponseEndBlock {

	app.latestBlockHeight = ctx.BlockHeight()
	app.mainKeeper.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))
	crisis.EndBlocker(ctx, app.crisisKeeper)
	gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates := staking.EndBlocker(ctx, app.stakingKeeper)
	app.stakingIndexKeeper.EndBlocker(ctx)

	bw.EndBlocker(ctx, app.paramsKeeper, app.bandwidthMeter)

	// RANK CALCULATION
	app.rankStateKeeper.EndBlocker(ctx, app.Logger())

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Events:           ctx.EventManager().ABCIEvents(),
	}
}

// Implements ABCI
func (app *CyberdApp) Commit() (res abci.ResponseCommit) {
	err := app.linkIndexedKeeper.Commit(uint64(app.latestBlockHeight))
	if err != nil {
		panic(err)
	}
	app.BaseApp.Commit()
	return abci.ResponseCommit{Data: app.appHash()}
}

// Implements ABCI
func (app *CyberdApp) Info(req abci.RequestInfo) abci.ResponseInfo {

	return abci.ResponseInfo{
		Data:             app.BaseApp.Name(),
		LastBlockHeight:  app.LastBlockHeight(),
		LastBlockAppHash: app.appHash(),
	}
}

func (app *CyberdApp) appHash() []byte {

	if app.LastBlockHeight() == 0 {
		return make([]byte, 0)
	}

	linkHash := app.linkIndexedKeeper.GetNetworkLinkHash()
	rankHash := app.rankStateKeeper.GetNetworkRankHash()

	result := make([]byte, len(linkHash))

	for i := 0; i < len(linkHash); i++ {
		result[i] = linkHash[i] ^ rankHash[i]
	}

	return result
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *CyberdApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}