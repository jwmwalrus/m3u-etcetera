package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
)

// All -.
func All() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		m20230514080315863_add_read_only_queries(),
		m20230515200631346_add_query_id_to_playlist(),
		m20230515223654066_add_lastplayedfor_to_playlist_track(),
		m20231218164345055_add_bucket_to_playlist(),
	}
}
