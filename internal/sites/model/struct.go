package model

type Site interface {
	Name() (service string)
	Search(term string) ([]SearchResult, error)
	Seasons(result SearchResult) ([]Season, error)
	Episodes(season Season) ([]Episode, error)
	Streams(episode Episode) ([]Stream, error)
	Download(streams Stream) error
}

type Score struct {
	Score float64 `json:"score"`
	Query string  `json:"query"`
	Href  string  `json:"href"`
}

type SearchResult struct {
	Service        string           `json:"service"`
	Name           string           `json:"name"`
	Href           string           `json:"href"`
	Description    string           `json:"description"`
	Cover          string           `json:"cover"`
	CoverPath      string           `json:"coverPath"`
	ProductionYear string           `json:"productionYear"`
	Score          map[string]Score `json:"scores"`
}

type Season struct {
	Service     string `json:"service"`
	Name        string `json:"name"`
	Href        string `json:"href"`
	SeasonNum   string `json:"seasonNum"`
	SeasonLabel string `json:"seasonLabel"`
}

type Episode struct {
	Service                 string `json:"service"`
	Name                    string `json:"name"`
	Href                    string `json:"href"`
	SeasonNum               string `json:"seasonNum"`
	EpisodeTitle            string `json:"episodeTitle"`
	EpisodeAlternativeTitle string `json:"episodeAltTitle"`
	EpisodeNum              string `json:"episodeNum"`
}

type Stream struct {
	Service                 string `json:"service"`
	Name                    string `json:"name"`
	Href                    string `json:"href"`
	SeasonNum               string `json:"seasonNum"`
	EpisodeTitle            string `json:"episodeTitle"`
	EpisodeAlternativeTitle string `json:"episodeAltTitle"`
	EpisodeNum              string `json:"episodeNum"`
	Hoster                  string `json:"hoster"`
	Language                string `json:"language"`
}
