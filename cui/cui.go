package cui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

const (
	ViewInfo   = "info"
	ViewKeys   = "keys"
	ViewCond   = "cond"
	ViewOption = "option"
	ViewData   = "data"
	ViewStatus = "status"
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

	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err = g.SetKeybinding(ViewData, gocui.KeyArrowLeft, gocui.ModNone, switchKeys); err != nil {
		log.Panicln(err)
	}
	if err = g.SetKeybinding(ViewKeys, gocui.KeyArrowRight, gocui.ModNone, switchData); err != nil {
		log.Panicln(err)
	}
	if err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, switchCond); err != nil {
		log.Panicln(err)
	}

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	x, y := g.Size()

	leftX := 24

	if v, err := g.SetView(ViewInfo, 0, 0, leftX, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Info"
		v.Write([]byte("123.123.123.123:6789"))
	}

	if v, err := g.SetView(ViewKeys, 0, 4, leftX, y-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"

		if _, err = setCurrentViewOnTop(g, ViewKeys); err != nil {
			return err
		}
	}

	if v, err := g.SetView(ViewOption, leftX+1, 0, x-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Option"
	}

	if v, err := g.SetView(ViewData, leftX+1, 4, x-1, y-2); err != nil {
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
		fmt.Fprintln(v, "status bar")
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
