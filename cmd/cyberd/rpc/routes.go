package rpc

import (
	"github.com/cybercongress/go-cyber/app"
	"github.com/tendermint/tendermint/rpc/core"
	"github.com/tendermint/tendermint/rpc/lib/server"
)

var cyberdApp *app.CyberdApp
var codec = app.MakeCodec()

func SetCyberdApp(cApp *app.CyberdApp) {
	cyberdApp = cApp
}

var Routes = map[string]*rpcserver.RPCFunc{
	"search":                  rpcserver.NewRPCFunc(Search, "cid,page,perPage"),
	"top":                     rpcserver.NewRPCFunc(Top, "page,perPage"),
	"rank":                    rpcserver.NewRPCFunc(Rank, "cid,proof"),
	"account":                 rpcserver.NewRPCFunc(Account, "address"),
	"account_bandwidth":       rpcserver.NewRPCFunc(AccountBandwidth, "address"),
	"is_link_exist":           rpcserver.NewRPCFunc(IsLinkExist, "from,to,address"),
	"current_bandwidth_price": rpcserver.NewRPCFunc(CurrentBandwidthPrice, ""),
	"index_stats":             rpcserver.NewRPCFunc(IndexStats, ""),

	"staking/validators": rpcserver.NewRPCFunc(StakingValidators, "page,limit,status"),
	"staking/pool":       rpcserver.NewRPCFunc(StakingPool, ""),

	// routes required cuz of amino.js not exists yet.
	"submit_signed_link": rpcserver.NewRPCFunc(SignedMsgHandler(UnmarshalLinkRequestFn), "data"),
	"submit_signed_send": rpcserver.NewRPCFunc(SignedMsgHandler(UnmarshalSendRequestFn), "data"),
}

func init() {
	for path, rpcFunc := range Routes {
		core.Routes[path] = rpcFunc
	}
}
