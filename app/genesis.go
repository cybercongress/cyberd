package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmwasm/wasmd/x/wasm"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cybercongress/cyberd/types/coin"
	"github.com/cybercongress/cyberd/x/bandwidth"
	"github.com/cybercongress/cyberd/x/rank"
)

const (
	defaultUnbondingTime = 60 * 60 * 24 * 3 * 7 * time.Second // 3 weeks
)

// State to Unmarshal
type GenesisState struct {
	AuthData      auth.GenesisState         `json:"auth"`
	BankData      bank.GenesisState         `json:"bank"`
	DistrData     distr.GenesisState        `json:"distribution"`
	MintData      mint.GenesisState         `json:"mint"`
	StakingData   staking.GenesisState      `json:"staking"`
	Pool          staking.Pool              `json:"pool"`
	SupplyData    supply.GenesisState       `json:"supply"`
	SlashingData  slashing.GenesisState     `json:"slashing"`
	GovData       gov.GenesisState          `json:"gov"`
	BandwidthData bandwidth.GenesisState    `json:"bandwidth"`
	RankData      rank.GenesisState         `json:"rank"`
	GenUtil       genutil.GenesisState      `json:"genutil"`
	Crisis        crisis.GenesisState       `json:"crisis"`
	Evidence      evidence.GenesisState		`json:"evidence"`
	WasmData      wasm.GenesisState         `json:"wasm"`
}

func (gs *GenesisState) GetAddresses() []sdk.AccAddress {
	addresses := make([]sdk.AccAddress, 0, len(gs.AuthData.Accounts))
	for _, acc := range gs.AuthData.Accounts {
		addresses = append(addresses, acc.GetAddress())
	}
	return addresses
}

func NewGenesisState(
	authData auth.GenesisState,
	stakingData staking.GenesisState, pool staking.Pool,
	mintData mint.GenesisState, distrData distr.GenesisState,
	govData gov.GenesisState, supplyData supply.GenesisState,
	slashingData slashing.GenesisState, bandwidthData bandwidth.GenesisState,
	rankData rank.GenesisState, crisisData crisis.GenesisState,
	evidenceData evidence.GenesisState, wasmData wasm.GenesisState,
) GenesisState {

	return GenesisState{
		AuthData:      authData,
		StakingData:   stakingData,
		Pool:          pool,
		SupplyData:    supplyData,
		MintData:      mintData,
		DistrData:     distrData,
		SlashingData:  slashingData,
		GovData:       govData,
		BandwidthData: bandwidthData,
		RankData:      rankData,
		Crisis:        crisisData,
		Evidence: 	   evidenceData,
		WasmData:      wasmData,
	}
}

