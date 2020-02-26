package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/x/bank"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"time"
)

type StateKeeper struct {
	*BaseRankKeeper
	cdc *codec.Codec

	networkCidRank types.Rank // array linksIndex is cid number
	nextCidRank    types.Rank // array linksIndex is cid number

	rankCalculationFinished bool
	cidCount                int64

	hasNewLinksForPeriod bool // indicates if new links where submitted for rank calc period

	rankCalcChan chan types.Rank
	rankErrChan  chan error
	allowSearch  bool
	computeUnit  types.ComputeUnit

	// keepers
	mainKeeper        store.MainKeeper
	stakeKeeper       types.StakeKeeper
	cidNumKeeper      types.CidNumberKeeper
	linkIndexedKeeper types.LinkIndexedKeeper
	paramsKeeper      params.Keeper

	// index
	index         types.SearchIndex
	getIndexError types.GetError
}

func NewStateKeeper(
	cdc *codec.Codec, paramSpace *params.Subspace, allowSearch bool, mainKeeper store.MainKeeper, stakeIndex bank.IndexedKeeper,
	linkIndexedKeeper types.LinkIndexedKeeper, cidNumKeeper types.CidNumberKeeper,
	pk params.Keeper, unit types.ComputeUnit,
) *StateKeeper {
	return &StateKeeper{
		cdc:            cdc, //
		BaseRankKeeper: NewBaseRankKeeper(paramSpace),
		allowSearch:    allowSearch,
		rankCalcChan:   make(chan types.Rank, 1),
		rankErrChan:    make(chan error),
		rankCalculationFinished: true,
		mainKeeper:     mainKeeper,
		stakeKeeper:    stakeIndex,
		linkIndexedKeeper: linkIndexedKeeper,
		cidNumKeeper:   cidNumKeeper,
		paramsKeeper:   pk,
		computeUnit:    unit,
		hasNewLinksForPeriod: true,
	}
}

func (s *StateKeeper) Load(ctx sdk.Context, log log.Logger) {
	s.networkCidRank = types.NewFromMerkle(s.mainKeeper.GetCidsCount(ctx), s.mainKeeper.GetLatestMerkleTree(ctx))
	s.nextCidRank = types.NewFromMerkle(s.mainKeeper.GetNextRankCidCount(ctx), s.mainKeeper.GetNextMerkleTree(ctx))
	s.cidCount = int64(s.mainKeeper.GetCidsCount(ctx))

	s.index = s.BuildSearchIndex(log)
	s.index.Load(s.linkIndexedKeeper.GetOutLinks())
	s.getIndexError = s.index.Run()

	paramset := s.GetParamSet(ctx)

	dampingFactor, err := strconv.ParseFloat(paramset.DampingFactor.String(), 64)
	if err != nil {
		panic(err)
	}

	tolerance, err := strconv.ParseFloat(paramset.Tolerance.String(), 64)
	if err != nil {
		panic(err)
	}

	// if we fell down and need to start new rank calculation
	if !s.mainKeeper.GetRankCalculationFinished(ctx) {
		s.startRankCalculation(ctx, dampingFactor, tolerance, log)
		s.rankCalculationFinished = false
	}
}

func (s *StateKeeper) BuildSearchIndex(logger log.Logger) types.SearchIndex {
	if s.allowSearch {
		return types.NewBaseSearchIndex(logger)
	} else {
		return types.NoopSearchIndex{}
	}
}

