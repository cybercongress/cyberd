package keys

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/mintkey"
	"github.com/cybercongress/go-cyber/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/btcd/btcec"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

const hashPrefix = "0x"

func importPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import_private <name>",
		Short: "Import private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			rootDir := viper.GetString(cli.HomeFlag)
			kb, _ := keys.NewKeyBaseFromDir(rootDir)

			name := args[0]

			bufStdin := bufio.NewReader(os.Stdin)
			_, err := kb.Get(name)
			if err == nil {
				// account exists, ask for user confirmation
				fmt.Println("This name already exists. Please choose another one.")
				return nil
			}

			encryptPassword, err := client.GetCheckPassword(
				"Enter a passphrase to encrypt your key to disk:",
				"Repeat the passphrase:", bufStdin)

			privateKey, err := client.GetString(
				"Enter your private key (hex encoded):", bufStdin)
			if err != nil {
				return err
			}

			if util.HasPrefixIgnoreCase(privateKey, hashPrefix) {
				privateKey = privateKey[len(hashPrefix):]
			}

			b, _ := hex.DecodeString(privateKey)

			privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), b)
			ethPubkey, _ := btcec.ParsePubKey(pubKey.SerializeUncompressed(), btcec.S256())

			var cbdPubKey [33]byte
			copy(cbdPubKey[:], ethPubkey.SerializeCompressed()[:])

			var cbdPribKey [32]byte
			copy(cbdPribKey[:], privKey.Serialize()[:])

			pkArmor := mintkey.EncryptArmorPrivKey(secp256k1.PrivKeySecp256k1(cbdPribKey), encryptPassword)

			buf := bytes.NewBuffer(nil)

			buf.Write([]byte{13, 173, 21, 61, 10})
			amino.EncodeString(buf, name)
			buf.Write([]byte{18, 38, 235, 90, 233, 135, 33})
			buf.Write(cbdPubKey[:])
			buf.Write([]byte{26})
			amino.EncodeString(buf, pkArmor)

			bz := buf.Bytes()
			bufRes := bytes.NewBuffer(nil)
			amino.EncodeUvarint(bufRes, uint64(len(bz)))
			bufRes.Write(bz)

			armorStr := mintkey.ArmorInfoBytes(bufRes.Bytes())
			// import armored bytesm
			return kb.Import(name, armorStr)
		},
	}

	return cmd
}
