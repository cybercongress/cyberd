package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ipfs/go-cid"

	//sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cybercongress/go-cyber/x/link"
)

var (
	_ sdk.Msg = &MsgAddJob{}
	_ sdk.Msg = &MsgRemoveJob{}

	_ sdk.Msg = &MsgChangeJobCID{}
	_ sdk.Msg = &MsgChangeJobLabel{}

	_ sdk.Msg = &MsgChangeJobCallData{}
	_ sdk.Msg = &MsgChangeJobGasPrice{}

	_ sdk.Msg = &MsgChangeJobPeriod{}
	_ sdk.Msg = &MsgChangeJobBlock{}
)

//______________________________________________________________________

type MsgAddJob struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Trigger  Trigger        `json:"trigger" yaml:"trigger"`
	Load     Load           `json:"load" yaml:"load"`
	Label    string         `json:"label" yaml:"label"`
	CID      link.Cid       `json:"cid" yaml:"cid"`
}

func NewMsgAddJob(
	address, contract sdk.AccAddress,
	trigger Trigger,
	load Load,
	label string,
	cid link.Cid,
) MsgAddJob {
	return MsgAddJob{
		Address:  address,
		Contract: contract,
		Trigger:  trigger,
		Load:     load,
		Label:    label,
		CID:      cid,
	}
}

func (msg MsgAddJob) Route() string { return RouterKey }

func (msg MsgAddJob) Type() string  { return "add_job" }

func (msg MsgAddJob) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgAddJob) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgAddJob) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Load.CallData) == 0 || len(msg.Load.CallData) > 256 {
		return ErrBadCallData
	}
	if msg.Load.GasPrice == 0 {
		return ErrBadGasPrice
	}
	if (msg.Trigger.Period == 0 && msg.Trigger.Block == 0) {
		return ErrBadTrigger
	}
	if (msg.Trigger.Period > 0 && msg.Trigger.Block > 0) {
		return ErrBadTrigger
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if _, err := cid.Decode(string(msg.CID)); err != nil {
		return link.ErrInvalidCid
	}

	return nil
}

//______________________________________________________________________

type MsgRemoveJob struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
}

func NewMsgRemoveJob(address, contract sdk.AccAddress, label string) MsgRemoveJob {
	return MsgRemoveJob{
		Address:  address,
		Contract: contract,
		Label: label,
	}
}

func (msg MsgRemoveJob) Route() string { return RouterKey }

func (msg MsgRemoveJob) Type() string  { return "remove_job" }

func (msg MsgRemoveJob) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgRemoveJob) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgRemoveJob) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobLabel struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	NewLabel string         `json:"new_label" yaml:"new_label"`
}

func NewMsgChangeLabel(
	address, contract sdk.AccAddress,
	label, newLabel string,
) MsgChangeJobLabel {
	return MsgChangeJobLabel{
		Address:  address,
		Contract: contract,
		Label:    label,
		NewLabel: newLabel,
	}
}

func (msg MsgChangeJobLabel) Route() string { return RouterKey }

func (msg MsgChangeJobLabel) Type() string { return "change_label" }

func (msg MsgChangeJobLabel) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobLabel) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobLabel) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if len(msg.NewLabel) == 0 || len(msg.NewLabel) > 32 {
		return ErrBadLabel
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobCID struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	CID link.Cid  		    `json:"cid" yaml:"cid"`
}

func NewMsgChangeCID(
	address, contract sdk.AccAddress,
	label string,
	cid link.Cid,
) MsgChangeJobCID {
	return MsgChangeJobCID{
		Address:  address,
		Contract: contract,
		Label: label,
		CID: cid,
	}
}

func (msg MsgChangeJobCID) Route() string { return RouterKey }

func (msg MsgChangeJobCID) Type() string { return "change_cid" }

func (msg MsgChangeJobCID) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobCID) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobCID) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if _, err := cid.Decode(string(msg.CID)); err != nil {
		return link.ErrInvalidCid
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobCallData struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	CallData string  		`json:"calldata" yaml:"calldata"`
}

func NewMsgChangeCallData(
	address, contract sdk.AccAddress,
	label string,
	calldata string,
) MsgChangeJobCallData {
	return MsgChangeJobCallData{
		Address:  address,
		Contract: contract,
		Label:    label,
		CallData: calldata,
	}
}

func (msg MsgChangeJobCallData) Route() string { return RouterKey }

func (msg MsgChangeJobCallData) Type() string { return "change_calldata" }

func (msg MsgChangeJobCallData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobCallData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobCallData) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if len(msg.CallData) == 0 || len(msg.CallData) > 256 {
		return ErrBadCallData
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobGasPrice struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	GasPrice uint64  		`json:"gasprice" yaml:"gasprice"`
}

func NewMsgChangeGasPrice(
	address, contract sdk.AccAddress,
	label string,
	gasprice uint64,
) MsgChangeJobGasPrice {
	return MsgChangeJobGasPrice{
		Address:  address,
		Contract: contract,
		Label:    label,
		GasPrice: gasprice,
	}
}

func (msg MsgChangeJobGasPrice) Route() string { return RouterKey }

func (msg MsgChangeJobGasPrice) Type() string { return "change_gasprice" }

func (msg MsgChangeJobGasPrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobGasPrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobGasPrice) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.GasPrice == 0 {
		return ErrBadGasPrice
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobPeriod struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	Period	 uint64  	    `json:"period" yaml:"period"`
}

func NewMsgChangeJobPeriod(
	address, contract sdk.AccAddress,
	label string,
	period uint64,
) MsgChangeJobPeriod {
	return MsgChangeJobPeriod{
		Address:  address,
		Contract: contract,
		Label:	  label,
		Period:   period,
	}
}

func (msg MsgChangeJobPeriod) Route() string { return RouterKey }

func (msg MsgChangeJobPeriod) Type() string { return "change_period" }

func (msg MsgChangeJobPeriod) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobPeriod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobPeriod) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Period == 0 {
		return ErrBadTrigger
	}

	return nil
}

//______________________________________________________________________

type MsgChangeJobBlock struct {
	Address  sdk.AccAddress `json:"address" yaml:"address"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`
	Label    string         `json:"label" yaml:"label"`
	Block	 uint64  		`json:"block" yaml:"block"`
}

func NewMsgChangeJobBlock(
	address, contract sdk.AccAddress,
	label string,
	block uint64,
) MsgChangeJobBlock {
	return MsgChangeJobBlock{
		Address:  address,
		Contract: contract,
		Label:    label,
		Block:    block,
	}
}

func (msg MsgChangeJobBlock) Route() string { return RouterKey }

func (msg MsgChangeJobBlock) Type() string { return "change_block" }

func (msg MsgChangeJobBlock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}

func (msg MsgChangeJobBlock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgChangeJobBlock) ValidateBasic() error {
	if len(msg.Address) == 0 {
		return ErrInvalidAddress
	}
	if msg.Contract.Empty() {
		return ErrEmptyContract
	}
	if len(msg.Label) == 0 || len(msg.Label) > 32 {
		return ErrBadLabel
	}
	if msg.Block == 0 {
		return ErrBadTrigger
	}

	return nil
}