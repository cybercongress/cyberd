package bank

import (
	"github.com/cybercongress/go-cyber/x/bank/exported"
	"github.com/cybercongress/go-cyber/x/bank/internal/keeper"
	"github.com/cybercongress/go-cyber/x/bank/internal/types"
)

type (
	Keeper        = exported.Keeper
	IndexedKeeper = exported.IndexedKeeper

	CoinsTransferHook = types.CoinsTransferHook
)

var (
	NewKeeper        = keeper.NewKeeper
	NewIndexedKeeper = keeper.NewIndexedKeeper
)
