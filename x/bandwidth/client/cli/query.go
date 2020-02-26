package cli

import (
	"fmt"

	"github.com/cybercongress/go-cyber/x/bandwidth/internal/types"
	"github.com/spf13/cobra"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the cli query commands for the bandwidth module.
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	rankingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bandwidth module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	rankingQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryParams(cdc),
			GetCmdQueryDesirableBandwidth(cdc),
			GetCmdQueryMaxBlockBandwidth(cdc),
			GetCmdQueryRecoveryPeriod(cdc),
			GetCmdQueryAdjustPricePeriod(cdc),
			GetCmdQueryBaseCreditPrice(cdc),
			GetCmdQueryTxCost(cdc),
			GetCmdQueryLinkMsgCost(cdc),
			GetCmdQueryNonLinkMsgCost(cdc),
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

// GetCmdQueryDesirableBandwidth implements a command to return the current
// desirable bandwidth of network.
func GetCmdQueryDesirableBandwidth(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "window",
		Short: "Query the current desirable bandwidth",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDesirableBandwidth)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var band sdk.Int
			if err := cdc.UnmarshalJSON(res, &band); err != nil {
				return err
			}

			return cliCtx.PrintOutput(band)
		},
	}
}

// GetCmdQueryMaxBlockBandwidth implements a command to return the current max
// block bandwidth value.
func GetCmdQueryMaxBlockBandwidth(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "max-block",
		Short: "Query the current max block bandwidth",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMaxBlockBandwidth)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var factor sdk.Int
			if err := cdc.UnmarshalJSON(res, &factor); err != nil {
				return err
			}

			return cliCtx.PrintOutput(factor)
		},
	}
}

// GetCmdQueryRecoveryPeriod implements a command to return the current bandwidth
// recovery period.
func GetCmdQueryRecoveryPeriod(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "recovery-period",
		Short: "Query the current bandwidth recovery period",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecoveryPeriod)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var period sdk.Int
			if err := cdc.UnmarshalJSON(res, &period); err != nil {
				return err
			}

			return cliCtx.PrintOutput(period)
		},
	}
}

// GetCmdQueryAdjustPricePeriod implements a command to return the current bandwidth
// adjust price period.
func GetCmdQueryAdjustPricePeriod(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "price-period",
		Short: "Query the current bandwidth adjust price period",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAdjustPricePeriod)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var period sdk.Int
			if err := cdc.UnmarshalJSON(res, &period); err != nil {
				return err
			}

			return cliCtx.PrintOutput(period)
		},
	}
}

// GetCmdQueryBaseCreditPrice implements a command to return the current bandwidth
// base credit price.
func GetCmdQueryBaseCreditPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "credit-price",
		Short: "Query the current bandwidth base credit price",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBaseCreditPrice)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var price sdk.Dec
			if err := cdc.UnmarshalJSON(res, &price); err != nil {
				return err
			}

			return cliCtx.PrintOutput(price)
		},
	}
}

// GetCmdQueryTxCost implements a command to return the current bandwidth
// cost of Tx.
func GetCmdQueryTxCost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "tx-cost",
		Short: "Query the current bandwidth cost of Tx",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTxCost)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var cost sdk.Int
			if err := cdc.UnmarshalJSON(res, &cost); err != nil {
				return err
			}

			return cliCtx.PrintOutput(cost)
		},
	}
}


// GetCmdQueryLinkMsgCost implements a command to return the current bandwidth
// cost of link msg.
func GetCmdQueryLinkMsgCost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "link-cost",
		Short: "Query the current bandwidth cost of Link Msg",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLinkMsgCost)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var cost sdk.Int
			if err := cdc.UnmarshalJSON(res, &cost); err != nil {
				return err
			}

			return cliCtx.PrintOutput(cost)
		},
	}
}

// GetCmdQueryNonLinkMsgCost implements a command to return the current bandwidth
// cost of non-link Msg.
func GetCmdQueryNonLinkMsgCost(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "non-link-cost",
		Short: "Query the current bandwidth cost of non-Link Msg",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNonLinkMsgCost)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var cost sdk.Int
			if err := cdc.UnmarshalJSON(res, &cost); err != nil {
				return err
			}

			return cliCtx.PrintOutput(cost)
		},
	}
}