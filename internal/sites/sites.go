package sites

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"
	"mfg-dl/internal/sites/modules/serienstream"

	"charm.land/log/v2"
)

var Sites []model.Site

func Init() {
	if core.GetConfig().Services[aniworld.Name] {
		log.Debug("enabled", "name", aniworld.Name)
		Sites = append(Sites, aniworld.New())
	}

	if core.GetConfig().Services[serienstream.Name] {
		log.Debug("enabled", "name", serienstream.Name)
		Sites = append(Sites, serienstream.New())
	}
}
