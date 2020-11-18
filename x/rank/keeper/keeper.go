package keeper

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"strconv"

	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"time"
)

type StateKeeper struct {

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
	//mainKeeper        store.MainKeeper
	//stakeKeeper       types.StakeKeeper
	//cidNumKeeper      types.CidNumberKeeper
	//linkIndexedKeeper types.LinkIndexedKeeper
	stakeKeeper        types.StakeKeeper
	graphKeeper        types.GraphKeeper
	graphIndexedKeeper types.GraphIndexedKeeper

	storeKey          sdk.StoreKey
	paramSpace        params.Subspace

	index         types.SearchIndex
	getIndexError types.GetError
}

func NewKeeper(
	key sdk.StoreKey, paramSpace params.Subspace, allowSearch bool, stakeIndex types.StakeKeeper,
	graphIndexedKeeper types.GraphIndexedKeeper, graphKeeper types.GraphKeeper,
	unit types.ComputeUnit,
) *StateKeeper {
	return &StateKeeper{
		storeKey:       key,
		paramSpace: 	paramSpace.WithKeyTable(types.ParamKeyTable()),
		allowSearch:    allowSearch,
		rankCalcChan:   make(chan types.Rank, 1),
		rankErrChan:    make(chan error),
		rankCalculationFinished: true,
		stakeKeeper:    stakeIndex,
		graphIndexedKeeper: graphIndexedKeeper,
		graphKeeper:   graphKeeper,
		computeUnit:    unit,
		hasNewLinksForPeriod: true,
	}
}

