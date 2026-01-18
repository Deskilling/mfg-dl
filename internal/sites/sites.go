package sites

type Site struct {
	Name    string
	Baseurl string

	Search   func(term string) ([]SearchResult, error)
	Seasons  func(result SearchResult) ([]Season, error)
	Episodes func(season Season) ([]Episode, error)
	Streams  func(episode Episode) ([]Stream, error)
	Download func(stream Stream) error
}

type SearchResult struct {
	Name           string `json:"name"`
	Href           string `json:"link"`
	Description    string `json:"description"`
	Cover          string `json:"cover"`
	ProductionYear string `json:"productionYear"`
}

type Season struct {
	Name      string
	SeasonNum string
	Href      string
	Label     string
}

type Episode struct {
	Name       string
	EpisodeNum string
	SeasonNum  string
	Href       string
	Title      string
}

type Stream struct {
	Name       string
	EpisodeNum string
	SeasonNum  string
	Href       string
	Hoster     string
	Language   string
}
