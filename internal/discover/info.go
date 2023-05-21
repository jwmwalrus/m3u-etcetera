package discover

// Info defines the discovered information.
type Info struct {
	Duration int64  `json:"duration"`
	Live     bool   `json:"live"`
	Seekable bool   `json:"seekable"`
	URI      string `json:"uri"`
}
