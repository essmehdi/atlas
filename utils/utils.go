package utils

import "strings"

func ArrayContains[N comparable](needle N, haystack []N) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}
	return false
}

func IndentStringByLevel(level int, str string) string {
	prefix := strings.Repeat("  ", level) + "|"
	str = strings.ReplaceAll(str, "\n", "\n" + prefix)
	return prefix + str
}