package base

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
)

type Split struct {
	Target **Split
	Base
	Vertical    bool
	Ratio       float64 // (a-b)/(a+b)
	Kids        [2]Kid
	size        int
	rect        image.Rectangle
	drag        bool
	other       bool
	gutterClick image.Point
}

func NewSplit(kid1, kid2 ui.Widget) *Split {
	var s Split
	s.Kids = [2]Kid{Kid{Widget: kid1}, Kid{Widget: kid2}}
	return &s
}
func (s *Split) Draw(dst *image.RGBA, force bool) {
	if s.Target != nil {
		*s.Target = s
	}
	s.Rect = dst.Rect
	if force == false && s.draw == false {
		s.Kids[0].Draw(dst, force)
		s.Kids[1].Draw(dst, force)
		return
	}
	s.draw = false
	s.size = dst.Rect.Dy()
	if s.Vertical {
		s.size = dst.Rect.Dx()
	}
	a := (s.size + int(s.Ratio*float64(s.size))) / 2
	if a < 0 {
		a = 0
	}
	b := s.size - a
	if b < 0 {
		b = 0
	}
	if s.Vertical {
		s.Kids[0].Rectangle = dst.Rect.Intersect(dst.Rect.Sub(image.Point{b + 1, 0}))
		s.Kids[1].Rectangle = dst.Rect.Intersect(dst.Rect.Add(image.Point{a + 1, 0}))
		s.rect = image.Rect(dst.Rect.Min.X+a-1, dst.Rect.Min.Y, dst.Rect.Min.X+a+1, dst.Rect.Max.Y)
	} else {
		s.Kids[0].Rectangle = dst.Rect.Intersect(dst.Rect.Sub(image.Point{0, b + 1}))
		s.Kids[1].Rectangle = dst.Rect.Intersect(dst.Rect.Add(image.Point{0, a + 1}))
		s.rect = image.Rect(dst.Rect.Min.X, dst.Rect.Min.Y+a-1, dst.Rect.Max.X, dst.Rect.Min.Y+a+1)
	}
	s.Kids[0].Draw(dst, true)
	s.Kids[1].Draw(dst, true)
	c := SplitGutter
	if s.drag {
		c = SplitGutterActive
	}
	draw.Draw(dst, s.rect, c.Uniform(), image.ZP, draw.Src)
}
func (s *Split) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	gutter := s.rect
	gutter.Min = gutter.Min.Sub(image.Point{2, 2})
	gutter.Max = gutter.Max.Add(image.Point{2, 2})
	if pos.In(gutter) && but == 1 && dir > 0 { // but-1 click in gutter → start drag or flip
		s.gutterClick = pos
		s.drag = true
		return s.DrawSelf()
	} else if but == 1 && dir < 0 { // but-1 release: stop drag/flip
		if s.drag {
			s.drag = false
			if pos == s.gutterClick {
				s.Vertical = !s.Vertical
			}
			return s.DrawSelf()
		}
	} else if dir == 0 && s.drag { // gutter drag
		a := pos.Y - s.Rect.Min.Y
		if s.Vertical {
			a = pos.X - s.Rect.Min.X
		}
		if s.size > 0 {
			s.Ratio = float64(2*a-s.size) / float64(s.size)
		} else {
			s.Ratio = 0
		}
		return s.DrawSelf()
	}
	if but == 1 && dir > 0 { // but-1 click in kid: activate kid for key events
		s.other = false
		if pos.In(s.Kids[1].Rectangle) {
			s.other = true
		}
	}
	if pos.In(s.Kids[0].Rectangle) {
		return s.Kids[0].Mouse(pos, but, dir, mod)
	} else if pos.In(s.Kids[1].Rectangle) {
		return s.Kids[1].Mouse(pos, but, dir, mod)
	}
	return 0
}
func (s *Split) Key(r rune, code uint32, dir int, mod uint32) int {
	k := &s.Kids[0]
	if s.other {
		k = &s.Kids[1]
	}
	return k.Key(r, code, dir, mod)
}
