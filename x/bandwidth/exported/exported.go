package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

type Keeper interface {
	SetAccBandwidth(ctx sdk.Context, bandwidth types.AcсBandwidth)
	GetAccBandwidth(ctx sdk.Context, address sdk.AccAddress) (bw types.AcсBandwidth)

	SetParams(ctx sdk.Context, params types.Params)
	GetParams(ctx sdk.Context) (params types.Params)
}

type BlockSpentBandwidthKeeper interface {
	SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64)
	GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64
}
