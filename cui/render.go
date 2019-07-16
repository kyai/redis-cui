package cui

import (
	"fmt"
	"redis-cui/redis"
)

func renderKeys() (err error) {
	view, err := g.View(ViewCond)
	if err != nil {
		return
	}
	cond := view.ViewBufferLines()[0]

	conn := redis.Pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", cond))
	if err != nil {
		return
	}

	view, err = g.View(ViewKeys)
	if err != nil {
		return
	}

	view.Clear()
	for _, v := range keys {
		fmt.Fprintln(view, v)
	}

	return
}

func renderTest(s string) {
	v, _ := g.View(ViewData)
	v.Write([]byte(s))
}
