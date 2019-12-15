package common

import (
	"github.com/kyai/gocui"
)

type Extend struct {
	g      *gocui.Gui
	before map[string]string
}

func NewExtend(g *gocui.Gui) *Extend {
	return &Extend{
		g:      g,
		before: make(map[string]string),
	}
}

func (e *Extend) SetCurrentViewOnTop(name string) (*gocui.View, error) {
	if _, err := e.g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return e.g.SetViewOnTop(name)
}

func (e *Extend) Open(name string, x0, y0, x1, y1 int) (view *gocui.View, err error) {
	e.before[name] = e.g.CurrentView().Name()
	if view, err = e.g.SetView(name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		return e.SetCurrentViewOnTop(name)
	}
	return
}

func (e *Extend) OpenOnCenter(name string, w, h int) (view *gocui.View, err error) {
	x, y := e.g.Size()
	x0, y0 := (x-w)/2, (y-h)/2
	return e.Open(name, x0, y0, x0+w-1, y0+h-1)
}

func (e *Extend) Close(name string) (err error) {
	if err = e.g.DeleteView(name); err != nil && err != gocui.ErrUnknownView {
		return
	}
	bv, ok := e.before[name]
	if !ok {
		return nil
	}
	delete(e.before, name)
	_, err = e.SetCurrentViewOnTop(bv)
	return
}

func (e *Extend) GetCurrentLine(v *gocui.View) (string, error) {
	_, y := v.Cursor()
	return v.Line(y)
}
