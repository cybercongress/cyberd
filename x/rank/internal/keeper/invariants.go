package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cybercongress/go-cyber/x/rank/exported"
	"github.com/cybercongress/go-cyber/x/rank/internal/types"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k exported.StateKeeper) {
	ir.RegisterRoute(types.ModuleName, "index-error",
		IndexErrorInvariant(k))
}

func IndexErrorInvariant(keeper exported.StateKeeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var broken bool

		if err := keeper.GetIndexError(); err != nil {
			msg = fmt.Sprintf("index error: %s", err.Error())
			broken = true
		}

		return sdk.FormatInvariant(types.ModuleName, "index error", msg), broken
	}
}
