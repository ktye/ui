package main

import (
	"fmt"

	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
)

var win *ui.Window

type entry struct {
	s string
	c int // color index
}

func (e entry) String() string { return e.s }
func (e entry) Color() int     { return e.c }

func main() {
	base.SetFont(font.APL385(), 20)
	l := base.List{Single: false}
	e := make([]base.Stringer, 1000)
	for i := range e {
		c := 0
		if i == 0 {
			c = 1
		} else if i < 5 {
			c = 2
		} else if i%10 == 0 {
			c = 3
		} else if (i-1)%10 == 0 {
			c = 4
		}
		e[i] = &entry{
			s: fmt.Sprintf("this is line number %d", i),
			c: c,
		}
	}
	l.Colorsets = []base.Colorset{
		base.NewColors(base.White, base.Black),     // 1 title
		base.NewColors(base.LightGrey, base.White), // 2 light
		base.NewColors(base.Green, base.White),     // 3 green
		base.NewColors(base.Red, base.White),       // 4 red
	}
	l.Menu = base.NewMenu()
	l.List = e
	l.Execute = func() int { println("exec"); return 0 }

	win = ui.New(dpy.New(nil))
	win.Top = base.Scale{&l}
	done := win.Run()
	<-done
}