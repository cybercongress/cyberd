package keeper

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cybercongress/go-cyber/types"

	"github.com/cybercongress/go-cyber/x/cron/types"
	"github.com/cybercongress/go-cyber/x/link"
)

// Keeper of the power store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	wasmKeeper 	  wasm.Keeper
	supplyKeeper  types.SupplyKeeper
	accountKeeper types.AccountKeeper
	proxyKeeper   types.BankKeeper
	paramspace    types.ParamSubspace
}

func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey,
	wk wasm.Keeper,
	sk types.SupplyKeeper, bk types.BankKeeper,
	ak types.AccountKeeper, paramspace types.ParamSubspace,
) Keeper {

	keeper := Keeper{
		storeKey:   key,
		cdc:        cdc,
		wasmKeeper: wk,
		supplyKeeper: sk,
		proxyKeeper:   bk,
		accountKeeper: ak,
		paramspace: paramspace.WithKeyTable(types.ParamKeyTable()),
	}
	return keeper
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) AddJob(
	ctx sdk.Context, creator, contract sdk.AccAddress,
	trigger types.Trigger, load types.Load,
	label string, cid link.Cid,
) error {

	if trigger.Block != 0 && ctx.BlockHeight() > int64(trigger.Block) {
		return types.ErrBadTrigger
	}

	k.SetJob(ctx, types.NewJob(
		creator, contract,
		trigger, load,
		label, cid,
	))
	k.SetJobStats(ctx, contract, creator, label,
		types.NewStats(0,0, 0, uint64(ctx.BlockHeight())),
	)

	return nil
}

func (k Keeper) RemoveJob(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.DeleteJob(ctx, contract, creator, label)

	return nil
}

func (k Keeper) ChangeJobCID(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	cid link.Cid,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		job.Trigger, job.Load,
		job.Label, cid,
	))

	return nil
}

func (k Keeper) ChangeJobLabel(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	labelNew string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}
	jobStats, _ := k.GetJobStats(ctx, contract, creator, label)

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		job.Trigger, job.Load,
		labelNew, job.CID,
	))

	k.SetJobStats(ctx, job.Contract, job.Creator, labelNew,
		types.NewStats(
			jobStats.Calls, jobStats.Fees, jobStats.Fees, jobStats.LastBlock,
	))

	k.DeleteJob(ctx, contract, creator, label)
	k.DeleteJobStats(ctx, contract, creator, label)

	return nil
}

func (k Keeper) ChangeJobCallData(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	calldata string,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		job.Trigger, types.NewLoad(calldata, job.Load.GasPrice),
		job.Label, job.CID,
	))

	return nil
}

func (k Keeper) ChangeJobGasPrice(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	gasprice uint64,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		job.Trigger, types.NewLoad(job.Load.CallData, gasprice),
		job.Label, job.CID,
	))

	return nil
}

func (k Keeper) ChangeJobPeriod(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	period uint64,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		types.NewTrigger(period, job.Trigger.Block, sdk.ZeroDec()), job.Load,
		job.Label, job.CID,
	))

	return nil
}

func (k Keeper) ChangeJobBlock(
	ctx sdk.Context, creator, contract sdk.AccAddress, label string,
	block uint64,
) error {
	job, found := k.GetJob(ctx, contract, creator, label)
	if !found {
		return types.ErrJobNotExist
	}

	if string(job.Creator) != string(creator) {
		return types.ErrNotAuthorized
	}

	k.SetJob(ctx, types.NewJob(
		job.Creator, job.Contract,
		types.NewTrigger(job.Trigger.Period, block, sdk.ZeroDec()), job.Load,
		job.Label, job.CID,
	))

	return nil
}


//______________________________________________________________________


func (k Keeper) MaxJobs(ctx sdk.Context) (res uint16) {
	k.paramspace.Get(ctx, types.KeyMaxSlots, &res)
	return
}

func (k Keeper) MaxGas(ctx sdk.Context) (res uint32) {
	k.paramspace.Get(ctx, types.KeyMaxGas, &res)
	return
}

func (k Keeper) FeeTTL(ctx sdk.Context) (res uint16) {
	k.paramspace.Get(ctx, types.KeyFeeTTL, &res)
	return
}

//______________________________________________________________________

func (k Keeper) SetJob(ctx sdk.Context, job types.Job) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalJob(k.cdc, job)
	store.Set(types.GetJobKey(job.Contract, job.Creator, job.Label), b)
}

func (k Keeper) DeleteJob(ctx sdk.Context, contract, creator sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobKey(contract, creator, label))
}

func (k Keeper) SetJobs(ctx sdk.Context, jobs types.Jobs) error {
	for _, job := range jobs {
		k.SetJob(ctx, types.NewJob(
			job.Creator, job.Contract,
			job.Trigger, job.Load,
			job.Label, job.CID,
		))
	}
	return nil
}

