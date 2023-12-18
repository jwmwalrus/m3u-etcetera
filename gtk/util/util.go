package util

import (
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/glib/v2"
)

// IDListToGValue converts an ID list to a glib.Value.
func IDListToGValue(ids []int64) *glib.Value {
	return glib.NewValue(IDListToString(ids))
}

// IDListToString converts an ID list to a comma-separaterd string.
func IDListToString(ids []int64) (s string) {
	if len(ids) < 1 {
		return
	}
	s = strconv.FormatInt(ids[0], 10)
	for i := 1; i < len(ids); i++ {
		s += "," + strconv.FormatInt(ids[i], 10)
	}
	return
}

// StringToIDList parses the IDList column.
func StringToIDList(s string) (ids []int64, err error) {
	if len(s) == 0 {
		return
	}
	list := strings.Split(s, ",")
	for _, l := range list {
		var id int64
		id, err = strconv.ParseInt(l, 10, 64)
		if err != nil {
			return
		}
		ids = append(ids, id)
	}
	return
}
