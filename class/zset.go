package class

import "github.com/jroimartin/gocui"

type Zset struct {
	Base
}

func (e *Zset) Column() []string {
	return []string{"Row", "Value", "Score"}
}

func (e *Zset) Widths() []int {
	c := e.Column()
	w := e.Base.Widths(c)
	// x, _ := e.g.Size()

	w[0]++
	w[1]++

	return w
}

func (e *Zset) Render(v *gocui.View) {
	render(v, e.Rows, e.Column(), e.Widths())
}
