package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgAddJob{}, "cyber/MsgAddJob", nil)
	cdc.RegisterConcrete(MsgRemoveJob{}, "cyber/MsgRemoveJob", nil)
	cdc.RegisterConcrete(MsgChangeJobCID{}, "cyber/MsgChangeJobCID", nil)
	cdc.RegisterConcrete(MsgChangeJobLabel{}, "cyber/MsgChangeJobLabel", nil)
	cdc.RegisterConcrete(MsgChangeJobCallData{}, "cyber/MsgChangeJobCallData", nil)
	cdc.RegisterConcrete(MsgChangeJobGasPrice{}, "cyber/MsgChangeJobGasPrice", nil)
	cdc.RegisterConcrete(MsgChangeJobPeriod{}, "cyber/MsgChangeJobPeriod", nil)
	cdc.RegisterConcrete(MsgChangeJobBlock{}, "cyber/MsgChangeJobBlock", nil)

}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
