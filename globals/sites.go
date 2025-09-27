package globals

var streamingSites = []string{
	"aniworld",
}

func Sites() []string {
	sites := make([]string, len(streamingSites))
	copy(sites, streamingSites)
	return sites
}
