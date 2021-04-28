package types

import (
	"fmt"
	"sort"
	"strings"

	//"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link"
)

type Job struct {
	Creator  sdk.AccAddress `json:"creator" yaml:"creator"`
	Contract sdk.AccAddress `json:"contract" yaml:"contract"`

	Trigger  Trigger	    `json:"trigger" yaml:"trigger"`
	Load	 Load			`json:"load" yaml:"load"`

	Label string   			`json:"label" yaml:"label"`
	CID   link.Cid 			`json:"cid" yaml:"cid"`
}

func NewJob(
	creator, contract sdk.AccAddress,
	trigger Trigger,
	load Load,
	label string,
	cid link.Cid,
) Job {
	return Job{
		Creator:  creator,
		Contract: contract,
		Trigger:  trigger,
		Load:     load,
		Label:    label,
		CID:      cid,
	}
}

func MustMarshalJob(cdc *codec.Codec, route Job) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(route)
}

func MustUnmarshalJob(cdc *codec.Codec, value []byte) Job {
	job, err := UnmarshalJob(cdc, value)
	if err != nil {
		panic(err)
	}
	return job
}

func UnmarshalJob(cdc *codec.Codec, value []byte) (job Job, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &job)
	return job, err
}

func (r Job) GetCreator()       sdk.AccAddress { return r.Creator }
func (r Job) GetContract()      sdk.AccAddress { return r.Contract }
func (r Job) GetTrigger()       Trigger 	   { return r.Trigger }
func (r Job) GetLoad()          Load 		   { return r.Load }
func (r Job) GetAlias()      	string 		   { return r.Label }
func (r Job) GetCID() 			link.Cid 	   { return r.CID }

func (r Job) String() string {
	return fmt.Sprintf(`Job:
  Creator:  %s
  Contract: %s
  Trigger:  %s
  Load:     %s
  Label:    %s
  CID:      %s
`,
  r.Creator, r.Contract, r.Trigger, r.Load, r.Label, r.CID)
}

type Jobs []Job

func (d Jobs) String() (out string) {
	for _, job := range d {
		out += job.String() + "\n"
	}
	return strings.TrimSpace(out)
}
func (js Jobs) Sort() {
	sort.Sort(js)
}

func (js Jobs) Len() int { return len(js) }

func (js Jobs) Less(i, j int) bool {
	return js[i].GetLoad().GasPrice > js[j].GetLoad().GasPrice
}

func (js Jobs) Swap(i, j int) { js[i], js[j] = js[j], js[i] }

//______________________________________________________________________

type Trigger struct {
	Period  uint64	`json:"period" yaml:"period"`
	Block   uint64	`json:"block" yaml:"block"`
	Entropy sdk.Dec	`json:"entropy" yaml:"entropy"`
}

func NewTrigger (
	period, block uint64,
	entropy sdk.Dec,
) Trigger {
	return Trigger{
		Period: period,
		Block: block,
		Entropy: entropy,
	}
}

func (r Trigger) GetBlock()   uint64 { return r.Block }
func (r Trigger) GetPeriod()  uint64 { return r.Period }
func (r Trigger) GetEntropy() sdk.Dec { return r.Entropy }

func (r Trigger) String() string {
	return fmt.Sprintf(`Load:
  Block:  %s
  Period: %s
  Entropy: %s
`,
		r.Block, r.Period, r.Entropy)
}

//______________________________________________________________________

type JobStats struct {
	Calls     uint64	`json:"calls" yaml:"calls"`
	Fees      uint64	`json:"fees" yaml:"fees"`
	Gas       uint64	`json:"gas" yaml:"gas"`
	LastBlock uint64	`json:"last_block" yaml:"last_block"`
}

func NewStats (
	calls, fees, gas, block uint64,
) JobStats {
	return JobStats{
		Calls: calls,
		Fees: fees,
		Gas: gas,
		LastBlock: block,
	}
}

func MustMarshalJobStats(cdc *codec.Codec, stats JobStats) []byte {
	return cdc.MustMarshalBinaryLengthPrefixed(stats)
}

func MustUnmarshalJobStats(cdc *codec.Codec, value []byte) JobStats {
	stats, err := UnmarshalJobStats(cdc, value)
	if err != nil {
		panic(err)
	}
	return stats
}

func UnmarshalJobStats(cdc *codec.Codec, value []byte) (stats JobStats, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(value, &stats)
	return stats, err
}

func (r JobStats) GetCalls() uint64 { return r.Calls }
func (r JobStats) GetFees()  uint64 { return r.Fees }
func (r JobStats) GetGas()   uint64 { return r.Gas}
func (r JobStats) GetLastBlock()   uint64 { return r.LastBlock}

func (r JobStats) String() string {
	return fmt.Sprintf(`JobStats:
  Calls:  %s
  Fees: %s
  Gas:  %s
  LastBlock: %s
`,
	r.Calls, r.Fees, r.Gas, r.LastBlock)
}

type JobsStats []JobStats

func (d JobsStats) String() (out string) {
	for _, stat := range d {
		out += stat.String() + "\n"
	}
	return strings.TrimSpace(out)
}


//______________________________________________________________________

type Load struct {
	CallData string	`json:"calldata" yaml:"calldata"`
	GasPrice uint64	`json:"gasprice" yaml:"gasprice"`
}

func NewLoad (
	calldata string,
	gasprice uint64,
) Load {
	return Load{
		CallData: calldata,
		GasPrice: gasprice,
	}
}

func (r Load) GetCallData() string { return r.CallData }
func (r Load) GetGasPrice() uint64 { return r.GasPrice }

func (r Load) String() string {
	return fmt.Sprintf(`Load:
  CallData:  %s
  GasPrice: %s
`,
		r.CallData, r.GasPrice)
}
