package cui

import (
	"fmt"
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

func StreamFmt(values map[string]interface{}) string {
	var keys []string
	for k, _ := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := ""
	for _, key := range keys {
		s += fmt.Sprintf("%s:%s ", key, values[key])
	}
	return s
}
