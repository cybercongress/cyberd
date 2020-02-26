package keeper

import (
	"github.com/cybercongress/go-cyber/store"
	"github.com/cybercongress/go-cyber/util"
	"github.com/cybercongress/go-cyber/x/link/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/binary"
	"errors"
	"io"
)

const (
	CidLengthBytesSize = uint64(1)
	CidNumberBytesSize = uint64(8)
	CidCountBytesSize  = uint64(8)
)

type CidNumberKeeper struct {
	ms         store.MainKeeper
	key        *sdk.KVStoreKey
	reverseKey *sdk.KVStoreKey
}

func NewCidNumberKeeper(ms store.MainKeeper, key *sdk.KVStoreKey, reverseKey *sdk.KVStoreKey) *CidNumberKeeper {
	return &CidNumberKeeper{
		ms:         ms,
		key:        key,
		reverseKey: reverseKey,
	}
}

// Return cid number and true, if cid exists
func (k CidNumberKeeper) GetCidNumber(ctx sdk.Context, cid types.Cid) (types.CidNumber, bool) {
	cidsIndex := ctx.KVStore(k.key)
	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)
	if cidIndexAsBytes != nil {
		return types.CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes)), true
	}
	return 0, false
}

func (k CidNumberKeeper) GetCid(ctx sdk.Context, num types.CidNumber) types.Cid {
	index := ctx.KVStore(k.reverseKey)
	cidNumberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cidNumberAsBytes, uint64(num))
	cidAsBytes := index.Get(cidNumberAsBytes)
	return types.Cid(cidAsBytes)
}

// WARNING: use only for state import. Don't forget to set right cid count after
func (k CidNumberKeeper) PutCid(ctx sdk.Context, cid types.Cid, cidNumber types.CidNumber) {
	cidsIndex := ctx.KVStore(k.key)
	cidsReverseIndex := ctx.KVStore(k.reverseKey)

	cidNumberAsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(cidNumberAsBytes, uint64(cidNumber))

	cidsIndex.Set([]byte(cid), cidNumberAsBytes)
	cidsReverseIndex.Set(cidNumberAsBytes, []byte(cid))
}

// CIDs index is array of all added CIDs, sorted asc by first link time.
//   - for given link, CIDs added in order [CID1, CID2] (if they both new to chain)
// This method performs lookup of CIDs, returns index value, or create and put in index new value if not exists.
func (k CidNumberKeeper) GetOrPutCidNumber(ctx sdk.Context, cid types.Cid) types.CidNumber {

	cidsIndex := ctx.KVStore(k.key)
	cidsReverseIndex := ctx.KVStore(k.reverseKey)

	cidAsBytes := []byte(cid)
	cidIndexAsBytes := cidsIndex.Get(cidAsBytes)

	// new cid, get new index
	if cidIndexAsBytes == nil {

		lastIndex := k.GetCidsCount(ctx)
		lastIndexAsBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(lastIndexAsBytes, lastIndex)

		cidsIndex.Set(cidAsBytes, lastIndexAsBytes)
		cidsReverseIndex.Set(lastIndexAsBytes, cidAsBytes)
		k.ms.SetLastCidIndex(ctx, lastIndexAsBytes)
		return types.CidNumber(lastIndex)
	}

	return types.CidNumber(binary.LittleEndian.Uint64(cidIndexAsBytes))
}

// returns all added cids
func (k CidNumberKeeper) GetFullCidsNumbers(ctx sdk.Context) map[types.Cid]types.CidNumber {
	index := make(map[types.Cid]types.CidNumber)
	k.Iterate(ctx, func(cid types.Cid, number types.CidNumber) {
		index[cid] = number
	})
	return index
}

func (k CidNumberKeeper) Iterate(ctx sdk.Context, process func(types.Cid, types.CidNumber)) {
	iterator := ctx.KVStore(k.key).Iterator(nil, nil)
	defer iterator.Close()

	for iterator.Valid() {
		process(types.Cid(iterator.Key()), types.CidNumber(binary.LittleEndian.Uint64(iterator.Value())))
		iterator.Next()
	}
}

func (k CidNumberKeeper) GetCidsCount(ctx sdk.Context) uint64 {
	return k.ms.GetCidsCount(ctx)
}

// write CIDs to writer in binary format: <n><cid1_size><cid1><cid1_number><cid2_size><cid2><cid2_number>....<cidn_size><cidn><cidn_number>
func (k CidNumberKeeper) WriteCids(ctx sdk.Context, writer io.Writer) (err error) {
	uintAsBytes := make([]byte, 8) //common bytes array to convert uints

	cidsCount := k.GetCidsCount(ctx)
	binary.LittleEndian.PutUint64(uintAsBytes, cidsCount)
	_, err = writer.Write(uintAsBytes)
	if err != nil {
		return
	}

	k.Iterate(ctx, func(cid types.Cid, number types.CidNumber) {
		cidLength := len(cid)
		if cidLength > 255 {
			err = errors.New("cid length cannot be over 255")
			return
		}

		_, err = writer.Write([]byte{byte(cidLength)})
		if err != nil {
			return
		}
		_, err = writer.Write([]byte(cid))
		if err != nil {
			return
		}
		binary.LittleEndian.PutUint64(uintAsBytes, uint64(number))
		_, err = writer.Write(uintAsBytes)
		if err != nil {
			return
		}
	})
	return
}

func (k CidNumberKeeper) LoadFromReader(ctx sdk.Context, reader io.Reader) (err error) {
	cidCountBytes, err := util.ReadExactlyNBytes(reader, CidCountBytesSize)
	if err != nil {
		return
	}
	cidCount := binary.LittleEndian.Uint64(cidCountBytes)

	// Read all CIDs with their numbers
	for i := uint64(0); i < cidCount; i++ {
		cidLengthBytes, err := util.ReadExactlyNBytes(reader, CidLengthBytesSize)
		if err != nil {
			return err
		}

		cidBytes, err := util.ReadExactlyNBytes(reader, uint64(cidLengthBytes[0]))
		if err != nil {
			return err
		}
		cid := types.Cid(cidBytes)
		cidNumberBytes, err := util.ReadExactlyNBytes(reader, CidNumberBytesSize)
		if err != nil {
			return err
		}
		cidNumber := types.CidNumber(binary.LittleEndian.Uint64(cidNumberBytes))
		k.PutCid(ctx, cid, cidNumber)
	}

	lastCidIndex := cidCount - 1
	lastCidIndexBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lastCidIndexBytes, lastCidIndex)
	k.ms.SetLastCidIndex(ctx, lastCidIndexBytes)

	return
}
