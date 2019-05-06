package main

import (
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/editor"
)

func main() {
	var interp interp
	repl := &editor.Repl{Reply: true}
	repl.Edit.Nowrap = true

	interp.repl = repl
	repl.Interp = &interp

	win := ui.New(dpy.New(nil))
	win.Top = &base.Scale{repl}
	done := win.Run()
	<-done
}

type interp struct {
	repl *editor.Repl
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	i.repl.Write([]byte("> " + s + "\n"))
	i.repl.Write([]byte("\t"))
	i.repl.Edit.MarkAddr("$")
}

func (i *interp) Cancel() {}
