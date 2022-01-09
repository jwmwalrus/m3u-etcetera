package builder

import (
	"fmt"
)

func SetTextView(id, val string) (err error) {
	tv, err := GetTextView(id)
	if err != nil {
		return
	}

	buf, err := tv.GetBuffer()
	if err != nil {
		err = fmt.Errorf("Unable to get buffer from text view: %w", err)
		return
	}

	buf.SetText(val)
	return
}
