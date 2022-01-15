package seeds

import (
	"github.com/jwmwalrus/onerror"
	"github.com/jwmwalrus/seater"
	"gorm.io/gorm"
)

// All seeds all initial data
func All(db *gorm.DB) {
	h := seater.SeedHandlerNew(db)
	h.Add(seater.Seed{
		Name: "perspective",
		Run:  seedPerspective,
	})
	h.AddSome([]seater.Seed{
		{
			Name:     "collection",
			Run:      seedCollection,
			Requires: []string{"perspective"},
		},
		{
			Name:     "playlist",
			Run:      seedPlaylist,
			Requires: []string{"perspective"},
		},
		{
			Name:     "playbar",
			Run:      seedPlaybar,
			Requires: []string{"perspective"},
		},
		{
			Name:     "queue",
			Run:      seedQueue,
			Requires: []string{"perspective"},
		},
	})

	onerror.Panic(h.RunAll())
}
