package rest

import (
	"fmt"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/cybercongress/go-cyber/x/cron/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cron/params",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/cron/{%s}/{%s}/{%s}/job", Contract, Creator, Label),
		queryJobHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/cron/{%s}/{%s}/{%s}/job_stats", Contract, Creator, Label),
		queryJobStatsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/cron/jobs",
		queryJobsHandlerFn(cliCtx),
	).Methods("GET")

}

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryJobHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		contract  := vars[Contract]
		creator  := vars[Creator]
		label  := vars[Label]

		co, err := sdk.AccAddressFromBech32(contract)
		cr, err := sdk.AccAddressFromBech32(creator)
		params := types.NewQueryJobParams(co, cr, label)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryJob), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryJobStatsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		contract  := vars[Contract]
		creator  := vars[Creator]
		label  := vars[Label]

		co, err := sdk.AccAddressFromBech32(contract)
		cr, err := sdk.AccAddressFromBech32(creator)
		params := types.NewQueryJobParams(co, cr, label)
		bz, err := cliCtx.Codec.MarshalJSON(params)

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryJobStats), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryJobsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryJobs)

		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}