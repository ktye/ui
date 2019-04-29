package base

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
)

type Base struct {
	Menu *Menu
	Rect image.Rectangle
	draw bool
}

func (b *Base) DrawSelf() int {
	b.draw = true
	return 1
}

type Kid struct {
	ui.Widget
	image.Rectangle
}

func (k *Kid) Draw(dst *image.RGBA, force bool) {
	k.Widget.Draw(dst.SubImage(k.Rectangle).(*image.RGBA), force)
}

/*
func (k *Kid) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return k.Widget.Mouse(pos.Sub(k.Rectangle.Min), but, dir, mod)
}
*/

func clear(dst *image.RGBA) {
	draw.Draw(dst, dst.Bounds(), Colors[1], image.ZP, draw.Src)
}
