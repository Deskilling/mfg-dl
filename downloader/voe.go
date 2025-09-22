package downloader

import (
	parserHoster "mfg-dl/parser/hoster"
	"mfg-dl/request"
	"mfg-dl/util"

	"github.com/charmbracelet/log"
)

func Voe(voeUrl string) error {
	voeHtml, err := request.Get(voeUrl)
	if err != nil {
		log.Error(err)
		return err
	}

	parsed, err := parserHoster.Voe(voeHtml)
	if err != nil {
		log.Error(err)
		return err
	}

	// make sure its chill
	parsed.Directory = util.RemoveAfterSymbol(parsed.FileCode, "/")

	err = request.DownloadFile(parsed.Source, "./temp/"+parsed.Directory+"/master.m3u8")
	if err != nil {
		log.Error(err)
		return err
	}

	// read content -> parse into index downloder

	return nil
}
