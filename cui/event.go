package cui

import (
	"github.com/jroimartin/gocui"
)

func keybind() (err error) {
	if err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return
	}
	if err = g.SetKeybinding(ViewData, gocui.KeyArrowLeft, gocui.ModNone, switchKeys); err != nil {
		return
	}
	if err = g.SetKeybinding(ViewKeys, gocui.KeyArrowRight, gocui.ModNone, switchData); err != nil {
		return
	}
	if err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, switchCond); err != nil {
		return
	}
	if err = g.SetKeybinding(ViewKeys, gocui.KeyArrowUp, gocui.ModNone, handleKeysPrevLine); err != nil {
		return
	}
	if err = g.SetKeybinding(ViewKeys, gocui.KeyArrowDown, gocui.ModNone, handleKeysNextLine); err != nil {
		return
	}
	return
}

func switchKeys(g *gocui.Gui, v *gocui.View) error {
	_, err := setCurrentViewOnTop(g, ViewKeys)
	return err
}

func switchData(g *gocui.Gui, v *gocui.View) error {
	_, err := setCurrentViewOnTop(g, ViewData)
	return err
}

func switchCond(g *gocui.Gui, v *gocui.View) error {
	name := ViewCond
	if v.Name() == ViewCond {
		name = ViewKeys
		renderKeys()
	}
	_, err := setCurrentViewOnTop(g, name)
	return err
}

func handleKeysNextLine(g *gocui.Gui, v *gocui.View) error {
	return handleKeysSelect(g, v, false)
}

func handleKeysPrevLine(g *gocui.Gui, v *gocui.View) error {
	return handleKeysSelect(g, v, true)
}

func handleKeysSelect(g *gocui.Gui, v *gocui.View, up bool) error {
	if up {
		v.MoveCursor(0, -1, false)
	} else {
		v.MoveCursor(0, 1, false)
	}
	return nil
}
