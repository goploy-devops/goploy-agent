package utils

import (
	"strings"
)

// GetScriptExt return script extension default bash
func GetScriptExt(scriptMode string) string {
	switch scriptMode {
	case "sh", "zsh", "bash":
		return "sh"
	case "php":
		return "php"
	case "python":
		return "py"
	case "cmd":
		return "bat"
	default:
		return "sh"
	}
}

func ClearNewline(str string) string {
	return strings.TrimRight(strings.Replace(str, "\r\n", "\n", -1), "\n")
}
