package cron

import (
	"github.com/cybercongress/go-cyber/x/cron/keeper"
	"github.com/cybercongress/go-cyber/x/cron/types"
)

const (
	ModuleName                = types.ModuleName
	StoreKey				  = types.StoreKey
	DefaultParamspace 	      = types.DefaultParamspace
	QueryParams  			  = types.QueryParams
)

var (
	NewQuerier                = keeper.NewQuerier
	NewKeeper				  = keeper.NewKeeper
	DefaultGenesisState       = types.DefaultGenesisState
)

type (
	Keeper                    = keeper.Keeper
	Job				  	  	  = types.Job
	Jobs				  	  = types.Jobs
	JobStats				  = types.JobStats
	JobsStats				  = types.JobsStats
	Params 					  = types.Params
	GenesisState			  = types.GenesisState
)