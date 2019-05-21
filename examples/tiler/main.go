package main

import (
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
)

func main() {
	win := ui.New(dpy.New(nil))
	win.Top = base.NewTiler(base.NewLabel("drag any border to split widget"), func() ui.Widget { return base.NewLabel("new child widget") })
	done := win.Run()
	<-done
}
