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
	Score map[string]float64
	Query map[string]string
}

type SearchResult struct {
	Service        string
	Name           string
	Href           string
	Description    string
	Cover          string
	ProductionYear string
	Score          Score
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
