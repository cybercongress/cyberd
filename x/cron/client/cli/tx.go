package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cybercongress/go-cyber/x/cron/types"
	"github.com/cybercongress/go-cyber/x/link"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	cronTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cronTxCmd.AddCommand(flags.PostCommands(
		GetCmdAddJob(cdc),
		GetCmdRemoveJob(cdc),
		GetCmdChangeJobCID(cdc),
		GetCmdChangeJobLabel(cdc),
		GetCmdChangeJobCallData(cdc),
		GetCmdChangeJobGasPrice(cdc),
		GetCmdChangeJobPeriod(cdc),
		GetCmdChangeJobBlock(cdc),
	)...)

	return cronTxCmd
}

func GetCmdAddJob(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-job [contract] [trigger-period] [trigger-block] [load-calldata] [load-gasprice] [label] [cid]",
		Args:  cobra.ExactArgs(7),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron add-job cyber___ 10 0 "{"release":{}}" 1 "label" Qm___ --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			period, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[1])
			}

			block, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block %s not a valid uint, please input a valid block ", args[2])
			}

			gasprice, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("gas price %s not a valid uint, please input a valid block ", args[4])
			}

			jobCID := link.Cid(args[6])

			msg := types.NewMsgAddJob(
				creator, contract,
				types.NewTrigger(period, block, sdk.NewDec(0)),
				types.NewLoad(args[3], gasprice),
				args[5], jobCID,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdRemoveJob(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove-job [contract] [label]",
		Args:  cobra.ExactArgs(2),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron remove-job cyber___ "label" --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRemoveJob(
				creator, contract, args[1],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobCID(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-cid [contract] [label] [cid]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-cid cyber___ "label" Qm___ --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgChangeCID(
				creator, contract, args[1], link.Cid(args[2]),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobLabel(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-label [contract] [label] [newLabel]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-label cyber___ "label" "newLabel" --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgChangeLabel(
				creator, contract, args[1], args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobCallData(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-calldata [contract] [label] [calldata]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-calldata cyber___ "label" "{"release":{}}" --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgChangeCallData(
				creator, contract, args[1], args[2],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobGasPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-gasprice [contract] [label] [gasprice]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-gasprice cyber___ "label" 10 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			gasprice, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("gas price %s not a valid uint, please input a valid block ", args[4])
			}

			msg := types.NewMsgChangeGasPrice(
				creator, contract, args[1], gasprice,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobPeriod(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-period [contract] [label] [period]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-period cyber___ "label" 10 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			period, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block period %s not a valid uint, please input a valid block period", args[2])
			}

			msg := types.NewMsgChangeJobPeriod(
				creator, contract, args[1], period,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdChangeJobBlock(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "change-job-block [contract] [label] [block]",
		Args:  cobra.ExactArgs(3),
		Short: "Short",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Long.

Example:
$ %s tx cron change-job-block cyber___ "label" 10 --from mykey
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(auth.DefaultTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)

			creator := cliCtx.GetFromAddress()
			contract, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			block, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("block %s not a valid uint, please input a valid block ", args[2])
			}

			msg := types.NewMsgChangeJobBlock(
				creator, contract, args[1], block,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}