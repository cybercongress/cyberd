package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyTxCost             = []byte("TxCost")
	KeyLinkMsgCost 		  = []byte("LinkMsgCost")
	KeyNonLinkMsgCost     = []byte("NonLinkMsgCost")
	KeyRecoveryPeriod     = []byte("RecoveryPeriod")
	KeyAdjustPricePeriod  = []byte("AdjustPricePeriod")
	KeyBaseCreditPrice    = []byte("BaseCreditPrice")
	KeyDesirableBandwidth = []byte("DesirableBandwidth")
	KeyMaxBlockBandwidth  = []byte("MaxBlockBandwidth")
)

// Params defines the parameters for the bandwidth module.
// TODO move int64 -> uint64 for params
type Params struct {
	TxCost             int64   `json:"tx_cost" yaml:"tx_cost"`
	LinkMsgCost        int64   `json:"link_msg_cost" yaml:"link_msg_cost"`
	NonLinkMsgCost     int64   `json:"non_link_msg_cost" yaml:"non_link_msg_cost"`
	RecoveryPeriod     int64   `json:"recovery_period" yaml:"recovery_period"`
	AdjustPricePeriod  int64   `json:"adjust_price_period" yaml:"adjust_price_period"`
	BaseCreditPrice    sdk.Dec `json:"base_credit_price" yaml:"base_credit_price"`
	DesirableBandwidth int64   `json:"desirable_bandwidth" yaml:"desirable_bandwidth"`
	MaxBlockBandwidth  uint64  `json:"max_block_bandwidth" yaml:"max_block_bandwidth"`
}

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(
	txCost 			   int64,
	linkMsgCost        int64,
	nonLinkMsgCost     int64,
	recoveryPeriod     int64,
	adjustPricePeriod  int64,
	baseCreditPrice    sdk.Dec,
	desirableBandwidth int64,
	maxBlockBandwidth  uint64,
) Params {

	return Params{
		TxCost:             txCost,
		LinkMsgCost:        linkMsgCost,
		NonLinkMsgCost:     nonLinkMsgCost,
		RecoveryPeriod:     recoveryPeriod,
		AdjustPricePeriod:  adjustPricePeriod,
		BaseCreditPrice:    baseCreditPrice,
		DesirableBandwidth: desirableBandwidth,
		MaxBlockBandwidth:  maxBlockBandwidth,
	}
}

func DefaultParams() Params {
	return Params{
		TxCost:             int64(3000),
		LinkMsgCost:        int64(1000),
		NonLinkMsgCost:     int64(5000),
		RecoveryPeriod:     int64(1600000),
		AdjustPricePeriod:  int64(10),
		BaseCreditPrice:    sdk.NewDec(50),
		DesirableBandwidth: int64(200000000),
		MaxBlockBandwidth:  uint64(200000000*10/16000),
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateTxCost(p.TxCost); err != nil {
		return err
	}
	if err := validateLinkMsgCost(p.TxCost); err != nil {
		return err
	}
	if err := validateNonLinkMsgCost(p.NonLinkMsgCost); err != nil {
		return err
	}
	if err := validateRecoveryPeriod(p.RecoveryPeriod); err != nil {
		return err
	}
	if err := validateAdjustPricePeriod(p.AdjustPricePeriod); err != nil {
		return err
	}
	if err := validateBaseCreditPrice(p.BaseCreditPrice); err != nil {
		return err
	}
	if err := validateDesirableBandwidth(p.DesirableBandwidth); err != nil {
		return err
	}
	if err := validateMaxBlockBanwidth(p.MaxBlockBandwidth); err != nil {
		return err
	}

	return nil
}

func (p Params) String() string {

	return fmt.Sprintf(`Bandwidth params:
  LinkMsgCost:        %d
  TxCost:			  %d
  NonLinkMsgCost:	  %d
  RecoveryPeriod:     %d
  AdjustPricePeriod:  %d
  BaseCreditPrice:    %d
  DesirableBandwidth: %d
  MaxBlockBandidth:   %d
`,
		p.LinkMsgCost, p.RecoveryPeriod, p.AdjustPricePeriod,
		p.BaseCreditPrice, p.DesirableBandwidth, p.MaxBlockBandwidth,
		p.TxCost, p.NonLinkMsgCost,
	)
}

// TODO improve validations for parameters
func validateTxCost(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(10) {
		return fmt.Errorf("tx cost too low: %d", v)
	}

	return nil
}

func validateLinkMsgCost(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(10) {
		return fmt.Errorf("link msg cost too low: %d", v)
	}

	return nil
}

func validateNonLinkMsgCost(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(10) {
		return fmt.Errorf("non link msg too low: %d", v)
	}

	return nil
}

func validateRecoveryPeriod(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(100) {
		return fmt.Errorf("recovery period too low: %d", v)
	}

	return nil
}

func validateAdjustPricePeriod(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(2) {
		return fmt.Errorf("adjust price period too low: %d", v)
	}

	return nil
}

func validateBaseCreditPrice(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.OneDec()) {
		return fmt.Errorf("base credit price too low: %s", v)
	}

	return nil
}

func validateDesirableBandwidth(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(10000) {
		return fmt.Errorf("desirable bandwidth too low: %d", v)
	}

	return nil
}

func validateMaxBlockBanwidth(i interface{}) error {
	v, ok := i.(uint64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= uint64(100) {
		return fmt.Errorf("max block bandwidth too low: %d", v)
	}

	return nil
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyTxCost, &p.TxCost, validateTxCost),
		params.NewParamSetPair(KeyLinkMsgCost, &p.LinkMsgCost, validateLinkMsgCost),
		params.NewParamSetPair(KeyNonLinkMsgCost, &p.NonLinkMsgCost, validateNonLinkMsgCost),
		params.NewParamSetPair(KeyRecoveryPeriod, &p.RecoveryPeriod, validateRecoveryPeriod),
		params.NewParamSetPair(KeyAdjustPricePeriod, &p.AdjustPricePeriod, validateAdjustPricePeriod),
		params.NewParamSetPair(KeyBaseCreditPrice, &p.BaseCreditPrice, validateBaseCreditPrice),
		params.NewParamSetPair(KeyDesirableBandwidth, &p.DesirableBandwidth, validateDesirableBandwidth),
		params.NewParamSetPair(KeyMaxBlockBandwidth, &p.MaxBlockBandwidth, validateMaxBlockBanwidth),
	}
}
