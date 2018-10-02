package storage

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/cosmos/poc/app/coin"
)

// returns all added cids
func GetAllAccountsStakes(ctx sdk.Context, am auth.AccountMapper) map[AccountNumber]int64 {

	stakes := make(map[AccountNumber]int64)

	collect := func(acc auth.Account) bool {
		balance := acc.GetCoins().AmountOf(coin.CBD).Int64()
		stakes[AccountNumber(acc.GetAddress().String())] = balance
		return false
	}

	am.IterateAccounts(ctx, collect)
	return stakes
}