// NewDefaultGenesisState generates the default state for cyberd.
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		AuthData: auth.GenesisState{
			Params: auth.Params{
				MaxMemoCharacters: 256,
				TxSigLimit: 10,
				TxSizeCostPerByte: 1,
				SigVerifyCostED25519: 1,
				SigVerifyCostSecp256k1: 1,
			},
		},
		BankData: bank.GenesisState{
			SendEnabled: true,
		},
		MintData: mint.GenesisState{
			Minter: mint.Minter{
				Inflation:        sdk.NewDecWithPrec(3, 2),
				AnnualProvisions: sdk.NewDec(0),
			},
			Params: mint.Params{
				MintDenom:           coin.CYB,
				InflationRateChange: sdk.NewDecWithPrec(10, 2),
				InflationMax:        sdk.NewDecWithPrec(15, 2),
				InflationMin:        sdk.NewDecWithPrec(1, 2),
				GoalBonded:          sdk.NewDecWithPrec(88, 2),
				BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
			},
		},
		StakingData: staking.GenesisState{
			Params: types.Params{
				UnbondingTime: defaultUnbondingTime,
				MaxValidators: 7,
				MaxEntries:    7,
				BondDenom:     coin.CYB,
			},
		},
		Pool: staking.Pool{
			NotBondedTokens: sdk.ZeroInt(),
			BondedTokens:    sdk.ZeroInt(),
		},
		SupplyData: supply.GenesisState{
			Supply: sdk.NewCoins(),
		},
		SlashingData: slashing.GenesisState{
			Params: slashing.Params{
				SignedBlocksWindow:      60 * 4, // ~20min
				DowntimeJailDuration:    1,
				MinSignedPerWindow:      sdk.NewDecWithPrec(80, 2),         // 80%
				SlashFractionDoubleSign: sdk.NewDecWithPrec(5, 2),          // 5%
				SlashFractionDowntime:   sdk.NewDec(5).Quo(sdk.NewDec(10000)), // 0.05%
			},
		},
		DistrData: distr.GenesisState{
			FeePool:             distr.InitialFeePool(),
			Params:              distr.Params{
				CommunityTax:        sdk.NewDecWithPrec(10, 2), // 10%
				BaseProposerReward:  sdk.NewDecWithPrec(1, 2),  // 1%
				BonusProposerReward: sdk.NewDecWithPrec(5, 2),  // 5%
				WithdrawAddrEnabled: true,
			},
			PreviousProposer:    nil,
		},
		GovData: gov.GenesisState{
			StartingProposalID: 1,
			DepositParams: gov.DepositParams{
				MinDeposit:       sdk.Coins{coin.NewCybCoin(500 * coin.Giga)},
				MaxDepositPeriod: 7200 * time.Second, // 2 hours
			},
			VotingParams: gov.VotingParams{
				VotingPeriod: 7200 * time.Second, // 2 hours
			},
			TallyParams: gov.TallyParams{
				Quorum:    sdk.NewDecWithPrec(15, 2),
				Threshold: sdk.NewDecWithPrec(15, 2),
				Veto:      sdk.NewDecWithPrec(15, 2),
			},
		},
		BandwidthData: bandwidth.GenesisState{
			Params:		bandwidth.DefaultParams(),
		},
		RankData: rank.GenesisState{
			Params: 	rank.DefaultParams(),
		},
		Crisis: crisis.GenesisState{
			ConstantFee: sdk.NewCoin(coin.CYB, sdk.NewInt(10000000000)),
		},
		Evidence: evidence.DefaultGenesisState(),
	}
}

// Create the core parameters for genesis initialization for cyberd
// note that the pubkey input is this machines pubkey
func CyberdAppGenState(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	genesisState GenesisState, err error) {

	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return genesisState, errors.New("there must be at least one genesis tx")
	}

	for i, genTx := range appGenTxs {
		var tx auth.StdTx
		if err := cdc.UnmarshalJSON(genTx, &tx); err != nil {
			return genesisState, err
		}
		msgs := tx.GetMsgs()
		if len(msgs) != 1 {
			return genesisState, errors.New(
				"must provide genesis StdTx with exactly 1 CreateValidator message")
		}
		if _, ok := msgs[0].(staking.MsgCreateValidator); !ok {
			return genesisState, fmt.Errorf(
				"genesis transaction %v does not contain a MsgCreateValidator", i)
		}
	}

	//genesisState.GenTxs = appGenTxs
	return genesisState, nil
}

//todo should be here?
func CyberdAppGenStateJSON(cdc *codec.Codec, genDoc tmtypes.GenesisDoc, appGenTxs []json.RawMessage) (
	appState json.RawMessage, err error) {
	// create the final app state
	genesisState, err := CyberdAppGenState(cdc, genDoc, appGenTxs)
	if err != nil {
		return nil, err
	}
	return codec.MarshalJSONIndent(cdc, genesisState)
}

// validateGenesisState ensures that the genesis state obeys the expected invariants
func validateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.AuthData.Accounts); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}
	if err := distr.ValidateGenesis(genesisState.DistrData); err != nil {
		return err
	}
	if err := gov.ValidateGenesis(genesisState.GovData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := mint.ValidateGenesis(genesisState.MintData); err != nil {
		return err
	}
	if err := supply.ValidateGenesis(genesisState.SupplyData); err != nil {
		return err
	}
	if err := bandwidth.ValidateGenesis(genesisState.BandwidthData); err != nil {
		return err
	}
	if err := rank.ValidateGenesis(genesisState.RankData); err != nil {
		return err
	}
	return staking.ValidateGenesis(genesisState.StakingData)
}

// Ensures that there are no duplicate accounts in the genesis state,
func validateGenesisStateAccounts(accs []exported.GenesisAccount) (err error) {
	addrMap := make(map[string]bool, len(accs))
	for i := 0; i < len(accs); i++ {
		acc := accs[i]
		strAddr := string(acc.GetAddress())
		if _, ok := addrMap[strAddr]; ok {
			return fmt.Errorf("duplicate account in genesis state: Address %v", acc.GetAddress())
		}
		addrMap[strAddr] = true
	}
	return
}

func LoadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc tmtypes.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}
