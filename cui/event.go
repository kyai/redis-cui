package cui

import (
	"strconv"

	"github.com/kyai/gocui"
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
	{ViewKeys, gocui.KeyEnter, gocui.ModNone, switchCond},
	{ViewData, gocui.KeyEnter, gocui.ModNone, switchCond},
	{ViewCond, gocui.KeyEnter, gocui.ModNone, switchCond},
	{ViewKeys, gocui.KeyArrowUp, gocui.ModNone, handleKeysPrevLine},
	{ViewKeys, gocui.KeyArrowDown, gocui.ModNone, handleKeysNextLine},
	{ViewData, gocui.KeyArrowUp, gocui.ModNone, handleDataPrevLine},
	{ViewData, gocui.KeyArrowDown, gocui.ModNone, handleDataNextLine},
	{ViewKeys, 's', gocui.ModNone, handleDbSelect},
	{ViewData, 's', gocui.ModNone, handleDbSelect},
	{ViewKeys, 'r', gocui.ModNone, handleDataReload},
	{ViewData, 'r', gocui.ModNone, handleDataReload},
	{ViewKeys, 'm', gocui.ModNone, handleMenuToggle},
	{ViewData, 'm', gocui.ModNone, handleMenuToggle},
	{ViewMenu, gocui.KeyEsc, gocui.ModNone, handleClose},
	{ViewSelect, gocui.KeyEnter, gocui.ModNone, handleDbSelectDo},
	{ViewSelect, gocui.KeyEsc, gocui.ModNone, handleClose},
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
		g.SetKeybinding(v.viewname, v.key, v.mod, handleStatusBar)
	}
	return
}

func switchKeys(g *gocui.Gui, v *gocui.View) error {
	_, err := ext.SetCurrentViewOnTop(ViewKeys)
	return err
}

func switchData(g *gocui.Gui, v *gocui.View) error {
	_, err := ext.SetCurrentViewOnTop(ViewData)
	return err
}

func switchCond(g *gocui.Gui, v *gocui.View) error {
	name := ViewCond
	if v.Name() == ViewCond {
		name = ViewKeys
		renderKeys()
	}
	_, err := ext.SetCurrentViewOnTop(name)
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

// reload data for selected key
func handleDataReload(g *gocui.Gui, v *gocui.View) error {
	return renderData()
}

func handleStatusBar(g *gocui.Gui, v *gocui.View) error {
	return renderStatusBar()
}

func handleDbSelect(g *gocui.Gui, v *gocui.View) error {
	_, err := ext.OpenOnCenter(ViewSelect, 20, 3)
	if err != nil {
		return err
	}

	return renderSelect()
}

func handleDbSelectDo(g *gocui.Gui, v *gocui.View) error {
	val := v.ViewBufferLines()[0]
	db, err := strconv.Atoi(val)
	if err != nil {
		return nil // do not switch
	}

	err = ext.Close(v.Name())
	if err != nil {
		return err
	}

	selectDatabase(db)

	if err = renderInfo(); err != nil {
		return err
	}

	return renderKeys()
}

func handleMenuToggle(g *gocui.Gui, v *gocui.View) error {
	view, err := ext.OpenOnCenter(ViewMenu, 60, len(shortcuts)+4)
	if err != nil {
		return err
	}
	view.Title = "Menu"

	return renderMenu()
}

func handleClose(g *gocui.Gui, v *gocui.View) error {
	return ext.Close(v.Name())
}