func (k Keeper) SetJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string, stats types.JobStats) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalJobStats(k.cdc, stats)
	store.Set(types.GetJobStatsKey(contract, creator, label), b)
}

func (k Keeper) DeleteJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetJobStatsKey(contract, creator, label))
}

//______________________________________________________________________

func (k Keeper) GetJob(ctx sdk.Context, contract, creator sdk.AccAddress, label string) (job types.Job, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobKey(contract, creator, label)

	value := store.Get(key)
	if value == nil {
		return job, false
	}

	job = types.MustUnmarshalJob(k.cdc, value)

	return job, true
}

func (k Keeper) GetAllJobs(ctx sdk.Context) (jobs types.Jobs) {
	k.IterateAllJobs(ctx, func(job types.Job) bool {
		jobs = append(jobs, job)
		return false
	})

	return jobs
}

func (k Keeper) IterateAllJobs(ctx sdk.Context, cb func(job types.Job) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.JobKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		job := types.MustUnmarshalJob(k.cdc, iterator.Value())
		if cb(job) {
			break
		}
	}
}

func (k Keeper) GetJobStats(ctx sdk.Context, contract, creator sdk.AccAddress, label string) (stats types.JobStats, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetJobStatsKey(contract, creator, label)

	value := store.Get(key)
	if value == nil {
		return stats, false
	}

	stats = types.MustUnmarshalJobStats(k.cdc, value)

	return stats, true
}

//______________________________________________________________________

func (k Keeper) EndBlocker(ctx sdk.Context) {

	jobs := k.GetAllJobs(ctx)
	jobs.Sort()

	maxGas := k.MaxGas(ctx)
	gasBeforeCron := ctx.GasMeter().GasConsumed()
	gasUsedTotal := sdk.Gas(0)

	fmt.Println("[*] Jobs in queue: ", jobs.Len())

	for i, job := range jobs {
		if (job.Trigger.Period != 0 && ctx.BlockHeight()%int64(job.Trigger.Period) == 0) || (job.Trigger.Period == 0 && ctx.BlockHeight() == int64(job.Trigger.Block)) {
			price := job.Load.GasPrice

			fmt.Println("[!] Job started. ID: ", i)
			fmt.Println("-- Job gas price: ", price)

			cacheContext, writeFn := ctx.CacheContext()
			gasBefore := ctx.GasMeter().GasConsumed()
			js, _ := k.GetJobStats(ctx, job.Contract, job.Creator, job.Label)
			_, errExecute := k.wasmKeeper.Execute(
				cacheContext, job.Contract, job.Creator,
				[]byte(job.Load.CallData), sdk.Coins{},
			)
			gasUsed := ctx.GasMeter().GasConsumed() - gasBefore

			if uint64(gasUsedTotal) + uint64(gasUsed) > uint64(maxGas) {
				break
			} else {
				gasUsedTotal += gasUsed
			}

			amtGasFee := uint64(gasUsed) * price
			amtTTLFee := uint64(ctx.BlockHeight() - int64(js.LastBlock))*uint64(k.FeeTTL(ctx))
			amtTotalFee := amtGasFee + amtTTLFee

			fmt.Println("-- Gas used: ", gasUsed)
			fmt.Println("-- Gas Fee: ", amtGasFee)
			fmt.Println("-- TTL Fee: ", amtTTLFee)
			fmt.Println("-- Total Fee: ", amtTotalFee)

			fee := sdk.NewCoin(ctypes.CYB, sdk.NewInt(int64(amtTotalFee)))

			errSend := k.proxyKeeper.SendCoins(
				ctx, job.Contract, k.supplyKeeper.GetModuleAddress(auth.FeeCollectorName), sdk.NewCoins(fee))
			if errSend != nil {
				k.DeleteJob(ctx, job.Contract, job.Creator, job.Label)
				k.DeleteJobStats(ctx, job.Contract, job.Creator, job.Label)

				fmt.Println("[!] Not enough contract balance, state not applied, job killed! ID: ", i)
				continue
			}

			if errExecute != nil {
				fmt.Println("[!] Job failed, state not applied. ID: ", i)
				fmt.Println("-- Error: ", errExecute)
			} else {
				writeFn()
				fmt.Println("[!] Job finished, state applied. ID: ", i)
			}

			k.SetJobStats(ctx, job.Contract, job.Creator, job.Label,
				types.NewStats(
					js.Calls+1, js.Fees+amtTotalFee,
					js.Gas+gasUsed, uint64(ctx.BlockHeight())),
			)

			if ctx.BlockHeight() == int64(job.Trigger.Block) {
				k.DeleteJob(ctx, job.Contract, job.Creator, job.Label)
				k.DeleteJobStats(ctx, job.Contract, job.Creator, job.Label)

				fmt.Println("[!] Job executed at given block, killed! ID: ", i)
			}
		}
	}

	gasAfterCron := ctx.GasMeter().GasConsumed()
	fmt.Println("[*] Total cron gas used: ", gasAfterCron-gasBeforeCron)
}
