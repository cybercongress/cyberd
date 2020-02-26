package keeper

import (
	"github.com/cybercongress/go-cyber/merkle"
	"github.com/cybercongress/go-cyber/util"
	"github.com/cybercongress/go-cyber/x/link/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cmn "github.com/tendermint/tendermint/libs/common"

	"crypto/sha256"
	"encoding/binary"
	"io"
)

type IndexedKeeper struct {
	*Keeper

	// Actual links for current rank calculated state.
	currentRankInLinks  types.Links
	currentRankOutLinks types.Links

	// New links for the next rank calculation.
	nextRankInLinks  types.Links
	nextRankOutLinks types.Links

	currentBlockLinks []types.CompactLink
	MerkleTree        *merkle.Tree
}

func NewIndexedKeeper(keeper *Keeper) *IndexedKeeper {
	merkleTree := merkle.NewTree(sha256.New(), true)

	return &IndexedKeeper{Keeper: keeper, MerkleTree: merkleTree}
}

func (i *IndexedKeeper) Load(rankCtx sdk.Context, freshCtx sdk.Context) {
	inLinks, outLinks, err := i.Keeper.GetAllLinks(rankCtx)
	if err != nil {
		cmn.Exit(err.Error())
	}

	i.currentRankInLinks = inLinks
	i.currentRankOutLinks = outLinks

	newInLinks, newOutLinks, err := i.Keeper.GetAllLinksFiltered(freshCtx, func(l types.CompactLink) bool {
		return !i.currentRankOutLinks.IsLinkExist(l.From(), l.To(), l.Acc())
	})

	if err != nil {
		cmn.Exit(err.Error())
	}

	i.nextRankInLinks = newInLinks
	i.nextRankOutLinks = newOutLinks

	i.Iterate(freshCtx, func(link types.CompactLink) {
		linkAsBytes := link.MarshalBinary()
		i.MerkleTree.Push(linkAsBytes)
	})
}

func (i *IndexedKeeper) FixLinks() {
	i.currentRankInLinks.PutAll(i.nextRankInLinks)
	i.currentRankOutLinks.PutAll(i.nextRankOutLinks)
	i.nextRankInLinks = make(types.Links)
	i.nextRankOutLinks = make(types.Links)
}

// return true if this block has new links
func (i *IndexedKeeper) EndBlocker() bool {
	hasNewLinks := len(i.currentBlockLinks) > 0
	for _, link := range i.currentBlockLinks {
		i.nextRankOutLinks.Put(link.From(), link.To(), link.Acc())
		i.nextRankInLinks.Put(link.To(), link.From(), link.Acc())
	}
	i.currentBlockLinks = make([]types.CompactLink, 0, 1000) // todo: 1000 hardcoded value
	return hasNewLinks
}

func (i *IndexedKeeper) PutIntoIndex(link types.CompactLink) {
	i.currentBlockLinks = append(i.currentBlockLinks, link)
}

func (i *IndexedKeeper) PutLink(ctx sdk.Context, link types.CompactLink) {
	if !ctx.IsCheckTx() {
		i.currentBlockLinks = append(i.currentBlockLinks, link)
	}

	linkAsBytes := link.MarshalBinary()
	i.MerkleTree.Push(linkAsBytes)

	i.Keeper.PutLink(ctx, link)
}

func (i *IndexedKeeper) GetOutLinks() types.Links {
	return i.currentRankOutLinks
}

func (i *IndexedKeeper) GetInLinks() types.Links {
	return i.currentRankInLinks
}

func (i *IndexedKeeper) GetNextOutLinks() types.Links {
	return i.nextRankOutLinks
}

func (i *IndexedKeeper) GetCurrentBlockLinks() []types.CompactLink {
	return i.currentBlockLinks
}

func (i *IndexedKeeper) GetCurrentBlockNewLinks() []types.CompactLink {
	result := make([]types.CompactLink, 0, len(i.currentBlockLinks))
	for _, link := range i.currentBlockLinks {
		if !i.IsAnyLinkExist(link.From(), link.To()) {
			result = append(result, link)
		}
	}
	return result
}

func (i *IndexedKeeper) GetNetworkLinkHash() []byte {
	return i.MerkleTree.RootHash()
}

func (i *IndexedKeeper) IsAnyLinkExist(from types.CidNumber, to types.CidNumber) bool {
	return i.currentRankOutLinks.IsAnyLinkExist(from, to) || i.nextRankOutLinks.IsAnyLinkExist(from, to)
}

func (i *IndexedKeeper) IsLinkExist(link types.CompactLink) bool {
	return i.currentRankOutLinks.IsLinkExist(link.From(), link.To(), link.Acc()) ||
		i.nextRankOutLinks.IsLinkExist(link.From(), link.To(), link.Acc())
}

//todo: remove duplicated method (BaseLinksKeeper)
func (i *IndexedKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	linksCountBytes, err := util.ReadExactlyNBytes(reader, LinksCountBytesSize)
	if err != nil {
		return
	}
	linksCount := binary.LittleEndian.Uint64(linksCountBytes)

	for j := uint64(0); j < linksCount; j++ {
		linkBytes, err := util.ReadExactlyNBytes(reader, LinkBytesSize)
		if err != nil {
			return err
		}
		i.PutLink(ctx, types.UnmarshalBinaryLink(linkBytes))
	}
	return
}
