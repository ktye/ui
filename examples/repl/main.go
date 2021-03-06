package main

import (
	"github.com/ktye/ui"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/editor"
)

func main() {
	var interp interp
	repl := &editor.Repl{Reply: true}
	repl.Edit.Nowrap = true
	repl.Interp = &interp

	win := ui.New(dpy.New(nil))
	win.Top = repl
	done := win.Run()
	<-done
}

type interp struct {
	repl *editor.Repl
}

func (i *interp) Eval(s string) string {
	return "> " + s
}
