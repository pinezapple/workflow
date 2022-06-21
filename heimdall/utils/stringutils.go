package utils

import "strings"

func IsStringEmpty(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}
