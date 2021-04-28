package link

import (
	"github.com/cybercongress/go-cyber/x/link/keeper"
	"github.com/cybercongress/go-cyber/x/link/types"
)

const (
	ModuleName 		= types.ModuleName
	StoreKey	    = types.StoreKey
	QuerierRoute	= types.QuerierRoute
	RouterKey  		= types.RouterKey
)

var (
	NewKeeper          = keeper.NewKeeper
	NewIndexKeeper     = keeper.NewIndexKeeper
	NewQuerier         = keeper.NewQuerier

	NewMsgCyberlink	   = types.NewMsgCyberlink
	RegisterCodec      = types.RegisterCodec
	NewLink       	   = types.NewLink

	ErrCidNotFound 	   = types.ErrCidNotFound
	ErrInvalidCid	   = types.ErrInvalidCid
)

type (
	GraphKeeper     = keeper.GraphKeeper
	IndexKeeper   	= keeper.IndexKeeper

	MsgCyberlink    = types.MsgCyberlink
	//AccNumber   	= types.AccNumber
	Cid				= types.Cid
	CidNumber		= types.CidNumber
	CidLinks    	= types.CidLinks
	CidsFilter  	= types.CidsFilter
	Link        	= types.Link
	Links       	= types.Links
	CompactLink 	= types.CompactLink
	LinkFilter  	= types.LinkFilter
)
