package utils

import (
	"strings"
)

func ClearNewline(str string) string {
	return strings.TrimRight(strings.Replace(str, "\r\n", "\n", -1), "\n")
}
