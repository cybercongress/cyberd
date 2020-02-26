// +build !cuda

package keeper

import (
	"github.com/cybercongress/go-cyber/x/rank/internal/types"

	"github.com/tendermint/tendermint/libs/log"
)

func calculateRankGPU(ctx *types.CalculationContext, logger log.Logger) []float64 {
	panic("Daemon compiled without gpu support, but started in gpu mode")
}
