package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link"
	"github.com/cybercongress/go-cyber/x/link/types"
)

type StakeKeeper interface {
	DetectUsersStakeChange(ctx sdk.Context) bool
	GetTotalStakes() map[uint64]uint64
}

type GraphIndexedKeeper interface {
	FixLinks()
	EndBlocker() bool

	GetOutLinks() link.Links
	GetInLinks() link.Links

	GetLinksCount(sdk.Context) uint64
	GetCurrentBlockNewLinks() []link.CompactLink
	GetCidsCount(sdk.Context) uint64
}

type GraphKeeper interface {
	GetCidsCount(sdk.Context) uint64
	GetCidNumber(ctx sdk.Context, cid types.Cid) (types.CidNumber, bool)
	GetCid(ctx sdk.Context, num types.CidNumber) types.Cid
}
