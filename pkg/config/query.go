package config

const (
	// DefaultQueryLimit idem
	DefaultQueryLimit = 0

	// DefaultQueryMaxLimit idem
	DefaultQueryMaxLimit = 1023
)

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
