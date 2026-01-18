package util

import "flag"

type arguments struct {
	List bool

	Site   string
	Mode   string
	Search string

	DownloadSeason bool

	Episode string
	Season  string
}

func Flags() arguments {
	var flags arguments

	flag.BoolVar(&flags.List, "list", false, "allows you to list all available sites")
	flag.StringVar(&flags.Site, "site", "aniworld", "allows you to set the desired website")
	flag.StringVar(&flags.Mode, "mode", "search", "select mode valid are search and download")
	flag.StringVar(&flags.Search, "search", "", "allow you to search for the show on the given host (only available when mode is search)")

	flag.BoolVar(&flags.DownloadSeason, "downloadSeason", true, "select if the whole season should be downloaded")

	flag.StringVar(&flags.Season, "season", "1", "select the season you want to download")
	flag.StringVar(&flags.Episode, "episode", "1", "select the episode you want to download")

	flag.Parse()
	return flags
}
