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

	// DefaultPlayedThreshold -
	DefaultPlayedThreshold = 30

	// DefaultQueryLimit -
	DefaultQueryLimit = 0

	// DefaultQueryMaxLimit -
	DefaultQueryMaxLimit = 1023
)

// Server server-related config
type Server struct {
	Scheme     string   `json:"scheme"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	APIVersion string   `json:"apiVersion"`
	Database   Database `json:"database"`
	Playback   Playback `json:"playback"`
	Query      Query    `json:"query"`
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

	s.Database.SetDefaults()
	s.Playback.SetDefaults()
	s.Query.SetDefaults()
}

// GetAuthority returns the authority portion of the playback URI
func (s *Server) GetAuthority() string {
	return s.Host + ":" + strconv.Itoa(s.Port)
}

// GetURI returns the playback URI
func (s *Server) GetURI() string {
	return s.Scheme + "://" + s.GetAuthority()
}

// Database database-related config
type Database struct {
	Backup bool `json:"backup"`
}

// SetDefaults provides default settings
func (db *Database) SetDefaults() {
	db.Backup = true
}

// Playback playback-related config
type Playback struct {
	PlayedThreshold int `json:"playedThreshold"`
}

// SetDefaults provides default settings
func (pb *Playback) SetDefaults() {
	if pb.PlayedThreshold == 0 {
		pb.PlayedThreshold = DefaultPlayedThreshold
	}
}

// Query query-related config
type Query struct {
	Limit int `json:"limit"`
}

// SetDefaults provides default settings
func (q *Query) SetDefaults() {
	if q.Limit == 0 {
		q.Limit = DefaultQueryLimit
	}
}
