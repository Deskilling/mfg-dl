package sites

import (
	"mfg-dl/internal/sites/model"
	"mfg-dl/internal/sites/modules/aniworld"
)

var Sites = []model.Site{
	aniworld.New(),
}
