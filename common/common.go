package common

func FillLeft(s string, c rune, l int) string {
	for i, n := 0, l-len(s); i < n; i++ {
		s = string(c) + s
	}
	return s
}

func FillRight(s string, c rune, l int) string {
	for i, n := 0, l-len(s); i < n; i++ {
		s += string(c)
	}
	return s
}
