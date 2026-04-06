package sites

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"

	"charm.land/log/v2"
)

var Sites []model.Site

func Init() {
	if core.GetConfig().Services[aniworld.Name] {
		log.Debug("enabled", "name", aniworld.Name)
		Sites = append(Sites, aniworld.New())
	}
}
