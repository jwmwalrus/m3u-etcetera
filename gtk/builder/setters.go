package builder

import (
	"fmt"
)

// SetTextView sets the string value for the text view identified by id.
func SetTextView(id, val string) (err error) {
	tv, err := GetTextView(id)
	if err != nil {
		return
	}

	buf := tv.Buffer()
	if buf == nil {
		err = fmt.Errorf("Unable to get buffer from text view")
		return
	}

	buf.SetText(val)
	return
}
