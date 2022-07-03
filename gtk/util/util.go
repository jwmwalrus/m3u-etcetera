package util

import (
	"strconv"
	"strings"
)

// IDListToString converts an ID list to a comma-separaterd string
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

// StringToIDList parses the IDList column
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
