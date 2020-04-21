package keeper

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math"

	"github.com/tendermint/tendermint/libs/log"

	cbd "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)


func calculateRankCPU(ctx *types.CalculationContext, logger log.Logger) []float64 {

	inLinks := ctx.GetInLinks()
	tolerance := ctx.GetTolerance()
	dampingFactor := ctx.GetDampingFactor()

	size := ctx.GetCidsCount()
	if size == 0 {
		return []float64{}
	}

	rank := make([]float64, size)
	defaultRank := (1.0 - dampingFactor) / float64(size)
	danglingNodesSize := uint64(0)

	for i := range rank {
		rank[i] = defaultRank
		if len(inLinks[link.CidNumber(i)]) == 0 {
			danglingNodesSize++
		}
	}

	logger.Info("RANK", "defaultRank:", defaultRank)
	logger.Info("RANK", "danglingNodesSize:", danglingNodesSize)

	innerProductOverSize := defaultRank * (float64(danglingNodesSize) / float64(size))
	defaultRankWithCorrection := float64(dampingFactor*innerProductOverSize) + defaultRank

	logger.Info("RANK", "defaultRankWithCorrection:", defaultRankWithCorrection)

	change := tolerance + 1

	steps := 0
	prevrank := make([]float64, 0)
	prevrank = append(prevrank, rank...)
	for change > tolerance {
		rank = step(ctx, defaultRankWithCorrection, dampingFactor, prevrank)
		change = calculateChange(prevrank, rank)
		prevrank = rank
		steps++
		logger.Info("RANK", "step:", steps)
	}


	st := ctx.GetStakes()
	outl := ctx.GetOutLinks()
	inl := ctx.GetInLinks()
	saveStakesToBytesFile(&st, "./diffX/stakes.data")
	saveLinksToBytesFile(&outl, "./diffX/outLinks.data")
	saveLinksToBytesFile(&inl, "./diffX/inLinks.data")
	logger.Info("RANK: data saved", )

	return rank
}

func step(ctx *types.CalculationContext, defaultRankWithCorrection float64, dampingFactor float64, prevrank []float64) []float64 {

	rank := append(make([]float64, 0, len(prevrank)), prevrank...)

	for cid := range ctx.GetInLinks() {
		_, sortedCids, ok := ctx.GetSortedInLinks(cid)

		if !ok {
			continue
		} else {
			ksum := float64(0)
			for _, j := range sortedCids {
				linkStake := getOverallLinkStake(ctx, j, cid)
				jCidOutStake := getOverallOutLinksStake(ctx, j)
				weight := float64(linkStake) / float64(jCidOutStake)
				if math.IsNaN(weight) { weight = float64(0) }
				ksum = prevrank[j]*weight + ksum //force no-fma here by explicit conversion
			}
			rank[cid] = ksum*dampingFactor + defaultRankWithCorrection //force no-fma here by explicit conversion
		}
	}

	return rank
}

func getOverallLinkStake(ctx *types.CalculationContext, from link.CidNumber, to link.CidNumber) uint64 {

	stake := uint64(0)
	users := ctx.GetOutLinks()[from][to]
	for user := range users {
		stake += ctx.GetStakes()[user]
	}
	return stake
}

func getOverallOutLinksStake(ctx *types.CalculationContext, from link.CidNumber) uint64 {

	stake := uint64(0)
	for to := range ctx.GetOutLinks()[from] {
		stake += getOverallLinkStake(ctx, from, to)
	}
	return stake
}

func calculateChange(prevrank, rank []float64) float64 {

	maxDiff := 0.0
	diff := 0.0
	for i, pForI := range prevrank {
		if pForI > rank[i] {
			diff = pForI - rank[i]
		} else {
			diff = rank[i] - pForI
		}
		if diff > maxDiff {
			maxDiff = diff
		}
	}

	return maxDiff
}

func saveLinksToBytesFile(links *map[link.CidNumber]link.CidLinks, fileName string) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(links)
	if err != nil {
		fmt.Printf("encode error:", err)
	}
	err = ioutil.WriteFile(fileName, network.Bytes(), 0644)
	if err != nil {
		fmt.Printf("error on write links to file  err: %v", err)
	}

}

func saveStakesToBytesFile(stakes *map[cbd.AccNumber]uint64, fileName string) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(stakes)
	if err != nil {
		fmt.Printf("encode error:", err)
	}
	err = ioutil.WriteFile(fileName, network.Bytes(), 0644)
	if err != nil {
		fmt.Printf("error on write stakes to file  err: %v", err)
	}

}