package class

import "github.com/jroimartin/gocui"

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
	render(v, e.Rows, e.Column(), e.Widths())
}
