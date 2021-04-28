package app

import (
	"time"

	"github.com/CosmWasm/wasmd/x/wasm"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	ctypes "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/cron"
	"github.com/cybercongress/go-cyber/x/energy"
	"github.com/cybercongress/go-cyber/x/rank"
)

const (
	productionUnbondingTime = 60 * 60 * 24 * 3 * 7 * time.Second // 3 weeks
	developmentUnbondingTime = 60 * 60 * 2 * time.Second // 2 hours

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
	GenUtil       genutil.GenesisState      `json:"genutil"`
	Crisis        crisis.GenesisState       `json:"crisis"`
	Evidence      evidence.GenesisState		`json:"evidence"`
	BandwidthData bandwidth.GenesisState    `json:"bandwidth"`
	RankData      rank.GenesisState         `json:"rank"`
	EnergyData    energy.GenesisState       `json:"energy"`
	CronData      cron.GenesisState         `json:"cron"`
	WasmData      wasm.GenesisState         `json:"wasm"`
}

//func (gs *GenesisState) GetAddresses() []sdk.AccAddress {
//	addresses := make([]sdk.AccAddress, 0, len(gs.AuthData.Accounts))
//	for _, acc := range gs.AuthData.Accounts {
//		addresses = append(addresses, acc.GetAddress())
//	}
//	return addresses
//}

func NewGenesisState(
	authData auth.GenesisState,
	stakingData staking.GenesisState, pool staking.Pool,
	mintData mint.GenesisState, distrData distr.GenesisState,
	govData gov.GenesisState, supplyData supply.GenesisState,
	slashingData slashing.GenesisState, bandwidthData bandwidth.GenesisState,
	rankData rank.GenesisState, energyData energy.GenesisState, cronData cron.GenesisState,
	crisisData crisis.GenesisState, evidenceData evidence.GenesisState, wasmData wasm.GenesisState,
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
		EnergyData:    energyData,
		CronData:	   cronData,
		Crisis:        crisisData,
		Evidence: 	   evidenceData,
		WasmData:      wasmData,
	}
}

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
				MintDenom:           ctypes.CYB,
				InflationRateChange: sdk.NewDecWithPrec(10, 2),
				InflationMax:        sdk.NewDecWithPrec(15, 2),
				InflationMin:        sdk.NewDecWithPrec(1, 2),
				GoalBonded:          sdk.NewDecWithPrec(88, 2),
				BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
			},
		},
		StakingData: staking.GenesisState{
			Params: types.Params{
				UnbondingTime: developmentUnbondingTime,
				MaxValidators: 7,
				MaxEntries:    7,
				BondDenom:     ctypes.CYB,
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
				MinDeposit:       sdk.Coins{ctypes.NewCybCoin(500 * ctypes.Giga)},
				MaxDepositPeriod: 7200 * time.Second, // 2 hours
			},
			VotingParams: gov.VotingParams{
				VotingPeriod: 7200 * time.Second, // 2 hours
			},
			TallyParams: gov.TallyParams{
				Quorum:    sdk.NewDecWithPrec(334, 3),
				Threshold: sdk.NewDecWithPrec(5, 1),
				Veto:      sdk.NewDecWithPrec(334, 3),
			},
		},
		BandwidthData: bandwidth.GenesisState{
			Params:		bandwidth.DefaultParams(),
		},
		RankData: rank.GenesisState{
			Params: 	rank.DefaultParams(),
		},
		EnergyData: energy.DefaultGenesisState(),
		CronData: cron.DefaultGenesisState(),
		Crisis: crisis.GenesisState{
			ConstantFee: sdk.NewCoin(ctypes.CYB, sdk.NewInt(ctypes.Giga*10)),
		},
		Evidence: evidence.DefaultGenesisState(),
		WasmData: wasm.GenesisState{
				Params: wasm.DefaultParams(),
		},
	}
}