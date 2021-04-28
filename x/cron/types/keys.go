package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	ModuleName 		= "cron"

	StoreKey 		= ModuleName

  	RouterKey 		= ModuleName

	QuerierRoute 	= ModuleName
)

var (
	JobKey          = []byte{0x00}
	JobStatsKey		= []byte{0x01}
)

func GetJobKey(contract, creator sdk.AccAddress, label string) []byte {
	key := append(append(contract.Bytes(), creator.Bytes()...), []byte(label)...)
	return append(JobKey, key...)
}

func GetJobStatsKey(contract, creator sdk.AccAddress, label string) []byte {
	key := append(append(contract.Bytes(), creator.Bytes()...), []byte(label)...)
	return append(JobStatsKey, key...)
}

