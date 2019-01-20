package ui

import (
	"image"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type Mouse struct {
	mouse.Event
}

func (m Mouse) In(r image.Rectangle) bool {
	return image.Point{int(m.X), int(m.Y)}.In(r)
}

func (m Mouse) Sub(r image.Point) Mouse {
	mm := m
	mm.X -= float32(r.X)
	mm.Y -= float32(r.Y)
	return mm
}

// Mouse delivers a mouse event to the UI tree.
// Mouse is typically called by Input.
func (w *Window) Mouse(m mouse.Event) {
	// Shift+WheelUp/Down: scale font size.
	if m.Modifiers&key.ModShift != 0 {
		if m.Button == mouse.ButtonWheelUp {
			w.scale(true)
			w.Top.Layout = Dirty
		} else if m.Button == mouse.ButtonWheelDown {
			w.scale(false)
			w.Top.Layout = Dirty
		}
	}
	if m.Button == 0 || w.origMouse.Button == 0 {
		w.origMouse = Mouse{m}
	}
	w.mouse = Mouse{m}
	r := w.Top.W.Mouse(w, &w.Top, w.mouse, w.origMouse, image.ZP)
	w.apply(r)
}
