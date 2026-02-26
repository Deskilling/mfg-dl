package anime4k

import (
	"encoding/json"
	"fmt"
	"strings"

	"mfg-dl/internal/core"
	"mfg-dl/internal/request"
	"mfg-dl/pkg/filesystem"

	"github.com/charmbracelet/log"
)

const latestURL = "https://api.github.com/repos/bloc97/Anime4K/releases/latest"

type releaseResponse struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

type asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

func DownloadLatestRelease() (err error) {
	if !core.GetConfig().Shader.Enable {
		return nil
	}

	if !core.GetConfig().Shader.Autoupdate {
		if !filesystem.IsDirEmpty(core.GetConfig().Location.Shader) {
			return nil
		}
	}

	assetURL, assetName, tag, err := getLatestZipAsset()
	if err != nil {
		return err
	}

	versionPath := core.GetConfig().Location.Shader + "/version.txt"
	fileContent, err := filesystem.ReadFile(versionPath)
	if err == nil {
		fileContent = strings.TrimSpace(fileContent)
		if fileContent == tag {
			log.Info("Latest Anime4K version already downloaded")
			return nil
		} else if fileContent != "" {
			log.Info("Newser Version of Anime4k was found, to update delete the shaders directory")
			return nil
		}
	}

	path := fmt.Sprintf("%s/%s", core.GetConfig().Location.Temp, assetName)
	err = request.DownloadFile(assetURL, path)
	if err != nil {
		return err
	}

	tempPath := core.GetConfig().Location.Temp + "/anime4k/"
	err = filesystem.ExtractZip(path, tempPath)
	if err != nil {
		return err
	}

	err = filesystem.DeleteFile(path)
	if err != nil {
		return err
	}

	err = filesystem.CopyDirectory(tempPath, core.GetConfig().Location.Shader+"/")
	if err != nil {
		return err
	}

	err = filesystem.WriteFile(versionPath, []byte(tag))
	if err != nil {
		return err
	}

	return nil
}

func getLatestZipAsset() (assetURL string, name string, tag string, err error) {
	raw, err := request.Get(latestURL)
	if err != nil {
		return "", "", "", err
	}

	var data releaseResponse
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return "", "", "", err
	}

	if len(data.Assets) == 0 {
		return "", "", "", fmt.Errorf("no assets found")
	}

	for _, a := range data.Assets {
		if strings.HasPrefix(a.Name, "Anime4K") && strings.HasSuffix(a.Name, ".zip") {
			return a.URL, a.Name, data.TagName, nil
		}
	}

	return "", "", "", fmt.Errorf("zip asset not found")
}
