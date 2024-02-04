package utils

import "strings"

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}
