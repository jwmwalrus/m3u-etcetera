package config

import "strconv"

const (
	// DefaultServerScheme -
	DefaultServerScheme = "http"

	// DefaultServerHost -
	DefaultServerHost = "localhost"

	// DefaultServerPort -
	DefaultServerPort = 50099

	// DefaultServerAPIVersion -
	DefaultServerAPIVersion = "/api/v1"

	// DefaultQueryLimit -
	DefaultQueryLimit = 0

	// DefaultQueryMaxLimit -
	DefaultQueryMaxLimit = 1023
)

// Server server-related config
type Server struct {
	Scheme     string `json:"scheme"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	APIVersion string `json:"apiVersion"`

	Database struct {
		Backup bool `json:"backup"`
	} `json:"database"`

	Playback struct {
	} `json:"playback"`

	Query struct {
		Limit int `json:"limit"`
	} `json:"query"`

	Collection struct {
		Scanning struct {
			SkipCover bool `json:"skipCover"`
		} `json:"scanning"`
	} `json:"collection"`
}

// SetDefaults provides default settings
func (s *Server) SetDefaults() {
	if s.Scheme == "" {
		s.Scheme = DefaultServerScheme
	}

	if s.Host == "" {
		s.Host = DefaultServerHost
	}

	if s.Port == 0 {
		s.Port = DefaultServerPort
	}

	if s.APIVersion == "" {
		s.APIVersion = DefaultServerAPIVersion
	}

	s.Database.Backup = true

	if s.Query.Limit == 0 {
		s.Query.Limit = DefaultQueryLimit
	}
}

// GetAuthority returns the authority portion of the playback URI
func (s *Server) GetAuthority() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// GetURI returns the playback URI
func (s *Server) GetURI() string {
	return s.Scheme + "://" + s.GetAuthority()
}
