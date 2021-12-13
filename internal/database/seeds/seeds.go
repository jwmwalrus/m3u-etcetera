package seeds

import (
	"github.com/jwmwalrus/bnp/onerror"
	"github.com/jwmwalrus/seater"
	"gorm.io/gorm"
)

// All seeds all initial data
func All(db *gorm.DB) {
	h := seater.SeedHandlerNew(db)
	h.AddSome([]seater.Seed{
		seater.Seed{
			Name: "perspective",
			Run:  seedPerspective,
		},
		seater.Seed{
			Name:     "queue",
			Run:      seedQueue,
			Requires: []string{"perspective"},
		},
	})

	err := h.RunAll()
	onerror.Panic(err)
}
