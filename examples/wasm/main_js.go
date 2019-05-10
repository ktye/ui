// +build js

package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpw"
	"github.com/ktye/ui/fonts/apl385"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

// Build main.wasm (from this file) with:
//	GOOS=js GOARCH=wasm go build -o main.wasm
// Then serve to localhost (unset GOOS and GOARCH):
//	go run serve.go
func main() {
	println("this is wasm")

	var bp *base.ButtonBar
	var textButtons, iconButtons, bothButtons []*base.Button

	textButtons = []*base.Button{
		base.NewButton("Ok", "", nil),
		base.NewButton("Icons", "", func() int { bp.SetButtons(iconButtons); return -1 }),
		base.NewButton("Cancel", "", nil),
	}

	iconButtons = []*base.Button{
		base.NewButton("", "ActionDone", nil),
		base.NewButton("", "ActionSettingsInputComponent", func() int { bp.SetButtons(bothButtons); return -1 }),
		base.NewButton("", "NavigationCancel", nil),
	}

	bothButtons = []*base.Button{
		base.NewButton("Ok", "PlacesSmokeFree", nil),
		base.NewButton("Icons", "SocialPerson", func() int { bp.SetButtons(textButtons); return -1 }),
		base.NewButton("Ignore", "ToggleIndeterminateCheckBox", nil),
	}

	bar := base.NewButtonBar(&basic{}, textButtons)
	bar.Target = &bp
	win := ui.New(dpw.New("dpy"))
	win.Top = base.NewScale(bar)
	base.SetFont(apl385.TTF(), 20)
	done := win.Run()
	<-done

}

func init() {
	base.RegisterIcon("ActionDone", icons.ActionDone)
	base.RegisterIcon("ActionSettingsInputComponent", icons.ActionSettingsInputComponent)
	base.RegisterIcon("NavigationCancel", icons.NavigationCancel)
	base.RegisterIcon("PlacesSmokeFree", icons.PlacesSmokeFree)
	base.RegisterIcon("SocialPerson", icons.SocialPerson)
	base.RegisterIcon("ToggleIndeterminateCheckBox", icons.ToggleIndeterminateCheckBox)
}

type basic struct{}

func (b *basic) Draw(dst *image.RGBA, force bool) {
	println("draw", force)
	if force {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{0xFF, 0x00, 0xFF, 0xFF}}, image.ZP, draw.Src)
	}
}
func (b *basic) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	println("mouse", pos.String(), but, dir, mod)
	if but == 1 && dir > 0 {
		return 1
	}
	return 0
}
func (b *basic) Key(r rune, code uint32, dir int, mod uint32) int {
	println("key", code, dir, mod)
	return 0
}
