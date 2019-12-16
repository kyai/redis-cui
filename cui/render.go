package cui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/go-redis/redis/v7"
	"github.com/kyai/redis-cui/class"
	"github.com/kyai/redis-cui/common"
)

type Entry struct {
	Key  string
	TTL  float64
	Type string
	Size int
	Data interface{}
}

func renderInfo() (err error) {
	view, err := g.View(ViewInfo)
	if err != nil {
		return
	}
	view.Clear()
	opt := client.Options()
	_, err = view.Write([]byte(fmt.Sprintf("%s[%d]", opt.Addr, opt.DB)))
	return
}

func renderKeys() (err error) {
	view, err := g.View(ViewCond)
	if err != nil {
		return
	}
	cond := view.ViewBufferLines()[0]

	keys, err := client.Keys(cond).Result()
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

	if len(key) == 0 {
		vOption.Clear()
		vData.Clear()
		return nil
	}

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
		for _, v := range entry.Data.([]string) {
			e.AddRow(v)
		}
	case "HASH":
		for k, v := range entry.Data.(map[string]string) {
			e.AddRow(k, v)
		}
	case "ZSET":
		for _, v := range entry.Data.([]redis.Z) {
			e.AddRow(v.Member.(string), fmt.Sprint(v.Score))
		}
	}

	e.Render(vData)

	return
}

func getRedisEntry(key string) (entry *Entry, err error) {
	if client.Exists(key).Val() == 0 {
		return
	}

	entry = &Entry{}
	entry.Key = key
	entry.TTL = client.TTL(key).Val().Seconds()
	if entry.TTL < 0 {
		entry.TTL = -1
	}
	entry.Type = client.Type(key).Val()
	entry.Type = strings.ToUpper(entry.Type)

	switch entry.Type {
	case "STRING":
		entry.Data = client.Get(key).Val()
	case "LIST":
		entry.Data = client.LRange(key, 0, -1).Val()
		entry.Size = len(entry.Data.([]string))
	case "HASH":
		entry.Data = client.HGetAll(key).Val()
		entry.Size = len(entry.Data.(map[string]string))
	case "SET":
		entry.Data = client.SMembers(key).Val()
		entry.Size = len(entry.Data.([]string))
	case "ZSET":
		entry.Data = client.ZRangeWithScores(key, 0, -1).Val()
		entry.Size = len(entry.Data.([]redis.Z))
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
	case ViewMenu, ViewSelect:
		text = "esc: close"
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

	dbs, err := getDatabases()
	if err != nil {
		return err
	}

	view.Title = fmt.Sprintf("Index (0~%d)", dbs-1)
	view.Frame = true
	view.Editable = true
	view.Write([]byte(strconv.Itoa(client.Options().DB)))
	return nil
}

// shortcuts for menu display
var shortcuts = []struct {
	Key string
	Des string
}{
	{"h(←) l(→)", "Switch keys/data panel"},
	{"k(↑) j(↓)", "Select keys/data item"},
	{"enter", "Query condition"},
	{"s", "Select database"},
	{"m", "Display menu"},
	{"ctrl+c", "Quit"},
}

func renderMenu() error {
	var (
		view, _ = g.View(ViewMenu)
		x, _    = view.Size()
		w1, w2  int
		content string = "\n\n"
		cnum           = common.Characters
	)

	for _, s := range shortcuts {
		if w := cnum(s.Key); w > w1 {
			w1 = w
		}
		if w := cnum(s.Des); w > w2 {
			w2 = w
		}
	}
	w := w1 + w2 + 2
	b := (x - w) / 2

	for _, s := range shortcuts {
		content += common.FillLeft(s.Key, ' ', b+cnum(s.Key))
		content += common.FillLeft(s.Des, ' ', w1-cnum(s.Key)+cnum(s.Des)+2)
		content += "\n"
	}

	view.Clear()
	view.Write([]byte(content))
	return nil
}
