package keeper

import (
	"encoding/binary"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cybercongress/go-cyber/x/bandwidth/exported"
	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
)

var _ exported.Keeper = &BaseAccBandwidthKeeper{}

type BaseAccBandwidthKeeper struct {
	key        *sdk.KVStoreKey
	paramSpace *params.Subspace
}

func NewAccBandwidthKeeper(key *sdk.KVStoreKey, paramSpace *params.Subspace) *BaseAccBandwidthKeeper {
	return &BaseAccBandwidthKeeper{
		key:        key,
		paramSpace: paramSpace,
	}
}

func (bk *BaseAccBandwidthKeeper) SetAccBandwidth(ctx sdk.Context, bandwidth types.AcсBandwidth) {
	bwBytes, _ := json.Marshal(bandwidth)
	ctx.KVStore(bk.key).Set(bandwidth.Address, bwBytes)
}

func (bk *BaseAccBandwidthKeeper) GetAccBandwidth(ctx sdk.Context, addr sdk.AccAddress) (bw types.AcсBandwidth) {
	bwBytes := ctx.KVStore(bk.key).Get(addr)
	if bwBytes == nil {
		return types.AcсBandwidth{
			Address:          addr,
			RemainedValue:    0,
			LastUpdatedBlock: ctx.BlockHeight(),
			MaxValue:         0,
		}
	}
	err := json.Unmarshal(bwBytes, &bw)
	if err != nil {
		// should not happen
		panic("bandwidth index is broken")
	}
	return
}

type BaseBlockSpentBandwidthKeeper struct {
	key *sdk.KVStoreKey
}

func NewBlockSpentBandwidthKeeper(key *sdk.KVStoreKey) BaseBlockSpentBandwidthKeeper {
	return BaseBlockSpentBandwidthKeeper{key: key}
}

func (bk BaseBlockSpentBandwidthKeeper) SetBlockSpentBandwidth(ctx sdk.Context, blockNumber uint64, value uint64) {
	keyAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(keyAsBytes, blockNumber)
	valueAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(valueAsBytes, value)
	ctx.KVStore(bk.key).Set(keyAsBytes, valueAsBytes)
}

func (bk BaseBlockSpentBandwidthKeeper) GetValuesForPeriod(ctx sdk.Context, period int64) map[uint64]uint64 {

	store := ctx.KVStore(bk.key)

	windowStart := ctx.BlockHeight() - period + 1
	if windowStart <= 0 { // check needed cause it will be casted to uint and can cause overflow
		windowStart = 1
	}

	key := make([]byte, 8)
	result := make(map[uint64]uint64)
	for blockNumber := windowStart; blockNumber <= ctx.BlockHeight(); blockNumber++ {
		binary.LittleEndian.PutUint64(key, uint64(blockNumber))
		value := binary.LittleEndian.Uint64(store.Get(key))
		result[uint64(blockNumber)] = value
	}

	return result
}
