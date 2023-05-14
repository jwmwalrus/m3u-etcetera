package discover

import (
	"testing"

	"github.com/jwmwalrus/bnp/urlstr"
	"github.com/stretchr/testify/assert"
)

func TestExecDiscover(t *testing.T) {
	path := "../../data/testing/audio1/track01.ogg"

	location, err := urlstr.PathToURL(path)
	assert.NoError(t, err)

	_, err = Execute(location)
	assert.NoError(t, err)
}
