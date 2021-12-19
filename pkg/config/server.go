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

	// DefaultQueryLimit idem
	DefaultQueryLimit = 0

	// DefaultQueryMaxLimit idem
	DefaultQueryMaxLimit = 1023
)

// Server server-related config
type Server struct {
	Scheme     string   `json:"scheme"`
	Host       string   `json:"host"`
	Port       int      `json:"port"`
	APIVersion string   `json:"apiVersion"`
	Database   Database `json:"database"`
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
