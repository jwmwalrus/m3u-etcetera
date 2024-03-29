package seeds

import (
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/quorum/seater"
	"gorm.io/gorm"
)

// All seeds all initial data.
func All(db *gorm.DB) {
	h := seater.New(db)

	h.Add(seater.Seed{
		Name: "perspective",
		Run:  SeedPerspective,
	})
	h.Add(seater.Seed{
		Name: "query",
		Run:  SeedQuery,
	})
	h.AddSome([]seater.Seed{
		{
			Name:     "collection",
			Run:      SeedCollection,
			Requires: []string{"perspective"},
		},
		{
			Name:     "playlist",
			Run:      SeedPlaylist,
			Requires: []string{"perspective"},
		},
		{
			Name:     "playbar",
			Run:      SeedPlaybar,
			Requires: []string{"perspective"},
		},
		{
			Name:     "queue",
			Run:      SeedQueue,
			Requires: []string{"perspective"},
		},
	})

	onerror.Fatal(h.RunAll())
}
