package qparams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: test with leading condition
func TestParseParams(t *testing.T) {
	table := []struct {
		name     string
		params   string
		expected []QParam
		err      bool
	}{
		{
			"Empty",
			"",
			[]QParam{},
			true,
		},
		{
			"Key only=",
			"artist",
			[]QParam{},
			true,
		},
		{
			"Value only",
			"Prince and rock",
			[]QParam{},
			true,
		},
		{
			"Simple",
			"artist=Prince",
			[]QParam{{Key: "artist", Val: "Prince"}},
			false,
		},
		{
			"Leading cond",
			"and artist=Prince",
			[]QParam{{Key: "artist", Val: "Prince"}},
			false,
		},
		{
			"Two conditions",
			"artist=Prince and genre=rock",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
			},
			false,
		},
		{
			"With or",
			"artist=Prince or genre=rock",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Or: true, Key: "genre", Val: "rock"},
			},
			false,
		},
		{
			"With spaces",
			"artist =Prince and genre= rock or genre=pop",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
				{Or: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"With not",
			"artist =Prince and genre= rock not genre=pop",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Key: "genre", Val: "rock"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"Consecutive conds",
			"artist=Prince and not genre=pop",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"More consecutive conds",
			"artist=Prince and or not genre=pop",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Not: true, Key: "genre", Val: "pop"},
			},
			false,
		},
		{
			"CSV simple",
			"id=5050,23748,23761",
			[]QParam{
				{Key: "id", Val: "5050"},
				{Or: true, Key: "id", Val: "23748"},
				{Or: true, Key: "id", Val: "23761"},
			},
			false,
		},
		{
			"CSV with or",
			"or id=5050,23748,23761",
			[]QParam{
				{Or: true, Key: "id", Val: "5050"},
				{Or: true, Key: "id", Val: "23748"},
				{Or: true, Key: "id", Val: "23761"},
			},
			false,
		},
		{
			"CSV with not",
			"not id=5050,23748,23761",
			[]QParam{
				{Not: true, Key: "id", Val: "5050"},
				{Not: true, Key: "id", Val: "23748"},
				{Not: true, Key: "id", Val: "23761"},
			},
			false,
		},
		{
			"CSV with spaces",
			"artist=Prince, Cher, George Michael",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Or: true, Key: "artist", Val: "Cher"},
				{Or: true, Key: "artist", Val: "George Michael"},
			},
			false,
		},
		{
			"CSV with empty",
			"artist=Prince, Cher, George Michael, and genre=pop",
			[]QParam{
				{Key: "artist", Val: "Prince"},
				{Or: true, Key: "artist", Val: "Cher"},
				{Or: true, Key: "artist", Val: "George Michael"},
				{Key: "genre", Val: "pop"},
			},
			false,
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			qp, err := ParseParams(tc.params)
			assert.Equal(t, tc.err, err != nil)
			assert.Equal(t, len(tc.expected), len(qp))
			for i := range qp {
				assert.Contains(t, tc.expected, *qp[i])
			}
		})
	}
}
