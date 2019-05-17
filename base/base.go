package base

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
)

type Base struct {
	Menu  *Menu
	Rect  image.Rectangle
	Dirty bool
}

func (b *Base) DrawSelf() int {
	b.Dirty = true
	return 1
}

func (b *Base) ShowMenu(p image.Point) int {
	if b.Menu != nil {
		b.Menu.active = true
		b.Menu.point = p
		return 1
	}
	return 0
}
func (b *Base) DrawMenu(dst *image.RGBA) {
	if b.Menu != nil {
		b.Menu.drawMenu(dst)
	}
}
func (b *Base) MenuMouse(pos image.Point, but int, dir int, mod uint32) (int, bool) {
	if b.Menu != nil && b.Menu.active {
		r, o := b.Menu.mouse(pos, but, dir, mod), true
		if o && r != 0 {
			b.Dirty = true
		}
		return r, o
	}
	return 0, false
}
func (b *Base) MenuKey(r rune, code uint32, dir int, mod uint32) (int, bool) {
	if b.Menu != nil && b.Menu.active {
		r, o := b.Menu.key(r, code, dir, mod), true
		if o && r != 0 {
			b.Dirty = true
		}
		return r, o
	}
	return 0, false
}

type Kid struct {
	ui.Widget
	image.Rectangle
}

func (k *Kid) Draw(dst *image.RGBA, force bool) {
	k.Widget.Draw(dst.SubImage(k.Rectangle).(*image.RGBA), force)
}

func clear(dst *image.RGBA) {
	draw.Draw(dst, dst.Bounds(), Colors[1], image.ZP, draw.Src)
}

type Blank struct {
	Base
}

func (b *Blank) Draw(dst *image.RGBA, force bool) {
	if force {
		clear(dst)
	}
}
func (b *Blank) Mouse(pos image.Point, but int, dir int, mod uint32) int { return 0 }
func (b *Blank) Key(r rune, code uint32, dir int, mod uint32) int        { return 0 }

// Errlog can be overwritten with a with a user widget that presents errors.
var Errlog = func(err error) {
	println(err.Error())
}
