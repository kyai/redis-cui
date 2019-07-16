package cui

import (
	"github.com/jroimartin/gocui"
)

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
