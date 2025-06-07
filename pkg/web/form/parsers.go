package form

import (
	"strconv"
	"strings"
)


func ParseInt(v string) int {
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return r
}

func ParseBool(v string) bool {
	return strings.ToLower(v) == "on"
}
