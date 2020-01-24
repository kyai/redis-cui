package class

import "github.com/kyai/gocui"

type Stream struct {
	Base
}

func (e *Stream) Column() []string {
	return []string{"ID", "Values"}
}

func (e *Stream) Widths() []int {
	c := e.Column()
	w := e.Base.Widths(c)
	// x, _ := e.g.Size()

	w[0]++
	w[1]++

	return w
}

func (e *Stream) Render(v *gocui.View) {
	e.render(v, e.Rows, e.Column(), e.Widths())
}
