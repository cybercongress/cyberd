package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	DefaultParamspace = ModuleName

	DefaultMaxSlots   uint16 = 10
	DefaultMaxGas     uint32 = 1000000
	DefaultFeeTTL     uint16 = 50
)

var (
	KeyMaxSlots   = []byte("MaxSlots")
	KeyMaxGas     = []byte("MaxGas")
	KeyFeeTTL     = []byte("FeeTTL")
)

func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

type Params struct {
	MaxSlots   uint16 `json:"max_slots" yaml:"max_slots"`
	MaxGas     uint32 `json:"max_gas" yaml:"max_gas"`
	FeeTTL     uint16 `json:"fee_ttl" yaml:"fee_ttl"`
}

func NewParams(
	maxSlots uint16,
	maxGas   uint32,
	feeTTL   uint16,
) Params {
	return Params{
		MaxSlots: maxSlots,
		MaxGas:   maxGas,
		FeeTTL:   feeTTL,
	}
}

func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Max Slots:  %d
  Max Gas:    %d
  Fee TTL:    %d
`, p.MaxSlots, p.MaxGas, p.FeeTTL)
}

func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxSlots, &p.MaxSlots, validateMaxSlots),
		params.NewParamSetPair(KeyMaxGas, &p.MaxGas, validateMaxGas),
		params.NewParamSetPair(KeyFeeTTL, &p.FeeTTL, validateFeeTTL),
	}
}

func DefaultParams() Params {
	return NewParams(DefaultMaxSlots, DefaultMaxGas, DefaultFeeTTL)
}

func validateMaxSlots(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max slots must be positive: %d", v)
	}

	return nil
}

func validateMaxGas(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max gas must be positive: %d", v)
	}

	return nil
}

func validateFeeTTL(i interface{}) error {
	v, ok := i.(uint16)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("fee ttl must be positive: %d", v)
	}

	return nil
}
