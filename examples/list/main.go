package main

import (
	"fmt"

	"github.com/ktye/iv/cmd/lui/font"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
)

var win *ui.Window

func main() {
	base.SetFont(font.APL385(), 20)
	l := base.List{Single: false}
	s := make([]string, 1000)
	for i := range s {
		s[i] = fmt.Sprintf("this is line number %d", i)
	}
	l.Colorsets = []base.Colorset{
		base.NewColors(base.White, base.Black),     // 1 title
		base.NewColors(base.LightGrey, base.White), // 2 light
		base.NewColors(base.Green, base.White),     // 3 green
		base.NewColors(base.Red, base.White),       // 4 red
	}
	l.Colors = make([]int, len(s))
	l.Colors[0] = 1
	l.Menu = base.NewMenu()
	for i := range l.Colors {
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
		l.Colors[i] = c
	}
	l.SetList(s)

	/*
		l.Colors = []base.Colorset{
			base.NewColors(0xFF0000, 0x00FF00),
			base.NewColors(0xFFFFFF, 0x000000),
		}
	*/

	win = ui.New(dpy.New(nil))
	win.Top = base.Scale{&l}
	done := win.Run()
	<-done
}
