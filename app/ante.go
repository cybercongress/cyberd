package app

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/x/bank"

	//"github.com/Workiva/go-datastructures/threadsafe/err"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/link"
)

func NewAnteHandler(
	ak keeper.AccountKeeper,
	bankKeeper bank.Keeper,
	supplyKeeper types.SupplyKeeper,
	abk *bandwidth.BandwidthMeter,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(),
		ante.NewMempoolFeeDecorator(),
		ante.NewValidateBasicDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewSetPubKeyDecorator(ak),
		ante.NewValidateSigCountDecorator(ak),
		//ante.NewDeductFeeDecorator(ak, supplyKeeper),
		NewDeductBandwidthDecorator(ak, bankKeeper, supplyKeeper, abk),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak),
		ante.NewIncrementSequenceDecorator(ak),
	)
}

var (
	_ FeeTx = (*types.StdTx)(nil)
)

type FeeTx interface {
	sdk.Tx
	GetGas() uint64
	GetFee() sdk.Coins
	FeePayer() sdk.AccAddress
}


type DeductBandwidthDecorator struct {
	ak  auth.AccountKeeper
	bk 	bank.Keeper
	sk	types.SupplyKeeper
	bm	*bandwidth.BandwidthMeter
}

func NewDeductBandwidthDecorator(ak auth.AccountKeeper, bk bank.Keeper, sk types.SupplyKeeper, bm *bandwidth.BandwidthMeter) DeductBandwidthDecorator {
	return DeductBandwidthDecorator{
		ak: ak,
		bk: bk,
		sk: sk,
		bm: bm,
	}
}

func (dbd DeductBandwidthDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	nativeFlag := false
	wasmExecuteFlag := false
	ai2pay := sdk.AccAddress{}

	for _, msg := range tx.GetMsgs() {
		if (msg.Route() != link.RouterKey && msg.Route() != wasm.RouterKey) {
			nativeFlag = true
			break
		}
		if (msg.Route() == wasm.RouterKey && msg.Type() != "execute") {
			nativeFlag = true
			break
		}
		if (msg.Route() == wasm.RouterKey && msg.Type() == "execute") {
			executeTx, ok := msg.(wasm.MsgExecuteContract)
			if !ok {
				nativeFlag = true
				break
			}
			wasmExecuteFlag = true
			ai2pay = executeTx.Contract
		}

	}

	feeTx, ok := tx.(FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feePayerAcc := dbd.ak.GetAccount(ctx, feePayer)

	if feePayerAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	if nativeFlag {
		fmt.Println("[*] Native fee tx routing")
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(dbd.sk, dbd.bk, ctx, feePayerAcc, feeTx.GetFee(), nil)
			if err != nil {
				return ctx, err
			}
		}

		return next(ctx, tx, simulate)
	}
	if wasmExecuteFlag {
		fmt.Println("[*] Execute fee split tx routing")
		if !feeTx.GetFee().IsZero() {
			err = DeductFees(dbd.sk, dbd.bk, ctx, feePayerAcc, feeTx.GetFee(), ai2pay)
			if err != nil {
				return ctx, err
			}
		}

		return next(ctx, tx, simulate)
	}
	fmt.Println("[*] Bandwidth link tx routing")

	txCost := dbd.bm.GetPricedTxCost(ctx, tx)
	accountBandwidth := dbd.bm.GetCurrentAccountBandwidth(ctx, feePayerAcc.GetAddress())

	currentBlockSpentBandwidth := dbd.bm.GetCurrentBlockSpentBandwidth(ctx)
	maxBlockBandwidth := dbd.bm.GetMaxBlockBandwidth(ctx)

	if !accountBandwidth.HasEnoughRemained(txCost) {
		return ctx, bandwidth.ErrNotEnoughBandwidth
	} else if (uint64(txCost) + currentBlockSpentBandwidth) > maxBlockBandwidth  {
		return ctx, bandwidth.ErrExceededMaxBlockBandwidth
	} else {
		fmt.Println("-- bandwidth consumed: ", txCost)
		dbd.bm.ConsumeAccountBandwidth(ctx, accountBandwidth, txCost)
		dbd.bm.AddToBlockBandwidth(txCost)
	}

	return next(ctx, tx, simulate)
}

// DeductFees deducts fees from the given account.
//
// NOTE: We could use the BankKeeper (in addition to the AccountKeeper, because
// the BankKeeper doesn't give us accounts), but it seems easier to do this.
func DeductFees(supplyKeeper types.SupplyKeeper, bankKeeper bank.Keeper, ctx sdk.Context, acc exported.Account, fees sdk.Coins, ai sdk.AccAddress) error {
	blockTime := ctx.BlockHeader().Time
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	if ai == nil {
		fmt.Println("-- fee native payed: ", fees)
		err := supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	} else {
		feeInCYB := sdk.NewDec(fees.AmountOf(ctypes.CYB).Int64())
		toContract := feeInCYB.Mul(sdk.NewDecWithPrec(80,2))
		toValidators := feeInCYB.Sub(toContract)

		toValidatorsAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toValidators.RoundInt()))
		toContractAmount := sdk.NewCoins(sdk.NewCoin(ctypes.CYB, toContract.RoundInt()))

		fmt.Println("-- fee split contract payed: ", toContractAmount)
		err := bankKeeper.SendCoins(ctx, acc.GetAddress(), ai, toContractAmount)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}

		fmt.Println("-- fee split validator payed: ", toValidatorsAmount)
		err = supplyKeeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, toValidatorsAmount)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
		}
	}

	return nil
}

