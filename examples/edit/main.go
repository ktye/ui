package main

import (
	"strconv"

	"github.com/ktye/iv/cmd/lui/apl385"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/editor"
)

func main() {
	base.SetFont(apl385.TTF(), 20)
	e := editor.New(text)
	e.Menu = e.StandardMenu()
	e.Execute = func(e *editor.Edit, s string) int { println("exec:", s); return 0 } // select and middle click
	win := ui.New(dpy.New(nil))
	win.Top = base.NewScale(e)
	done := win.Run()
	<-done
}

var text = `This is the edit widget.

Select text with the mouse, press the right button for the menu.
Execute selected text with the middle button.
Execute can be programmed to be any action.
`

func init() {
	for i := 0; i < 1000; i++ {
		text = text + "\n" + strconv.Itoa(i)
	}
}
