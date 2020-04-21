package types

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"sort"
	"time"

	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/x/link"

	"github.com/tendermint/tendermint/libs/log"
)

type Rank struct {
	Values     []float64
	MerkleTree *merkle.Tree
	CidCount   uint64
	TopCIDs	   []RankedCidNumber
}

func NewRank(values []float64, logger log.Logger, fullTree bool) Rank {
	start := time.Now()
	merkleTree := merkle.NewTree(sha256.New(), fullTree)
	for _, f64 := range values {
		rankBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(rankBytes, math.Float64bits(f64))
		merkleTree.Push(rankBytes)
	}
	logger.Info("Rank constructing tree", "time", time.Since(start))

	// NOTE fulltree true if search index enabled
	var newSortedCIDs []RankedCidNumber
	if (fullTree == true) {
		newSortedCIDs = BuildTop(values, 1000)
	}
	//for cid, rank := range newSortedCIDs {
	//	logger.Info("TOP", "CID", cid, "rank", rank)
	//}

	return Rank{Values: values, MerkleTree: merkleTree, CidCount: uint64(len(values)), TopCIDs: newSortedCIDs}
}

func NewFromMerkle(cidCount uint64, treeBytes []byte) Rank {
	rank := Rank{
		Values:     nil,
		MerkleTree: merkle.NewTree(sha256.New(), false),
		CidCount:   cidCount,
		TopCIDs:    nil,
	}

	rank.MerkleTree.ImportSubtreesRoots(treeBytes)
	return rank
}

func (r Rank) IsEmpty() bool {
	return (r.Values == nil || len(r.Values) == 0) && r.MerkleTree == nil
}

func (r *Rank) Clear() {
	r.Values = nil
	r.MerkleTree = nil
	r.CidCount = 0
	r.TopCIDs = nil
}

func (r *Rank) CopyWithoutTree() Rank {

	if r.Values == nil {
		return Rank{Values: nil, MerkleTree: nil, CidCount: 0, TopCIDs: nil}
	}

	copiedValues := make([]float64, len(r.Values))
	n := copy(copiedValues, r.Values)
	// SHOULD NOT HAPPEN
	if n != len(r.Values) {
		panic("Not all rank values have been copied")
	}
	return Rank{Values: copiedValues, MerkleTree: nil, CidCount: r.CidCount, TopCIDs: r.TopCIDs}
}

// TODO: optimize. Possible solution: adjust capacity of rank values slice after rank calculation
func (r *Rank) AddNewCids(currentCidCount uint64) {

	newCidsCount := currentCidCount - r.CidCount

	// add new cids with rank = 0
	if r.Values != nil {
		r.Values = append(r.Values, make([]float64, newCidsCount)...)
	}

	// extend merkle tree
	zeroRankBytes := make([]byte, 8)
	for i := uint64(0); i < newCidsCount; i++ {
		r.MerkleTree.Push(zeroRankBytes)
	}

	r.CidCount = currentCidCount
}

func BuildTop(values []float64, size int) []RankedCidNumber {
	newSortedCIDs := make(sortableCidNumbers, 0, len(values))
	for cid, rank := range values {
		if (rank == 0) { continue } // NOTE math.IsNaN(rank) check removed after rank fix
		newRankedCid := RankedCidNumber{link.CidNumber(cid), rank}
		newSortedCIDs = append(newSortedCIDs, newRankedCid)
	}
	sort.Stable(sort.Reverse(newSortedCIDs))
	if (len(values) > size) {
		newSortedCIDs = newSortedCIDs[0:(size-1)]
	}
	return newSortedCIDs
}
