package cli

import (
	"fmt"

	"github.com/cybercongress/go-cyber/x/rank/internal/types"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands for the minting module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	rankingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the rank module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rankingQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryParams(cdc),
			GetCmdQueryCalculationWindow(cdc),
			GetCmdQueryDampingFactor(cdc),
			GetCmdQueryTolerance(cdc),
		)...,
	)

	return rankingQueryCmd
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current rank parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cdc.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryCalculationWindow implements a command to return the current rank
// calculation window.
func GetCmdQueryCalculationWindow(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "window",
		Short: "Query the current rank calculation window",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCalculationWindow)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var window sdk.Int
			if err := cdc.UnmarshalJSON(res, &window); err != nil {
				return err
			}

			return cliCtx.PrintOutput(window)
		},
	}
}

// GetCmdQueryDampingFactor implements a command to return the current rank
// damping factor value.
func GetCmdQueryDampingFactor(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "damping-factor",
		Short: "Query the current rank damping factor value",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDampingFactor)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var factor sdk.Dec
			if err := cdc.UnmarshalJSON(res, &factor); err != nil {
				return err
			}

			return cliCtx.PrintOutput(factor)
		},
	}
}

// GetCmdQueryDampingFactor implements a command to return the current rank
// damping factor value.
func GetCmdQueryTolerance(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tolerance",
		Short: "Query the current rank tolerance",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTolerance)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var tolerance sdk.Dec
			if err := cdc.UnmarshalJSON(res, &tolerance); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tolerance)
		},
	}
}
