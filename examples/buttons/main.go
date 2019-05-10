package main

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/fonts/apl385"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var win *ui.Window

func init() {
	base.RegisterIcon("ActionDone", icons.ActionDone)
	base.RegisterIcon("ActionSettingsInputComponent", icons.ActionSettingsInputComponent)
	base.RegisterIcon("NavigationCancel", icons.NavigationCancel)
	base.RegisterIcon("PlacesSmokeFree", icons.PlacesSmokeFree)
	base.RegisterIcon("SocialPerson", icons.SocialPerson)
	base.RegisterIcon("ToggleIndeterminateCheckBox", icons.ToggleIndeterminateCheckBox)
}

func main() {
	var bp *base.ButtonBar
	var textButtons, iconButtons, bothButtons []*base.Button

	textButtons = []*base.Button{
		base.FillButton(),
		base.NewButton("Flip", "", func() int { bp.Flip = !bp.Flip; return -1 }),
		base.FillButton(),
		base.NewButton("Icons", "", func() int { bp.SetButtons(iconButtons); return -1 }),
		base.FillButton(),
		base.NewButton("Vertical", "", func() int { bp.Vertical = !bp.Vertical; return -1 }),
		base.FillButton(),
	}

	iconButtons = []*base.Button{
		base.NewButton("", "ActionDone", func() int { bp.Flip = !bp.Flip; return -1 }),
		base.NewButton("", "ActionSettingsInputComponent", func() int { bp.SetButtons(bothButtons); return -1 }),
		base.FillButton(),
		base.NewButton("", "NavigationCancel", func() int { bp.Vertical = !bp.Vertical; return -1 }),
	}

	bothButtons = []*base.Button{
		base.NewButton("Flip", "PlacesSmokeFree", func() int { bp.Flip = !bp.Flip; return -1 }),
		base.NewButton("Icons", "SocialPerson", func() int { bp.SetButtons(textButtons); return -1 }),
		base.NewButton("Vertical", "ToggleIndeterminateCheckBox", func() int { bp.Vertical = !bp.Vertical; return -1 }),
	}

	bar := base.NewButtonBar(&panel{}, textButtons)
	bar.Target = &bp
	win = ui.New(dpy.New(nil))
	win.Top = base.NewScale(bar)
	base.SetFont(apl385.TTF(), 20)
	done := win.Run()
	<-done
}

// Panel is a dummy widget that will be the main widget (Kid) of the ButtonBar.
type panel struct{}

func (p *panel) Draw(dst *image.RGBA, force bool) {
	if force {
		draw.Draw(dst, dst.Bounds(), base.DeepOrange.Uniform(), image.ZP, draw.Src)
	}
}
func (p *panel) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if but == 1 && dir > 0 {
		return 1
	}
	return 0
}
func (p *panel) Key(r rune, code uint32, dir int, mod uint32) int {
	return 0
}
