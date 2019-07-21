package class

import "github.com/jroimartin/gocui"

type Hash struct {
	Base
}

func (e *Hash) Column() []string {
	return []string{"Row", "Key", "Value"}
}

func (e *Hash) Widths() []int {
	c := e.Column()
	w := e.Base.Widths(c)
	x, _ := e.g.Size()

	w[0] += 1
	w[1] = x - 2 - w[0]

	return w
}

func (e *Hash) Render(v *gocui.View) {
	render(v, e.Rows, e.Column(), e.Widths())
}
