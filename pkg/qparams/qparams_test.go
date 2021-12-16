package qparams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseParams(t *testing.T) {
	table := []struct {
		name     string
		params   string
		expected []*QParam
		err      bool
	}{
		{
			"Empty",
			"",
			nil,
			true,
		},
		{
			"Key only",
			"artist",
			nil,
			true,
		},
		{
			"No key",
			"Prince and rock",
			nil,
			true,
		},
		{
			"Simple",
			"artist=Prince",
			[]*QParam{{Key: "artist", Val: "Prince"}},
			false,
		},
		{
			"Two conditions",
			"artist=Prince and genre=rock",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
			},
			false,
		},
		{
			"With or",
			"artist=Prince or genre=rock",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Or: true, Key: "genre", Val: "rock"},
			},
			false,
		},
		{
			"With spaces",
			"artist =Prince and genre= rock or genre=pop",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
				{Or: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"With not",
			"artist =Prince and genre= rock not genre=pop",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"Consecutive conds",
			"artist=Prince and not genre=pop",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"More consecutive conds",
			"artist=Prince and or not genre=pop",
			[]*QParam{
				{Key: "artist", Val: "Prince"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			qp, err := ParseParams(tc.params)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, len(tc.expected), len(qp))
			for i := range tc.expected {
				assert.Equal(t, tc.expected[i], qp[i])
			}
		})
	}
}
