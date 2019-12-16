package common

func FillLeft(s string, c rune, l int) string {
	for i, n := 0, l-Characters(s); i < n; i++ {
		s = string(c) + s
	}
	return s
}

func FillRight(s string, c rune, l int) string {
	for i, n := 0, l-Characters(s); i < n; i++ {
		s += string(c)
	}
	return s
}

func Characters(s string) int {
	n := 0
	for _, _ = range s {
		n++
	}
	return n
}
