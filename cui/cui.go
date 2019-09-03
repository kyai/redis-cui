package cui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/kyai/redis-cui/app"
)

const (
	ViewInfo   = "info"
	ViewKeys   = "keys"
	ViewData   = "data"
	ViewCond   = "cond"
	ViewOption = "option"
	ViewStatus = "status"
	ViewMenu   = "menu"
)

var g *gocui.Gui

func New() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	g.Update(func(*gocui.Gui) error {
		if err := renderStatusBar(); err != nil {
			return err
		}
		return renderKeys()
	})

	if err = keybind(); err != nil {
		log.Panicln(err)
	}

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	x, y := g.Size()

	leftX := 32

	if v, err := g.SetView(ViewInfo, 0, 0, leftX, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		v.Write([]byte(fmt.Sprintf("%s (db%d)", app.RedisHost, app.RedisDB)))
	}

	if v, err := g.SetView(ViewKeys, 0, 3, leftX, y-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"
		v.Highlight = true

		if _, err = setCurrentViewOnTop(g, ViewKeys); err != nil {
			return err
		}
	}

	if v, err := g.SetView(ViewOption, leftX+1, 0, x-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Detail"
	}

	if v, err := g.SetView(ViewData, leftX+1, 3, x-1, y-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Data"
	}

	if v, err := g.SetView(ViewCond, 0, y-2, leftX, y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Editable = true
		fmt.Fprintln(v, "*")
	}

	if v, err := g.SetView(ViewStatus, leftX+1, y-2, x, y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorBlue
	}

	if v, err := g.SetView("version", x-10, y-2, x, y); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorCyan
		fmt.Fprintln(v, fillAtLeft(app.VERSION, 9))
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func getCurrentLine(v *gocui.View) (string, error) {
	_, y := v.Cursor()
	return v.Line(y)
}

func fillAtLeft(s string, l int) string {
	n := l - len(s)
	for i := 0; i < n; i++ {
		s = " " + s
	}
	return s
}
