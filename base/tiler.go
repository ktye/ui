package base

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
)

// Tiler is a widget container with a draggable border.
// If a border is dragged, the tiler's child widget is replaced by a splitter
// containing the child and a new widget created by the New function.
type Tiler struct {
	Widget ui.Widget
	New    func() ui.Widget
	rect   image.Rectangle
	drag   int
	start  bool
	sp     image.Point
	store  *image.RGBA
}

func NewTiler(child ui.Widget, f func() ui.Widget) *Tiler {
	return &Tiler{Widget: child, New: f}
}
func (t *Tiler) Draw(dst *image.RGBA, force bool) {
	if t.rect != dst.Rect {
		t.rect = dst.Rect
	}
	if t.start {
		t.start = false
		t.store = image.NewRGBA(dst.Rect)
		draw.Draw(t.store, dst.Rect, dst, dst.Rect.Min, draw.Src)
	}
	t.Widget.Draw(dst, force)
}
func (t *Tiler) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if t.drag > 0 {
		if but == 1 && dir < 0 {
			t.drag = 0
		} else if but == 0 && dir == 0 {
			d := [4]int{pos.X - t.sp.X, t.sp.X - pos.X, pos.Y - t.sp.Y, t.sp.Y - pos.Y}
			if d[t.drag-1] > 10 {
				a, b := t.Widget, t.New()
				if t.drag == 1 || t.drag == 3 {
					a, b = b, a
				}
				s := NewSplit(a, b)
				if t.drag < 3 {
					s.Vertical = true
				}
				t.drag = 0
				s.gutterClick = t.sp
				s.drag = true
				s.store = t.store
				t.Widget = s
				return -1
			}
		}
	}
	r := 0
	if but == 1 && dir > 0 {
		near := func(x int) bool {
			return x > 0 && x < 10
		}
		t.sp = pos
		t.start = true
		switch {
		case near(pos.X - t.rect.Min.X):
			t.drag = 1
		case near(t.rect.Max.X - pos.X):
			t.drag = 2
		case near(pos.Y - t.rect.Min.Y):
			t.drag = 3
		case near(t.rect.Max.Y - pos.Y):
			t.drag = 4
		default:
			t.start = false
		}
		r = 1
	}
	if res := t.Widget.Mouse(pos, but, dir, mod); res == 0 {
		return r
	} else {
		return res
	}
}
func (t *Tiler) Key(r rune, code uint32, dir int, mod uint32) int {
	return t.Widget.Key(r, code, dir, mod)
}
