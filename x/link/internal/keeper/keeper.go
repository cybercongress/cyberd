package keeper

import (
	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/util"
	"github.com/cybercongress/go-cyber/x/link/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/binary"
	"io"
)

const (
	LinkBytesSize       = uint64(24)
	LinksCountBytesSize = uint64(8)
)

type LinkFilter func(l types.CompactLink) bool

var DefaultLinkFilter = func(l types.CompactLink) bool { return true }

type Keeper struct {
	ms      store.MainKeeper
	key     *sdk.KVStoreKey
	storage store.Storage
}

func NewLinkKeeper(ms store.MainKeeper, key *sdk.KVStoreKey) *Keeper {
	storage, err := store.NewBaseStorage(key.Name(), util.RootifyPath("data/"), LinkBytesSize)
	if err != nil {
		panic("Failed to load links DB")
	}
	return &Keeper{
		key:     key,
		ms:      ms,
		storage: storage,
	}
}

func (lk Keeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	if ctx.BlockHeight() == lk.storage.LastVersion() {
		return
	}
	linkAsBytes := link.MarshalBinary()
	err := lk.storage.Put(linkAsBytes)
	if err != nil {
		panic(err)
	}
	lk.ms.IncrementLinksCount(ctx)
}

func (lk Keeper) GetAllLinks(ctx sdk.Context) (types.Links, types.Links, error) {
	return lk.GetAllLinksFiltered(ctx, DefaultLinkFilter)
}

func (lk Keeper) GetAllLinksFiltered(ctx sdk.Context, filter LinkFilter) (types.Links, types.Links, error) {

	inLinks := make(map[types.CidNumber]types.CidLinks)
	outLinks := make(map[types.CidNumber]types.CidLinks)

	lk.Iterate(ctx, func(link types.CompactLink) {
		if filter(link) {
			types.Links(outLinks).Put(link.From(), link.To(), link.Acc())
			types.Links(inLinks).Put(link.To(), link.From(), link.Acc())
		}
	})

	return inLinks, outLinks, nil
}

func (lk Keeper) GetLinksCount(ctx sdk.Context) uint64 {
	return lk.ms.GetLinksCount(ctx)
}

func (lk Keeper) Iterate(ctx sdk.Context, process func(link types.CompactLink)) {

	err := lk.storage.IterateTillVersion(func(bytes []byte) {
		process(types.UnmarshalBinaryLink(bytes))
	}, ctx.BlockHeight())

	if err != nil {
		panic("Error during links iterating")
	}
}

// write links to writer in binary format: <links_count><cid_number_from><cid_number_to><acc_number>...
func (lk Keeper) WriteLinks(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8) //common bytes array to convert uints

	linksCount := lk.GetLinksCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, linksCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	err = lk.storage.IterateTillVersion(func(bytes []byte) {
		_, _ = writer.Write(bytes)
	}, lk.storage.LastVersion())

	return err
}

func (lk Keeper) Commit(blockHeight uint64) error {
	return lk.storage.Commit(blockHeight)
}
