package cron

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/x/cron/keeper"
	"github.com/cybercongress/go-cyber/x/cron/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgAddJob:
			return handleMsgAddJob(ctx, k, msg)
		case types.MsgRemoveJob:
			return handleMsgRemoveJob(ctx, k, msg)
		case types.MsgChangeJobCID:
			return handleMsgChangeJobCID(ctx, k, msg)
		case types.MsgChangeJobLabel:
			return handleMsgChangeJobAlias(ctx, k, msg)
		case types.MsgChangeJobCallData:
			return handleMsgChangeJobCallData(ctx, k, msg)
		case types.MsgChangeJobGasPrice:
			return handleMsgChangeJobGasPrice(ctx, k, msg)
		case types.MsgChangeJobPeriod:
			return handleMsgChangeJobPeriod(ctx, k, msg)
		case types.MsgChangeJobBlock:
			return handleMsgChangeJobBlock(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgAddJob(ctx sdk.Context, k keeper.Keeper, msg types.MsgAddJob) (*sdk.Result, error) {
	err := k.AddJob(
		ctx, msg.Address, msg.Contract,
		msg.Trigger, msg.Load,
		msg.Label, msg.CID,
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddJob,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobTrigger, msg.Trigger.String()),
			sdk.NewAttribute(types.AttributeKeyJobLoad, msg.Load.String()),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
			sdk.NewAttribute(types.AttributeKeyJobCID, string(msg.CID)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgRemoveJob(ctx sdk.Context, k keeper.Keeper, msg types.MsgRemoveJob) (*sdk.Result, error) {
	err := k.RemoveJob(ctx, msg.Address, msg.Contract, msg.Label)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveJob,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobCID(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobCID) (*sdk.Result, error) {
	err := k.ChangeJobCID(ctx, msg.Address, msg.Contract, msg.Label, msg.CID)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobCID,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract,msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobCID, string(msg.CID)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobAlias(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobLabel) (*sdk.Result, error) {
	err := k.ChangeJobLabel(ctx, msg.Address, msg.Contract, msg.Label, msg.NewLabel)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobLabel,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobLabel, msg.Label),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobCallData(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobCallData) (*sdk.Result, error) {
	err := k.ChangeJobCallData(ctx, msg.Address, msg.Contract, msg.Label, msg.CallData)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobCallData,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobCallData, msg.CallData),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobGasPrice(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobGasPrice) (*sdk.Result, error) {
	err := k.ChangeJobGasPrice(ctx, msg.Address, msg.Contract, msg.Label, msg.GasPrice)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobGasPrice,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobGasPrice, string(msg.GasPrice)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobPeriod(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobPeriod) (*sdk.Result, error) {
	err := k.ChangeJobPeriod(ctx, msg.Address, msg.Contract, msg.Label, msg.Period)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobPeriod,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobPeriod, string(msg.Period)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgChangeJobBlock(ctx sdk.Context, k keeper.Keeper, msg types.MsgChangeJobBlock) (*sdk.Result, error) {
	err := k.ChangeJobBlock(ctx, msg.Address, msg.Contract, msg.Label, msg.Block)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeChangeJobBlock,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyJobContract, msg.Contract.String()),
			sdk.NewAttribute(types.AttributeKeyJobBlock, string(msg.Block)),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

