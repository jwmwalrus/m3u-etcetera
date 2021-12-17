package config

import "strconv"

const (
	// DefaultServerScheme idem
	DefaultServerScheme = "http"

	// DefaultServerHost idem
	DefaultServerHost = "localhost"

	// DefaultServerPort idem
	DefaultServerPort = 50099

	// DefaultServerAPIVersion idem
	DefaultServerAPIVersion = "/api/v1"
)

// Server server-related config
type Server struct {
	Scheme     string `json:"scheme"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	APIVersion string `json:"apiVersion"`
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
}

// GetAuthority returns the authority portion of the playback URI
func (s *Server) GetAuthority() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// GetURI returns the playback URI
func (s *Server) GetURI() string {
	return s.Scheme + "://" + s.GetAuthority()
}
