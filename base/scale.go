package base

import (
	"image"

	"github.com/ktye/ui"
)

// Scale is a widget that scales the global font size on mouse wheel events.
// It can be used to wrap the top widget.
type Scale struct {
	ui.Widget
	Funcs []func() // custom notification after font reset
}

func NewScale(w ui.Widget) Scale {
	return Scale{Widget: w}
}
func (s Scale) Draw(dst *image.RGBA, force bool) {
	s.Widget.Draw(dst, force)
}
func (s Scale) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if mod&1 != 0 {
		if but == -1 {
			scale(true)
			s.notify()
			return -1
		} else if but == -2 {
			scale(false)
			s.notify()
			return -1
		}
	}
	return s.Widget.Mouse(pos, but, dir, mod)
}
func (s Scale) Key(r rune, code uint32, dir int, mod uint32) int {
	return s.Widget.Key(r, code, dir, mod)
}

func (s Scale) notify() {
	if s.Funcs != nil {
		for _, f := range s.Funcs {
			f()
		}
	}
}

// Scale is the callback for shift-wheel up/down.
// It changes the font size.
func scale(up bool) {
	if Font.TTF == nil {
		return
	}
	if up {
		Font.size++
	} else if Font.size < 5 {
		return
	} else {
		Font.size--
	}
	SetFont(Font.TTF, Font.size)
}
