package ext

import (
	"github.com/kyai/gocui"
)

var (
	g *gocui.Gui

	before map[string]string
)

func Init(gui *gocui.Gui) {
	if g == nil {
		g = gui
	}
	before = make(map[string]string)
}

func SetCurrentViewOnTop(name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func Open(name string, x0, y0, x1, y1 int) (view *gocui.View, err error) {
	before[name] = g.CurrentView().Name()
	if view, err = g.SetView(name, x0, y0, x1, y1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		return SetCurrentViewOnTop(name)
	}
	return
}

func OpenOnCenter(name string, w, h int) (view *gocui.View, err error) {
	x, y := g.Size()
	x0, y0 := (x-w)/2, (y-h)/2
	return Open(name, x0, y0, x0+w-1, y0+h-1)
}

func Close(name string) (err error) {
	if err = g.DeleteView(name); err != nil && err != gocui.ErrUnknownView {
		return
	}
	bv, ok := before[name]
	if !ok {
		return nil
	}
	delete(before, name)
	_, err = SetCurrentViewOnTop(bv)
	return
}
