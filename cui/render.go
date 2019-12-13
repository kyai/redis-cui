package cui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kyai/redis-cui/app"
	"github.com/kyai/redis-cui/class"
	"github.com/kyai/redis-cui/redis"
	"github.com/mattn/go-runewidth"
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

func renderInfo() (err error) {
	view, err := g.View(ViewInfo)
	if err != nil {
		return
	}
	view.Clear()
	_, err = view.Write([]byte(fmt.Sprintf("%s (db%d)", app.RedisHost, app.RedisDB)))
	return
}

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
		fmt.Fprint(view, " "+v)
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
	key = strings.Trim(key, " ")

	vOption, _ := g.View(ViewOption)
	vData, _ := g.View(ViewData)

	entry, _ := getRedisEntry(key)
	if entry == nil {
		return
	}

	textLeft := fmt.Sprintf("%s %s", color.New(color.FgMagenta).Sprint(entry.Type), key)
	textRight := ""
	if entry.Type != "STRING" {
		textRight = fmt.Sprintf("Size:%v ", entry.Size)
	}
	textRight += fmt.Sprintf("TTL:%v", entry.TTL)
	textRight = color.New(color.FgMagenta).Sprint(textRight)

	textWidth, _ := vOption.Size()
	textBlank := textWidth - StringLen(textLeft) - StringLen(textRight)
	textSpace := strings.Join(make([]string, textBlank+1), " ")

	vOption.Clear()
	fmt.Fprintf(vOption, "%s%s%s", textLeft, textSpace, textRight)

	vData.SetOrigin(0, 0)
	vData.SetCursor(0, 0)

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

func renderStatusBar() error {
	var text string
	switch g.CurrentView().Name() {
	case ViewKeys:
		text = "m: menu, r: reload"
	case ViewData:
		text = "m: menu, r: reload, f: format"
	case ViewMenu:
		text = "m: quit"
	}
	view, _ := g.View(ViewStatus)
	view.Clear()
	view.Write([]byte(text))
	return nil
}

// render db select panel
func renderSelect() error {
	view, err := g.View(ViewSelect)
	if err != nil {
		return err
	}

	conn := redis.Pool.Get()
	defer conn.Close()

	res, err := redis.Strings(conn.Do("config", "get", "databases"))
	if err != nil {
		return err
	}
	dbs := 16
	if res[0] == "databases" {
		dbs, _ = strconv.Atoi(res[1])
	}

	view.Title = fmt.Sprintf("Index (0~%d)", dbs-1)
	view.Frame = true
	view.Editable = true
	view.Write([]byte(strconv.Itoa(app.RedisDB)))
	return nil
}

// shortcuts for menu display
var shortcuts = [][]string{
	[]string{"h(←) l(→)", "switch keys/data panel"},
	[]string{"k(↑) j(↓)", "select keys/data item"},
	[]string{"enter", "query condition"},
	[]string{"m", "toggle menu panel"},
	[]string{"s", "switch database"},
	[]string{"ctrl+c", "quit"},
}

func renderMenu() error {
	var (
		view, _ = g.View(ViewMenu)
		x, _    = view.Size()
		w1, w2  int
		content string = "\n\n"
	)

	for _, ss := range shortcuts {
		if w := runewidth.StringWidth(ss[0]); w > w1 {
			w1 = w
		}
		if w := runewidth.StringWidth(ss[1]); w > w2 {
			w2 = w
		}
	}
	w := w1 + w2 + 2
	b := (x - w) / 2

	for _, ss := range shortcuts {
		s1, s2 := ss[0], ss[1]
		for i := 0; i < b; i++ {
			content += " "
		}
		content += s1
		for i := 0; i < w1-len(s1); i++ {
			content += " "
		}
		content += "  " + s2 + "\n"
	}

	view.Clear()
	view.Write([]byte(content))
	return nil
}
