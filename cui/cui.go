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

func New() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
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
		fmt.Fprintln(v, "status bar")
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
