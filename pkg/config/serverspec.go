package config

import "strconv"

const (
	// DefaultServerBaseURL idem
	DefaultServerBaseURL = "http://localhost"

	// DefaultServerPort idem
	DefaultServerPort = 50099

	// DefaultServerAPIVersion idem
	DefaultServerAPIVersion = "/api/v1"
)

// ServerSpec Server-related config
type ServerSpec struct {
	BaseURL    string `json:"baseUrl"`
	Port       int    `json:"port"`
	APIVersion string `json:"apiVersion"`
}

// SetDefaults provides default settings
func (s *ServerSpec) SetDefaults() {
	if s.BaseURL == "" {
		s.BaseURL = DefaultServerBaseURL
	}

	if s.Port == 0 {
		s.Port = DefaultServerPort
	}

	if s.APIVersion == "" {
		s.APIVersion = DefaultServerAPIVersion
	}
}

// GetURL Returns the playback URL
func (s ServerSpec) GetURL() (url string) {
	url = s.BaseURL + ":" + strconv.Itoa(s.Port)
	return
}
