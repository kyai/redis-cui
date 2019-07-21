package cui

import (
	"github.com/jroimartin/gocui"
)

var keyboard = []struct {
	viewname string
	key      interface{}
	mod      gocui.Modifier
	handler  func(*gocui.Gui, *gocui.View) error
}{
	{"", gocui.KeyCtrlC, gocui.ModNone, quit},
	{ViewData, gocui.KeyArrowLeft, gocui.ModNone, switchKeys},
	{ViewKeys, gocui.KeyArrowRight, gocui.ModNone, switchData},
	{"", gocui.KeyEnter, gocui.ModNone, switchCond},
	{ViewKeys, gocui.KeyArrowUp, gocui.ModNone, handleKeysPrevLine},
	{ViewKeys, gocui.KeyArrowDown, gocui.ModNone, handleKeysNextLine},
	{ViewData, gocui.KeyArrowUp, gocui.ModNone, handleDataPrevLine},
	{ViewData, gocui.KeyArrowDown, gocui.ModNone, handleDataNextLine},
}

func init() {
	// Support vim key
	for i, l := 0, len(keyboard); i < l; i++ {
		v := keyboard[i]
		if key, ok := v.key.(gocui.Key); ok {
			switch key {
			case gocui.KeyArrowUp:
				v.key = 'k'
			case gocui.KeyArrowDown:
				v.key = 'j'
			case gocui.KeyArrowLeft:
				v.key = 'h'
			case gocui.KeyArrowRight:
				v.key = 'l'
			default:
				break
			}
			keyboard = append(keyboard, v)
		}
	}
}

func keybind() (err error) {
	for _, v := range keyboard {
		if err = g.SetKeybinding(v.viewname, v.key, v.mod, v.handler); err != nil {
			return
		}
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
	return renderData()
}

func handleDataNextLine(g *gocui.Gui, v *gocui.View) error {
	return handleDataSelect(g, v, false)
}

func handleDataPrevLine(g *gocui.Gui, v *gocui.View) error {
	return handleDataSelect(g, v, true)
}

func handleDataSelect(g *gocui.Gui, v *gocui.View, up bool) error {
	if up {
		v.MoveCursor(0, -1, false)
	} else {
		v.MoveCursor(0, 1, false)
	}
	return nil
}
