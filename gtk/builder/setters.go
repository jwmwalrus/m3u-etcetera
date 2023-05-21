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

	buf, err := tv.GetBuffer()
	if err != nil {
		err = fmt.Errorf("Unable to get buffer from text view: %v", err)
		return
	}

	buf.SetText(val)
	return
}
