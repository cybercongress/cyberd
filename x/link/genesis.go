package link

import (
	"bufio"
	"os"
	"path/filepath"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cybercongress/go-cyber/util"
)

const (
	LinksFileName       = "config/links"
	LinksExportFileName = "export/links"
)

func InitGenesis(
	ctx sdk.Context, cidNumKeeper CidNumberKeeper, linkIndexedKeeper IndexedKeeper, logger log.Logger,
) (err error) {

	linksFilePath := util.RootifyPath(LinksFileName)
	linksFile, err := os.Open(linksFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("File with links not found. Empty set will be used")
			return nil
		}
		return
	}
	reader := bufio.NewReader(linksFile) // 4096 bytes chunk size

	// initialize slices to read data
	err = cidNumKeeper.LoadFromReader(ctx, reader)
	if err != nil {
		return
	}

	// Read all links
	err = linkIndexedKeeper.LoadFromReader(ctx, reader)
	if err != nil {
		return
	}

	return
}

func WriteGenesis(
	ctx sdk.Context, cidNumKeeper CidNumberKeeper, linkIndexedKeeper IndexedKeeper, logger log.Logger,
) (err error) {

	linksFilePath := util.RootifyPath(LinksExportFileName)
	dirName := filepath.Dir(linksFilePath)
	if _, err := os.Stat(dirName); err != nil {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	}

	linksFile, err := os.Create(linksFilePath)
	if err != nil {
		return
	}

	writer := bufio.NewWriter(linksFile) // 4096 byte chunk
	err = cidNumKeeper.WriteCids(ctx, writer)
	if err != nil {
		return
	}
	err = linkIndexedKeeper.WriteLinks(ctx, writer)
	if err != nil {
		return
	}

	err = writer.Flush()
	if err != nil {
		return
	}
	err = linksFile.Close()

	logger.Info("Cids and links exported. File created.", "path", linksFilePath)
	return
}
