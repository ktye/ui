package main

import (
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/fonts/apl385"
	"github.com/ktye/ui/property"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var data *app
var top *base.ButtonBar

func main() {
	data = NewApp()
	base.SetFont(apl385.TTF(), data.General.FontSize)
	top = buttonBar(base.NewLabel("click on a button for the property dialog"))
	win := ui.New(dpy.New(nil))
	win.Top = base.NewScale(top)
	done := win.Run()
	<-done
}

func buttonBar(kid ui.Widget) *base.ButtonBar {
	buttons := []*base.Button{
		&base.Button{Icon: "general", Text: "settings", Execute: setProps("General", 0)},
		base.FillButton(),
		&base.Button{Icon: "clients", Text: "clients", Execute: setTable("Clients")},
		&base.Button{Icon: "orders", Text: "orders", Execute: setTable("Orders")},
	}
	bar := base.NewButtonBar(kid, buttons)
	bar.Flip = true // show on top
	return bar
}

func setProps(name string, index int) func() int {
	return func() int {
		top.Kid.Widget = property.NewList(data, name, index)
		return -1 // redraw everything
	}
}
func setTable(name string) func() int {
	return func() int {
		t := property.NewTable(data, name)
		t.DefaultButtons()
		t.Edit = func(idx []int) int {
			if len(idx) == 0 {
				return 0
			} else if len(idx) == 1 {
				f := setProps(name, idx[0])
				return f()
			} else {
				println("TODO: multi props")
			}
			return 0
		}
		top.Kid.Widget = t
		return -1 // redraw everything
	}
}

func init() {
	base.RegisterIcon("general", icons.ActionSettings)
	base.RegisterIcon("clients", icons.ActionFace)
	base.RegisterIcon("orders", icons.ActionAssignment)
}
