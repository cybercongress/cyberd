package keeper

import (
	"github.com/cybercongress/go-cyber/x/rank/internal/types"

	"github.com/tendermint/tendermint/libs/log"

	"encoding/hex"
	"time"
)

func CalculateRank(ctx *types.CalculationContext, unit types.ComputeUnit, logger log.Logger) (rank types.Rank) {
	start := time.Now()
	if unit == types.CPU {
		//used only for development
		rank = types.NewRank(calculateRankCPU(ctx, logger), logger, ctx.FullTree)
	} else {
		rank = types.NewRank(calculateRankGPU(ctx, logger), logger, ctx.FullTree)
	}
	logger.Info(
		"Rank calculated", "time", time.Since(start), "links", ctx.LinksCount, "objects", ctx.CidsCount,
		"hash", hex.EncodeToString(rank.MerkleTree.RootHash()),
	)
	return
}

func CalculateRankInParallel(
	ctx *types.CalculationContext, rankChan chan types.Rank, err chan error, unit types.ComputeUnit, logger log.Logger,
) {
	defer func() {
		if r := recover(); r != nil {
			err <- r.(error)
		}
	}()

	rank := CalculateRank(ctx, unit, logger)
	rankChan <- rank
}
