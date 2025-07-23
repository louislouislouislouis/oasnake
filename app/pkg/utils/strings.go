package utils

import (
	"strings"
	"unicode"
)

func CapitalizeFirstOnly(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func TrimTrailingSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}

func RemoveBackTicks(s string) string {
	return strings.ReplaceAll(s, "`", `"`)
}

func GoCodeString(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}
