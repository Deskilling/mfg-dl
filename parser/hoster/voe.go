package parserHoster

type VoeStream struct {
	Key                     string   `json:"key"`
	Sharing                 bool     `json:"sharing"`
	LogoEnabled             bool     `json:"logo_enabled"`
	LogoPath                string   `json:"logo_path"`
	LogoURL                 string   `json:"logo_url"`
	LogoPosition            string   `json:"logo_position"`
	Thumbnail               string   `json:"thumbnail"`
	ShowTitle               bool     `json:"show_title"`
	Airplay                 bool     `json:"airplay"`
	Check                   bool     `json:"check"`
	FileCode                string   `json:"file_code"`
	MetadataPreload         string   `json:"metadata_preload"`
	BufferLength            int      `json:"buffer_length"`
	BufferSize              int64    `json:"buffer_size"`
	DisableTimeSlider       bool     `json:"disable_timeslider"`
	Title                   string   `json:"title"`
	Source                  string   `json:"source"`
	Fallback                []string `json:"fallback"`
	Captions                []string `json:"captions"`
	DefaultCaptionsLanguage string   `json:"default_captions_language"`
	Request                 string   `json:"request"`
	DirectAccessAllowed     bool     `json:"direct_access_allowed"`
	DirectAccessURL         string   `json:"direct_access_url"`
	SDKVersion              string   `json:"sdk_version"`
	SiteName                string   `json:"site_name"`
}

// TODO Implement Voe parser
