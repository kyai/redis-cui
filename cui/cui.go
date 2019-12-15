package cui

import (
	"log"

	"github.com/kyai/gocui"
	"github.com/kyai/redis-cui/common"
)

const VERSION = "v0.2.0"

const (
	ViewInfo   = "info"
	ViewKeys   = "keys"
	ViewData   = "data"
	ViewCond   = "cond"
	ViewOption = "option"
	ViewStatus = "status"
	ViewSelect = "select"
	ViewMenu   = "menu"
)

var (
	g *gocui.Gui

	// Extended gui
	ext *common.Extend

	// The default redis query condition
	Query string
)

func New() {
	var err error
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	ext = common.NewExtend(g)

	g.Highlight = true
	g.Cursor = true
	g.InputEsc = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)

	g.Update(func(*gocui.Gui) error {
		if err := renderInfo(); err != nil {
			return err
		}
		if err := renderKeys(); err != nil {
			return err
		}
		return renderStatusBar()
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
		v.Title = "Server"
	}

	if v, err := g.SetView(ViewKeys, 0, 3, leftX, y-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Keys"
		v.Highlight = true

		if _, err = ext.SetCurrentViewOnTop(ViewKeys); err != nil {
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
		v.WriteString(Query)
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
		v.WriteString(common.FillLeft(VERSION, ' ', 9))
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
