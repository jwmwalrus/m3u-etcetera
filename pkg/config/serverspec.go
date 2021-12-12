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

// ServerSpec Server-related config
type ServerSpec struct {
	Scheme     string `json:"scheme"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	APIVersion string `json:"apiVersion"`
}

// SetDefaults provides default settings
func (s *ServerSpec) SetDefaults() {
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
func (s *ServerSpec) GetAuthority() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// GetURI returns the playback URI
func (s *ServerSpec) GetURI() string {
	return s.Scheme + "://" + s.GetAuthority()
}
