package utils

import "strings"

func UCFirst(str string) string {
	firstChar := strings.ToUpper(str[:1])
	return firstChar + str[1:]
}