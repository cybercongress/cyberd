package commands

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	cbd "github.com/cybercongress/go-cyber/types"
	"github.com/cybercongress/go-cyber/x/link"

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

			txCtx := authtypes.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc)

			cidFrom := link.Cid(viper.GetString(flagCidFrom))
			cidTo := link.Cid(viper.GetString(flagCidTo))

			if _, err := cid.Decode(string(cidFrom)); err != nil {
				return cbd.ErrInvalidCid()
			}

			if _, err := cid.Decode(string(cidTo)); err != nil {
				return cbd.ErrInvalidCid()
			}

			signAddr := cliCtx.GetFromAddress()

			// ensure that account exists in chain
			_, err := authtypes.NewAccountRetriever(cliCtx).GetAccount(signAddr)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := link.NewMsg(signAddr, []link.Link{{From: cidFrom, To: cidTo}})

			return utils.CompleteAndBroadcastTxCLI(txCtx, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagCidFrom, "", "Content id to link from")
	cmd.Flags().String(flagCidTo, "", "Content id to link to")

	return cmd
}
