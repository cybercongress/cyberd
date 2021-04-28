package rest
// The packages below are commented out at first to prevent an error if this file isn't initially saved.
import (
	// "bytes"
	// "net/http"

	"bytes"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cybercongress/go-cyber/x/cron/types"
	"github.com/cybercongress/go-cyber/x/link"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/cron/add_job",
		postAddJobHandlerFn(cliCtx),
	).Methods("POST")
	r.HandleFunc(
		"/cron/remove_job",
		postRemoveJobHandlerFn(cliCtx),
	).Methods("DELETE")
	r.HandleFunc(
		"/cron/change_job_cid",
		postChangeJobCIDHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_label",
		postChangeJobLabelHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_calldata",
		postChangeJobCallDataHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_gasprice",
		postChangeJoGasPriceHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_period",
		postChangeJobPeriodHandlerFn(cliCtx),
	).Methods("PUT")
	r.HandleFunc(
		"/cron/change_job_block",
		postChangeJobBlockHandlerFn(cliCtx),
	).Methods("PUT")
}

type AddJobReq struct {
	BaseReq rest.BaseReq 		`json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Block 		 	 uint64 	`json:"block" yaml:"block"`
	Period 		 	 uint64 	`json:"period" yaml:"period"`
	CallData 		 string 	`json:"calldata" yaml:"calldata"`
	GasPrice 		 uint64 	`json:"gasprice" yaml:"gasprice"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CID 		 	 string 	`json:"cid" yaml:"cid"`

}

type RemoveJobReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	  `json:"contract" yaml:"contract"`
	Creator  		 string	      `json:"creator" yaml:"creator"`
	Label 		 	 string 	  `json:"label" yaml:"label"`
}

type ChangeJobLabelReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	NewLabel 		 string 	`json:"new_label" yaml:"new_label"`
}

type ChangeJobCIDReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CID 		 	 string 	`json:"cid" yaml:"cid"`
}

type ChangeJobCallDataReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	CallData 		 string 	`json:"calldata" yaml:"calldata"`
}

type ChangeJobGasPriceReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	GasPrice 		 uint64 	`json:"gasprice" yaml:"gasprice"`
}

type ChangeJobPeriodReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	Period 		 	 uint64 	`json:"period" yaml:"period"`
}

type ChangeJobBlockReq struct {
	BaseReq 		 rest.BaseReq `json:"base_req" yaml:"base_req"`
	Contract 		 string 	`json:"contract" yaml:"contract"`
	Creator  		 string	    `json:"creator" yaml:"creator"`
	Label 		 	 string 	`json:"label" yaml:"label"`
	Block 		 	 uint64 	`json:"block" yaml:"block"`
}


func postAddJobHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddJobReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgAddJob(
			cr, co,
			types.NewTrigger(req.Period, req.Block, sdk.NewDec(0)),
			types.NewLoad(req.CallData, req.GasPrice),
			req.Label, link.Cid(req.CID),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postRemoveJobHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveJobReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgRemoveJob(
			cr, co, req.Label,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJobCIDHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobCIDReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeCID(
			cr, co, req.Label, link.Cid(req.CID),
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJobLabelHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobLabelReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeLabel(
			cr, co, req.Label, req.NewLabel,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJobCallDataHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobCallDataReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeCallData(
			cr, co, req.Label, req.CallData,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJoGasPriceHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobGasPriceReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeGasPrice(
			cr, co, req.Label, req.GasPrice,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJobBlockHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobBlockReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeJobBlock(
			cr, co, req.Label, req.Block,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func postChangeJobPeriodHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ChangeJobPeriodReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		cr, err := sdk.AccAddressFromBech32(req.Creator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		co, err := sdk.AccAddressFromBech32(req.Contract)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgChangeJobPeriod(
			cr, co, req.Label, req.Period,
		)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if !bytes.Equal(fromAddr, cr) {
			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own source address")
			return
		}

		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

