package cui

import (
	"fmt"
	"redis-cui/class"
	"redis-cui/redis"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type (
	Entry struct {
		Key  string
		TTL  int
		Type string
		Size int
		Data interface{}
	}

	Hash struct {
		Key   string
		Value string
	}

	Zset struct {
		Value string
		Score int
	}
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
	view.SetOrigin(0, 0)
	view.SetCursor(0, 0)

	return renderData()
}

func renderData() (err error) {
	vKeys, _ := g.View(ViewKeys)
	_, cy := vKeys.Cursor()
	key, _ := vKeys.Line(cy)

	vOption, _ := g.View(ViewOption)
	vData, _ := g.View(ViewData)

	entry, _ := getRedisEntry(key)
	if entry == nil {
		return
	}

	textLeft := fmt.Sprintf("%s %s", color.New(color.FgBlue).Sprint(entry.Type), key)
	textRight := ""
	if entry.Type != "STRING" {
		textRight = fmt.Sprintf("Size:%v ", entry.Size)
	}
	textRight += fmt.Sprintf("TTL:%v", entry.TTL)
	textRight = color.New(color.FgBlue).Sprint(textRight)

	textWidth, _ := vOption.Size()
	textBlank := textWidth - StringLen(textLeft) - StringLen(textRight)
	textSpace := strings.Join(make([]string, textBlank+1), " ")

	vOption.Clear()
	fmt.Fprintf(vOption, "%s%s%s", textLeft, textSpace, textRight)

	e := class.NewClass(entry.Type, g)

	switch entry.Type {
	case "STRING":
		e.AddRow(entry.Data.(string))
	case "LIST", "SET":
		for k, v := range entry.Data.([]string) {
			e.AddRow(strconv.Itoa(k+1), v)
		}
	case "HASH":
		for k, v := range entry.Data.([]*Hash) {
			e.AddRow(strconv.Itoa(k+1), v.Key, v.Value)
		}
	case "ZSET":
		for k, v := range entry.Data.([]*Zset) {
			e.AddRow(strconv.Itoa(k+1), v.Value, strconv.Itoa(v.Score))
		}
	}

	e.Render(vData)

	return
}

func getRedisEntry(key string) (entry *Entry, err error) {
	conn := redis.Pool.Get()
	defer conn.Close()

	exist, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil || !exist {
		return
	}

	entry = &Entry{}
	entry.Key = key
	entry.TTL, _ = redis.Int(conn.Do("TTL", key))
	entry.Type, _ = redis.String(conn.Do("TYPE", key))
	entry.Type = strings.ToUpper(entry.Type)

	switch entry.Type {
	case "STRING":
		entry.Data, _ = redis.String(conn.Do("GET", key))
	case "LIST":
		entry.Data, _ = redis.Strings(conn.Do("LRANGE", key, 0, -1))
		entry.Size = len(entry.Data.([]string))
	case "HASH":
		data, _ := redis.StringMap(conn.Do("HGETALL", key))
		keys := MapKeys(data)
		hashs := make([]*Hash, len(keys))
		for k, v := range keys {
			hashs[k] = &Hash{v, data[v]}
		}
		entry.Data = hashs
		entry.Size = len(hashs)
	case "SET":
		entry.Data, _ = redis.Strings(conn.Do("SMEMBERS", key))
		entry.Size = len(entry.Data.([]string))
	case "ZSET":
		data, _ := redis.Strings(conn.Do("ZRANGE", key, 0, -1, "WITHSCORES"))
		zsets := make([]*Zset, len(data)/2)
		for i := 0; i < len(data); i += 2 {
			score, _ := strconv.Atoi(data[i+1])
			zsets[i/2] = &Zset{data[i], score}
		}
		entry.Data = zsets
		entry.Size = len(zsets)
	}

	return
}

func renderTest(s string) {
	v, _ := g.View(ViewData)
	v.Clear()
	v.Write([]byte(s))
}
