package ui

import (
	"image"
	"image/draw"

	"golang.org/x/mobile/event/key"
)

// Split is a horizontal or vertical split of the available space, for 2 widgets.
type Split struct {
	// Space between the widgets in pixels.
	// If >0, users can drag the gutter.
	Gutter   int
	Vertical bool
	Ratio    float64 // (a-b)/(a+b)
	Kids     []*Kid  // 2 Kids that should be split.
	size     image.Point
	sizes    [2]int
}

func (s *Split) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	if KidsLayout(w, self, s.Kids, force) {
		return
	}
	self.R = rect(sizeAvail)

	size := sizeAvail.X
	if s.Vertical {
		size = sizeAvail.Y
	}
	if size > s.Gutter {
		size -= s.Gutter
	}
	a := (size + int(s.Ratio*float64(size))) / 2
	if a < 0 {
		a = 0
	}
	b := size - a
	sz := [2]int{a, b}
	s.sizes = sz
	if s.Vertical {
		s.Kids[0].W.Layout(w, s.Kids[0], image.Pt(sizeAvail.X, sz[0]), true)
		s.Kids[1].W.Layout(w, s.Kids[1], image.Pt(sizeAvail.X, sz[1]), true)
		s.Kids[1].R = s.Kids[1].R.Add(image.Pt(0, sz[0]+s.Gutter))
	} else {
		s.Kids[0].W.Layout(w, s.Kids[0], image.Pt(sz[0], sizeAvail.Y), true)
		s.Kids[1].W.Layout(w, s.Kids[1], image.Pt(sz[1], sizeAvail.Y), true)
		s.Kids[1].R = s.Kids[1].R.Add(image.Pt(sz[0]+s.Gutter, 0))
	}
	s.size = sizeAvail
}

func (s *Split) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	pt := orig.Add(image.Point{s.sizes[0], 0})
	gutter := image.Rectangle{pt, pt.Add(image.Point{s.Gutter, s.size.Y})}
	if s.Vertical {
		pt := orig.Add(image.Point{0, s.sizes[0]})
		gutter = image.Rectangle{pt, pt.Add(image.Point{s.size.X, s.Gutter})}
	}
	draw.Draw(img, gutter, &image.Uniform{w.Gutter}, image.ZP, draw.Src)
	KidsDraw(w, self, s.Kids, s.size, nil, img, orig, m, force)
}

func (s *Split) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	return KidsMouse(w, self, s.Kids, m, origM, orig)
}

func (s *Split) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (r Result) {
	return KidsKey(w, self, s.Kids, k, m, orig)
}

func (s *Split) FirstFocus(w *Window, self *Kid) *image.Point {
	return KidsFirstFocus(w, self, s.Kids)
}

func (s *Split) Focus(w *Window, self *Kid, o Widget) *image.Point {
	return KidsFocus(w, self, s.Kids, o)
}

func (s *Split) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return KidsMark(self, s.Kids, o, forLayout)
}
