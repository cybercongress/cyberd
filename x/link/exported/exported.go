package exported

import (
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cybercongress/go-cyber/x/link/internal/keeper"
	"github.com/cybercongress/go-cyber/x/link/internal/types"
)

type KeeperI interface {
	GetAllLinks(sdk.Context) (types.Links, types.Links, error)
	GetAllLinksFiltered(sdk.Context, keeper.LinkFilter) (types.Links, types.Links, error)

	GetLinksCount(sdk.Context) uint64
	Iterate(sdk.Context, func(types.CompactLink))

	PutLink(sdk.Context, types.CompactLink)
	WriteLinks(sdk.Context, io.Writer) error

	Commit(blockHeight uint64) error
}

type IndexedKeeperI interface {
	KeeperI

	Load(rankCtx sdk.Context, freshCtx sdk.Context)
	FixLinks()
	EndBlocker() bool

	PutIntoIndex(types.CompactLink)

	GetOutLinks() types.Links
	GetInLinks() types.Links

	GetNextOutLinks() types.Links

	GetCurrentBlockLinks() []types.CompactLink
	GetCurrentBlockNewLinks() []types.CompactLink

	GetNetworkLinkHash() []byte

	IsAnyLinkExist(from types.CidNumber, to types.CidNumber) bool
	IsLinkExist(types.CompactLink) bool

	LoadFromReader(sdk.Context, io.Reader) error
}

type CidNumberKeeperI interface {
	GetCidNumber(sdk.Context, types.Cid) (types.CidNumber, bool)
	GetCid(sdk.Context, types.CidNumber) types.Cid

	GetOrPutCidNumber(sdk.Context, types.Cid) types.CidNumber
	GetFullCidsNumbers(sdk.Context) map[types.Cid]types.CidNumber

	GetCidsCount(sdk.Context) uint64
	Iterate(sdk.Context, func(types.Cid, types.CidNumber))

	PutCid(sdk.Context, types.Cid, types.CidNumber)
	WriteCids(sdk.Context, io.Writer) error

	LoadFromReader(sdk.Context, io.Reader) error
}
