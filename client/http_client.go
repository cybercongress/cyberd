package client

import (
	"errors"
	"fmt"
	"os"

	cli "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cskeys "github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	tdmClient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/lib/client"

	"github.com/cybercongress/go-cyber/app"
	"github.com/cybercongress/go-cyber/cmd/cyberd/rpc"
	"github.com/cybercongress/go-cyber/x/bandwidth"
	"github.com/cybercongress/go-cyber/x/link"
)

var _ CyberdClient = &HttpCyberdClient{}

type HttpCyberdClient struct {
	// tdm client
	tdmClient tdmClient.Client
	// transport client
	httpClient rpcclient.JSONRPCClient

	// general fields
	nodeUrl string
	chainId string

	// fields used by local keys store to sign transactions
	passphrase  string
	fromAddress sdk.AccAddress
	cliCtx      *cli.CLIContext
	txBuilder   *authtypes.TxBuilder
}

func NewHttpCyberdClient(nodeUrl string, passphrase string, signAddr string) *HttpCyberdClient {
	tdmHttpClient := tdmClient.NewHTTP(nodeUrl, "/websocket")
	httpClient := rpcclient.NewJSONRPCClient(nodeUrl)
	status, err := tdmHttpClient.Status()
	if err != nil {
		panic(err)
	}

	cdc := app.MakeCodec()
	app.SetPrefix()
	addr, cliAddrName := accountFromAddress(signAddr)
	verifier := &NoopVerifier{ChainId: status.NodeInfo.Network}
	cliCtx := cli.CLIContext{
		Client:    tdmHttpClient,
		NodeURI:   nodeUrl,
		From:      cliAddrName,
		TrustNode: true,
		Verifier:  verifier,
	}.WithCodec(cdc)

	accountNumber, accountSequence, err := authtypes.NewAccountRetriever(cliCtx).GetAccountNumberSequence(addr)

	if err != nil {
		panic(err)
	}

	txBuilder := authtypes.NewTxBuilder(
		utils.GetTxEncoder(cdc), accountNumber, accountSequence, 111111, 1.0, false, status.NodeInfo.Network,
		"", sdk.Coins{}, sdk.NewDecCoins(sdk.Coins{}),
	).WithTxEncoder(utils.GetTxEncoder(cdc))

	return &HttpCyberdClient{
		tdmClient:  tdmHttpClient,
		httpClient: *httpClient,

		chainId: status.NodeInfo.Network,

		passphrase:  passphrase,
		fromAddress: addr,
		cliCtx:      &cliCtx,
		txBuilder:   &txBuilder,
	}
}

func (c *HttpCyberdClient) GetChainId() string {
	return c.chainId
}

func (c *HttpCyberdClient) IsLinkExist(from link.Cid, to link.Cid, addr sdk.AccAddress) (result bool, err error) {

	_, err = c.httpClient.Call("is_link_exist",
		map[string]interface{}{"from": from, "to": to, "address": addr.String()},
		&result,
	)

	return
}

func (c *HttpCyberdClient) IsAnyLinkExist(from link.Cid, to link.Cid) (result bool, err error) {
	_, err = c.httpClient.Call("is_link_exist",
		map[string]interface{}{"from": from, "to": to},
		&result,
	)
	return
}

func (c *HttpCyberdClient) GetCurrentBandwidthCreditPrice() (float64, error) {
	result := &rpc.ResultBandwidthPrice{}
	_, err := c.httpClient.Call("current_bandwidth_price", map[string]interface{}{}, &result)
	return result.Price, err
}

func (c *HttpCyberdClient) GetAccountBandwidth() (result bandwidth.AcсBandwidth, err error) {
	_, err = c.httpClient.Call("account_bandwidth",
		map[string]interface{}{"address": c.fromAddress.String()}, &result)
	return
}

func (c *HttpCyberdClient) SubmitLinkSync(links link.Link) error {
	return c.SubmitLinksSync([]link.Link{links}, false)
}

func (c *HttpCyberdClient) SubmitLinksSync(links []link.Link, submitOnlyNew bool) error {

	// used to remove duplicated items
	var filter = make(link.CidsFilter)
	msges := make([]sdk.Msg, 0)

	for _, l := range links {

		if filter.Contains(l.From, l.To) {
			continue
		}

		var err error
		var exists bool
		if submitOnlyNew {
			exists, err = c.IsAnyLinkExist(l.From, l.To)
		} else {
			exists, err = c.IsLinkExist(l.From, l.To, c.fromAddress)
		}

		if err != nil {
			return err
		}
		if !exists {
			msges = append(msges, link.NewMsg(c.fromAddress, []link.Link{{From: l.From, To: l.To}}))
		}
		filter.Put(l.From, l.To)
	}
	return c.BroadcastTx(msges)
}

func (c *HttpCyberdClient) BroadcastTx(msgs []sdk.Msg) error {

	if len(msgs) == 0 {
		return nil
	}

	txBytes, err := c.txBuilder.BuildAndSign(c.cliCtx.From, c.passphrase, msgs)
	if err != nil {
		panic(err)
	}

	result, err := c.cliCtx.BroadcastTxSync(txBytes)
	if err != nil {
		println("Error during broadcasting tx. Rebrodcasting ...")
		println(err.Error())

		_ = c.BroadcastTx(msgs)
	}

	if result.Code != 0 {
		return errors.New(string(result.Logs.String()))
	}
	newBuilder := c.txBuilder.WithSequence(c.txBuilder.Sequence() + 1)
	c.txBuilder = &newBuilder
	return nil
}

func accountFromAddress(from string) (fromAddr sdk.AccAddress, fromName string) {
	if from == "" {
		return nil, ""
	}

	keybase, err := keys.NewKeyBaseFromHomeFlag()
	if err != nil {
		fmt.Println("no keybase found")
		os.Exit(1)
	}

	var info cskeys.Info
	if addr, err := sdk.AccAddressFromBech32(from); err == nil {
		info, err = keybase.GetByAddress(addr)
		if err != nil {
			fmt.Printf("could not find key %s\n", from)
			os.Exit(1)
		}
	} else {
		info, err = keybase.Get(from)
		if err != nil {
			fmt.Printf("could not find key %s\n", from)
			os.Exit(1)
		}
	}

	fromAddr = info.GetAddress()
	fromName = info.GetName()
	return
}
