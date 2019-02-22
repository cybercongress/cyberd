package commands

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	cbd "github.com/cybercongress/cyberd/types"
	"github.com/cybercongress/cyberd/x/link"
	cbdlink "github.com/cybercongress/cyberd/x/link/types"
	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCidFrom = "cid-from"
	flagCidTo   = "cid-to"
)

// LinkTxCmd will create a link tx and sign it with the given key.
func LinkTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Create and sign a link tx",
		RunE: func(cmd *cobra.Command, args []string) error {

			txCtx := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			cidFrom := cbdlink.Cid(viper.GetString(flagCidFrom))
			cidTo := cbdlink.Cid(viper.GetString(flagCidTo))

			if _, err := cid.Decode(string(cidFrom)); err != nil {
				return cbd.ErrInvalidCid()
			}

			if _, err := cid.Decode(string(cidTo)); err != nil {
				return cbd.ErrInvalidCid()
			}

			from := cliCtx.GetFromAddress()

			// ensure that account exists in chain
			_, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := link.NewMsg(from, []cbdlink.Link{{From: cidFrom, To: cidTo}})

			return utils.CompleteAndBroadcastTxCLI(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCidFrom, "", "Content id to link from")
	cmd.Flags().String(flagCidTo, "", "Content id to link to")

	return cmd
}
