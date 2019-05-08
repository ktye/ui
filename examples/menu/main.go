package main

import (
	"fmt"

	"github.com/ktye/iv/cmd/lui/apl385"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
)

func main() {
	base.SetFont(apl385.TTF(), 20)
	s := make([]string, 1000)
	for i := range s {
		s[i] = fmt.Sprintf("this is line number %d", i)
	}
	l := base.NewList(s)

	buttons := []string{"File", "New", "Delete", "Accept"}
	m := make([]*base.Button, len(buttons))
	for i := range m {
		m[i] = base.NewButton(buttons[i], "", nil)
	}
	m[2].Execute = func() int {
		var s []base.Stringer
		for i := range l.Sel {
			if l.Sel[i] == false {
				s = append(s, l.List[i])
			}
		}
		l.List = s
		return -1
	}
	l.Menu = base.NewMenu(m)

	win := ui.New(dpy.New(nil))
	win.Top = base.NewScale(l)
	done := win.Run()
	<-done
}
