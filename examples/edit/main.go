package main

import (
	"github.com/ktye/iv/cmd/lui/apl385"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/editor"
)

func main() {
	base.SetFont(apl385.TTF(), 20)
	e := editor.New(text)
	e.Execute = func(e *editor.Edit, s string) int { println("exec:", s); return 0 } // select and middle click
	win := ui.New(dpy.New(nil))
	win.Top = base.Scale{e}
	done := win.Run()
	<-done
}

const text = `This is the edit widget.

Search text with the right mouse button, execute with the middle.
`