func (s *StateKeeper) EndBlocker(ctx sdk.Context, log log.Logger) {
	currentCidsCount := s.mainKeeper.GetCidsCount(ctx)

	s.index.PutNewLinks(s.linkIndexedKeeper.GetCurrentBlockNewLinks())

	blockHasNewLinks := s.linkIndexedKeeper.EndBlocker()
	s.hasNewLinksForPeriod = s.hasNewLinksForPeriod || blockHasNewLinks

	paramset := s.GetParamSet(ctx)

	if ctx.BlockHeight()%paramset.CalculationPeriod == 0 || ctx.BlockHeight() == 1 {

		dampingFactor, err := strconv.ParseFloat(paramset.DampingFactor.String(), 64)
		if err != nil {
			panic(err)
		}

		tolerance, err := strconv.ParseFloat(paramset.Tolerance.String(), 64)
		if err != nil {
			panic(err)
		}

		s.checkRankCalcFinished(ctx, true, log)
		s.applyNextRank()

		s.cidCount = int64(currentCidsCount)
		stakeChanged := s.stakeKeeper.FixUserStake(ctx)

		// start new calculation
		if s.hasNewLinksForPeriod || stakeChanged {
			s.linkIndexedKeeper.FixLinks()
			s.rankCalculationFinished = false
			s.hasNewLinksForPeriod = false
			s.mainKeeper.StoreRankCalculationFinished(ctx, false)
			s.startRankCalculation(ctx, dampingFactor, tolerance, log)
		}
	} else {
		s.checkRankCalcFinished(ctx, false, log)
	}
	s.networkCidRank.AddNewCids(currentCidsCount)
	s.mainKeeper.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())
}

func (s *StateKeeper) Search(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error) {
	return s.index.Search(cidNumber, page, perPage)
}

func (s *StateKeeper) Top(page, perPage int) ([]types.RankedCidNumber, int, error) {
	return s.index.Top(page, perPage)
}

func (s *StateKeeper) GetRankValue(cidNumber link.CidNumber) float64 {
	return s.index.GetRankValue(cidNumber)
}

func (s *StateKeeper) startRankCalculation(ctx sdk.Context, dampingFactor float64, tolerance float64, log log.Logger) {


	calcCtx := types.NewCalcContext(ctx, s.linkIndexedKeeper, s.cidNumKeeper, s.stakeKeeper, s.allowSearch, dampingFactor, tolerance)
	go CalculateRankInParallel(calcCtx, s.rankCalcChan, s.rankErrChan, s.computeUnit, log)
}

func (s *StateKeeper) checkRankCalcFinished(ctx sdk.Context, block bool, log log.Logger) {

	if !s.rankCalculationFinished {
		for {
			select {
			case newRank := <-s.rankCalcChan:
				s.handleNextRank(ctx, newRank)
				return
			case err := <-s.rankErrChan:
				// DUMB ERROR HANDLING
				log.Error("Error during rank calculation " + err.Error())
				panic(err.Error())
			default:
				if !block {
					return
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *StateKeeper) handleNextRank(ctx sdk.Context, newRank types.Rank) {
	s.nextCidRank = newRank
	s.mainKeeper.StoreNextMerkleTree(ctx, s.getNextMerkleTreeAsBytes())
	s.mainKeeper.StoreNextRankCidCount(ctx, newRank.CidCount)
	s.rankCalculationFinished = true
	s.mainKeeper.StoreRankCalculationFinished(ctx, true)
}

func (s *StateKeeper) applyNextRank() {

	if !s.nextCidRank.IsEmpty() {
		s.networkCidRank = s.nextCidRank
		s.index.PutNewRank(s.networkCidRank)
	}
	s.nextCidRank.Clear()
}

//
// GETTERS
//

func (s *StateKeeper) GetNetworkRankHash() []byte {
	return s.networkCidRank.MerkleTree.RootHash()
}

func (s *StateKeeper) getNetworkMerkleTreeAsBytes() []byte {
	return s.networkCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *StateKeeper) getNextMerkleTreeAsBytes() []byte {
	return s.nextCidRank.MerkleTree.ExportSubtreesRoots()
}

func (s *StateKeeper) GetLastCidNum() link.CidNumber {
	return link.CidNumber(len(s.networkCidRank.Values) - 1)
}

func (s *StateKeeper) GetMerkleTree() *merkle.Tree {
	return s.networkCidRank.MerkleTree
}

func (s *StateKeeper) GetIndexError() error {
	return s.getIndexError()
}

func (s *StateKeeper) GetParamSet(ctx sdk.Context) (params types.Params) {
	subspace, ok := s.paramsKeeper.GetSubspace(DefaultParamspace)
	if !ok {
		panic("rank params subspace is not found")
	}
	subspace.GetParamSet(ctx, &params)
	return params
}
