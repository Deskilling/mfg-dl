package sites

import (
	"mfg-dl/internal/core"
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"
	"mfg-dl/internal/sites/modules/serienstream"
	"slices"

	"charm.land/log/v2"
)

var Modules []model.Module
var Languages model.Languages

func Init() {
	if core.GetConfig().Services[aniworld.Name] {
		log.Debug("enabled", "name", aniworld.Name)
		Modules = append(Modules, aniworld.New())
	}

	if core.GetConfig().Services[serienstream.Name] {
		log.Debug("enabled", "name", serienstream.Name)
		Modules = append(Modules, serienstream.New())
	}

	for _, m := range Modules {
		l := m.Languages()
		Languages.Subbed = slices.Concat(Languages.Subbed, l.Subbed)
		Languages.Dubbed = slices.Concat(Languages.Dubbed, l.Dubbed)
	}

	slices.Sort(Languages.Subbed)
	Languages.Subbed = slices.Compact(Languages.Subbed)

	slices.Sort(Languages.Dubbed)
	Languages.Dubbed = slices.Compact(Languages.Dubbed)

	log.Debug(Languages)
}
