package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Parameter keys
var (
	KeyCalculationPeriod = []byte("CalculationPeriod")
	KeyDampingFactor     = []byte("DampingFactor")
	KeyTolerance		 = []byte("Tolerance")
)

// Params defines the parameters for the rank module.
type Params struct {
	CalculationPeriod int64   `json:"calculation_period" yaml:"calculation_period"`
	DampingFactor 	  sdk.Dec `json:"damping_factor" yaml:"damping_factor"`
	Tolerance		  sdk.Dec `json:"tolerance" yaml:"tolerance"`
}

// ParamKeyTable for rank module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object
func NewParams(
	calculationPeriod int64,
	dampingFactor sdk.Dec,
	tolerance sdk.Dec,
) Params {

	return Params{
		CalculationPeriod: calculationPeriod,
		DampingFactor:     dampingFactor,
		Tolerance:		   tolerance,
	}
}

// NewDefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CalculationPeriod: int64(10),
		DampingFactor:	   sdk.NewDecWithPrec(85, 2),
		Tolerance:         sdk.NewDecWithPrec(1, 3),
	}
}

func (p Params) Validate() error {
	if err := validateCalculationPeriod(p.CalculationPeriod); err != nil {
		return err
	}
	if err := validateDampingFactor(p.DampingFactor); err != nil {
		return err
	}
	if err := validateTolerance(p.Tolerance); err != nil {
		return err
	}

	return nil
}

func (p Params) String() string {

	return fmt.Sprintf(`Rank params:
  CalculationPeriod: %d
  DampingFactor:	 %d
  Tolerance:		 %d
`,
		p.CalculationPeriod, p.DampingFactor, p.Tolerance,
	)
}

func validateCalculationPeriod(i interface{}) error {
	v, ok := i.(int64)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= int64(2) {
		return fmt.Errorf("calculation period too low: %d", v)
	}

	return nil
}

func validateDampingFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("damping factor should be positive: %s", v)
	}

	if v.GTE(sdk.OneDec()) {
		return fmt.Errorf("damping factor should be < 1, is: %s", v)
	}

	return nil
}

func validateTolerance(i interface{}) error {
	v, ok := i.(sdk.Dec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.GT(sdk.NewDecWithPrec(1, 3)) {
		return fmt.Errorf("tolerance too low: %s", v)
	}

	if v.LT(sdk.NewDecWithPrec(1, 5)) {
		return fmt.Errorf("tolerance too big: %s", v)
	}

	return nil
}

func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		params.NewParamSetPair(KeyCalculationPeriod, &p.CalculationPeriod, validateCalculationPeriod),
		params.NewParamSetPair(KeyDampingFactor, &p.DampingFactor, validateDampingFactor),
		params.NewParamSetPair(KeyTolerance, &p.Tolerance, validateTolerance),
	}
}
