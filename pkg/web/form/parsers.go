package form

import "strconv"

func ParseInt(v string) int {
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return r
}

func ParseBool(v string) bool {
	return v == "on"
}