func (s *StateKeeper) Load(ctx sdk.Context, log log.Logger) {
	s.networkCidRank = types.NewFromMerkle(s.graphKeeper.GetCidsCount(ctx), s.GetLatestMerkleTree(ctx))
	s.nextCidRank = types.NewFromMerkle(s.GetNextRankCidCount(ctx), s.GetNextMerkleTree(ctx))
	s.cidCount = int64(s.graphKeeper.GetCidsCount(ctx))

	s.index = s.BuildSearchIndex(log)
	s.index.Load(s.graphIndexedKeeper.GetOutLinks()) // was wrong with GetNextOutLinks
	s.getIndexError = s.index.Run()

	params := s.GetParams(ctx)

	dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
	if err != nil {
		panic(err)
	}

	tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
	if err != nil {
		panic(err)
	}

	// if we fell down and need to start new rank calculation
	if !s.GetRankCalculationFinished(ctx) {
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
	s.StoreLatestBlockNumber(ctx, uint64(ctx.BlockHeight()))
	currentCidsCount := s.graphKeeper.GetCidsCount(ctx)

	s.index.PutNewLinks(s.graphIndexedKeeper.GetCurrentBlockNewLinks())

	blockHasNewLinks := s.graphIndexedKeeper.EndBlocker()
	s.hasNewLinksForPeriod = s.hasNewLinksForPeriod || blockHasNewLinks

	params := s.GetParams(ctx)

	if ctx.BlockHeight()%params.CalculationPeriod == 0 || ctx.BlockHeight() == 1 {

		dampingFactor, err := strconv.ParseFloat(params.DampingFactor.String(), 64)
		if err != nil {
			panic(err)
		}

		tolerance, err := strconv.ParseFloat(params.Tolerance.String(), 64)
		if err != nil {
			panic(err)
		}

		s.checkRankCalcFinished(ctx, true, log)
		s.applyNextRank()

		s.cidCount = int64(currentCidsCount)
		stakeChanged := s.stakeKeeper.DetectUsersStakeChange(ctx)

		// start new calculation
		if s.hasNewLinksForPeriod || stakeChanged {
			s.graphIndexedKeeper.FixLinks()
			s.rankCalculationFinished = false
			s.hasNewLinksForPeriod = false
			s.StoreRankCalculationFinished(ctx, false)
			s.startRankCalculation(ctx, dampingFactor, tolerance, log)
		}
	} else {
		s.checkRankCalcFinished(ctx, false, log)
	}
	s.networkCidRank.AddNewCids(currentCidsCount)
	s.StoreLatestMerkleTree(ctx, s.getNetworkMerkleTreeAsBytes())
}

func (s *StateKeeper) Search(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error) {
	return s.index.Search(cidNumber, page, perPage)
}

func (s *StateKeeper) Backlinks(cidNumber link.CidNumber, page, perPage int) ([]types.RankedCidNumber, int, error) {
	return s.index.Backlinks(cidNumber, page, perPage)
}

func (s *StateKeeper) Top(page, perPage int) ([]types.RankedCidNumber, int, error) {
	return s.index.Top(page, perPage)
}

func (s *StateKeeper) Accounts(account uint64, page, perPage int) (map[link.CidNumber]link.CidNumber, int, error) {
	return s.index.Accounts(account, page, perPage)
}

func (s *StateKeeper) GetRankValue(cidNumber link.CidNumber) uint64 {
	return s.index.GetRankValue(cidNumber)
}

func (s *StateKeeper) startRankCalculation(ctx sdk.Context, dampingFactor float64, tolerance float64, log log.Logger) {

	calcCtx := types.NewCalcContext(ctx, s.graphIndexedKeeper, s.graphKeeper, s.stakeKeeper, s.allowSearch, dampingFactor, tolerance)
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
	s.StoreNextMerkleTree(ctx, s.getNextMerkleTreeAsBytes())
	s.StoreNextRankCidCount(ctx, newRank.CidCount)
	s.rankCalculationFinished = true
	s.StoreRankCalculationFinished(ctx, true)
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

func (s *StateKeeper) GetRankValues(number uint64) uint64 {
	return s.networkCidRank.Values[number]
}

func (s *StateKeeper) GetNetworkRankHash() []byte {
	fmt.Println("GetNetworkRankHash: ", hex.EncodeToString(s.networkCidRank.MerkleTree.RootHash()))
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

func (s *StateKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	s.paramSpace.GetParamSet(ctx, &params)
	return params
}

func (s *StateKeeper) SetParams(ctx sdk.Context, params types.Params) {
	s.paramSpace.SetParamSet(ctx, &params)
}

func (k *StateKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}


func (sk StateKeeper) GetLatestBlockNumber(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.LatestBlockNumber)
	if numberAsBytes == nil {
		return 0
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreLatestBlockNumber(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsBytes, number)
	store.Set(types.LatestBlockNumber, numberAsBytes)
}

func (sk StateKeeper) GetLatestMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(sk.storeKey)
	return store.Get(types.LatestMerkleTree)
}

func (sk StateKeeper) StoreLatestMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(sk.storeKey)
	store.Set(types.LatestMerkleTree, treeAsBytes)
}

func (sk StateKeeper) GetNextMerkleTree(ctx sdk.Context) []byte {
	store := ctx.KVStore(sk.storeKey)
	return store.Get(types.NextMerkleTree)
}

func (sk StateKeeper) StoreNextMerkleTree(ctx sdk.Context, treeAsBytes []byte) {
	store := ctx.KVStore(sk.storeKey)
	store.Set(types.NextMerkleTree, treeAsBytes)
}

func (sk StateKeeper) StoreRankCalculationFinished(ctx sdk.Context, finished bool) {
	store := ctx.KVStore(sk.storeKey)
	var byteFlag byte
	if finished {
		byteFlag = 1
	}
	store.Set(types.RankCalcFinished, []byte{byteFlag})
}

func (sk StateKeeper) GetRankCalculationFinished(ctx sdk.Context) bool {
	store := ctx.KVStore(sk.storeKey)
	bytes := store.Get(types.RankCalcFinished)
	if bytes == nil || bytes[0] == 1 {
		return true
	}
	return false
}

func (sk StateKeeper) GetNextRankCidCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := store.Get(types.NextRankCidCount)
	if numberAsBytes == nil {
		return sk.graphKeeper.GetCidsCount(ctx)
	}
	return binary.LittleEndian.Uint64(numberAsBytes)
}

func (sk StateKeeper) StoreNextRankCidCount(ctx sdk.Context, number uint64) {
	store := ctx.KVStore(sk.storeKey)
	numberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(numberAsBytes, number)
	store.Set(types.NextRankCidCount, numberAsBytes)
}