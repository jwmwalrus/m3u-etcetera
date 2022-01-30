package impexp

import "io"

type PLS struct {
	*playlist
}

func (pi *PLS) Format() (io.Writer, error) {
	return nil, nil
}

func (pi *PLS) Parse(f io.Reader) error {
	return nil
}

func (*PLS) Type() string {
	return PLSPlaylist.String()
}
