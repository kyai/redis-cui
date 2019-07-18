package cui

import (
	"fmt"
	"redis-cui/redis"
	"sort"
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
	sort.Strings(keys)
	for k, v := range keys {
		fmt.Fprint(view, v)
		if k < len(keys)-1 {
			fmt.Fprintln(view)
		}
	}

	return
}

func renderTest(s string) {
	v, _ := g.View(ViewData)
	v.Write([]byte(s))
}
