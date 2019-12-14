package class

import "github.com/kyai/gocui"

type List struct {
	Base
}

func (e *List) Column() []string {
	return []string{"Row", "Value"}
}

func (e *List) Widths() []int {
	c := e.Column()
	w := e.Base.Widths(c)
	x, _ := e.g.Size()

	w[0] += 1
	w[1] = x - 2 - w[0]

	return w
}

func (e *List) Render(v *gocui.View) {
	e.rownum()
	e.render(v, e.Rows, e.Column(), e.Widths())
}
