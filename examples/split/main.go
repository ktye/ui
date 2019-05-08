package main

import (
	"fmt"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
)

var win *ui.Window

func main() {
	top := base.NewSplit(base.NewList(strings(12)), base.NewList(strings(100)))

	win = ui.New(dpy.New(nil))
	win.Top = base.NewScale(top)
	done := win.Run()
	<-done
}

func strings(n int) []string {
	s := make([]string, n)
	for i := range s {
		s[i] = fmt.Sprintf("this is line %d", i)
	}
	return s
}
