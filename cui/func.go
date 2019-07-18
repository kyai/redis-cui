package cui

import (
	"regexp"
	"sort"
)

func MapKeys(m map[string]string) (keys []string) {
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return
}

func StringRaw(s string) string {
	return regexp.MustCompile(`\x1b\[\d+m|\x1b\[0m`).ReplaceAllString(s, "")
}

func StringLen(s string) int {
	return len(StringRaw(s))
}
