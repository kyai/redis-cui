package cui

import (
	"redis-cui/redis"
	"strings"
)

func renderKeys() (err error) {
	v, err := g.View(ViewCond)
	if err != nil {
		return
	}
	cond := v.ViewBuffer()

	conn := redis.Pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", cond))
	if err != nil {
		return
	}
	renderTest(strings.Join(keys, "|"))
	return
}

func renderTest(s string) {
	v, _ := g.View(ViewData)
	v.Write([]byte(s))
}
