package model

type Site struct {
	Name    string
	Baseurl string

	Search           func(term string) ([]SearchResult, error)
	Seasons          func(result SearchResult) ([]Season, error)
	Episodes         func(season Season) ([]Episode, error)
	Streams          func(episode Episode) ([]Stream, error)
	Download         func(stream Stream) error
	DownloadMultiple func(streams []Stream) error
}

type SearchResult struct {
	Name           string
	Href           string
	Description    string
	Cover          string
	ProductionYear string
}

type Season struct {
	Name        string
	Href        string
	SeasonNum   string
	SeasonLabel string
}

type Episode struct {
	Name                    string
	Href                    string
	SeasonNum               string
	EpisodeTitle            string
	EpisodeAlternativeTitle string
	EpisodeNum              string
}

type Stream struct {
	Name                    string
	Href                    string
	SeasonNum               string
	EpisodeTitle            string
	EpisodeAlternativeTitle string
	EpisodeNum              string
	Hoster                  string
	Language                string
}
