package class

import (
	"encoding/json"
	"strings"

	"github.com/kyai/gocui"
)

type String struct {
	Base
}

func (e *String) Column() []string {
	return []string{"Value"}
}

func (e *String) Widths() []int {
	x, _ := e.g.Size()
	return []int{x - 2}
}

func (e *String) Render(v *gocui.View) {
	// format json
	if s := e.Rows[0][0]; len(s) > 0 {
		var j interface{}
		if err := json.Unmarshal([]byte(s), &j); err == nil {
			if b, err := json.MarshalIndent(j, "", "  "); err == nil {
				e.Rows[0][0] = strings.ReplaceAll(string(b), "\n", "\n ")
			}
		}
	}
	e.render(v, e.Rows, e.Column(), e.Widths())
}
