package ui

import (
	"image"
	"image/draw"

	"golang.org/x/mobile/event/key"
)

// Kid holds a UI and its layout/draw state.
type Kid struct {
	W      Widget          // widget this state is about.
	R      image.Rectangle // Location and size within this widget.
	Draw   State           // Whether widget or its children need a draw.
	Layout State           // Whether widget or its children need a layout.
}

// NewKids turns widgets into Kids containing those widgets.
func NewKids(uis ...Widget) []*Kid {
	kids := make([]*Kid, len(uis))
	for i, ui := range uis {
		kids[i] = &Kid{W: ui}
	}
	return kids
}

// Mark checks if o is its Widget, and if so marks it as needing a layout or draw (forLayout false).
func (k *Kid) Mark(o Widget, forLayout bool) (marked bool) {
	if o != k.W {
		return false
	}
	if forLayout {
		k.Layout = Dirty
	} else {
		k.Draw = Dirty
	}
	return true
}

func propagateEvent(self *Kid, r *Result, e Event) {
	if e.NeedLayout {
		self.Layout = Dirty
	}
	if e.NeedDraw {
		self.Draw = Dirty
	}
	r.Consumed = e.Consumed || r.Consumed
}

func propagateResult(w *Window, self, k *Kid) {
	if k.Layout != Clean {
		if k.Layout == DirtyKid {
			// panic("kid propagated layout kids")
			k.Layout = Dirty // xxx
		}
		nk := *k
		k.W.Layout(w, &nk, k.R.Size(), false)
		if nk.R.Size() != k.R.Size() {
			self.Layout = Dirty
		} else {
			self.Layout = Clean
			k.Layout = Clean
			nk.R = nk.R.Add(k.R.Min)
			k.Draw = Dirty
			self.Draw = DirtyKid
		}
	} else if k.Draw != Clean {
		self.Draw = DirtyKid
	}
}

// KidsLayout is called by layout widgets before they do their own layouts.
// KidsLayout returns whether there is any work left to do, determined by looking at self.Layout.
// Children will be layed out if necessary. KidsLayout updates layout and draw state of self and kids.
func KidsLayout(w *Window, self *Kid, kids []*Kid, force bool) (done bool) {
	if force {
		self.Layout = Clean
		self.Draw = Dirty
		return false
	}
	switch self.Layout {
	case Clean:
		return true
	case Dirty:
		self.Layout = Clean
		self.Draw = Dirty
		return false
	}
	for _, k := range kids {
		if k.Layout == Clean {
			continue
		}
		k.W.Layout(w, k, k.R.Size(), false)
		switch k.Layout {
		case Dirty:
			self.Layout = Dirty
			self.Draw = Dirty
			return false
		case DirtyKid:
			panic("layout of kid results in kid.Layout = DirtKid")
		case Clean:
		}
	}
	self.Layout = Clean
	self.Draw = Dirty
	return true
}

// KidsDraw draws a widget by drawing all its kids.
// uiSize is the size of the entire UI, used in case it has to be redrawn entirely.
// Bg can override the default duit background color.
// Img is the whether the widget should be drawn on, with orig as origin (offset).
// M is used for passing a mouse position to the kid's widget draw, for possibly drawing hover states.
// KidsDraw only draws if draw state indicates a need for drawing, or if force is set.
func KidsDraw(w *Window, self *Kid, kids []*Kid, size image.Point, bg image.Image, img draw.Image, orig image.Point, m Mouse, force bool) {
	force = force || self.Draw == Dirty
	if force {
		self.Draw = Dirty
	}

	if bg == nil {
		bg = w.Background
	}
	if force {
		draw.Draw(img, rect(size).Add(orig), bg, image.ZP, draw.Src)
	}
	for _, k := range kids {
		if !force && k.Draw == Clean {
			continue
		}
		if !force && k.Draw == Dirty {
			draw.Draw(img, k.R.Add(orig), bg, image.ZP, draw.Src)
		}

		mm := m
		mm = mm.Sub(k.R.Min)
		if force {
			k.Draw = Dirty
		}
		k.W.Draw(w, k, img, orig.Add(k.R.Min), mm, force)
		k.Draw = Clean
	}
	self.Draw = Clean
}

// KidsMouse delivers mouse event m to the widget at origM (often the same, but not in case button is held pressed).
// Mouse positions are always relative to their own origin. Orig is passed so UIs can calculate locations to warp the mouse to.
func KidsMouse(w *Window, self *Kid, kids []*Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	for _, k := range kids {
		if !origM.In(k.R) {
			continue
		}
		origM = origM.Sub(k.R.Min)
		m = m.Sub(k.R.Min)
		r = k.W.Mouse(w, k, m, origM, orig.Add(k.R.Min))
		if r.Hit == nil {
			r.Hit = k.W
		}
		propagateResult(w, self, k)
		return
	}
	return Result{}
}

// KidsKey delivers key event key to the widget at m.
// Orig is passed so widgets can calculate locations to warp the mouse to.
func KidsKey(w *Window, self *Kid, kids []*Kid, key key.Event, m Mouse, orig image.Point) (r Result) {
	for i, k := range kids {
		if !m.In(k.R) {
			continue
		}
		m = m.Sub(k.R.Min)
		r = k.W.Key(w, k, key, m, orig.Add(k.R.Min))
		if !r.Consumed && key.Rune == '\t' { // TODO: direction?
			for next := i + 1; next < len(kids); next++ {
				k := kids[next]
				first := k.W.FirstFocus(w, k)
				if first != nil {
					p := first.Add(orig).Add(k.R.Min)
					r.Warp = &p
					r.Consumed = true
					r.Hit = k.W
					break
				}
			}
		}
		if r.Hit == nil {
			r.Hit = self.W
		}
		propagateResult(w, self, k)
		return
	}
	return Result{}
}

// KidsFirstFocus delivers the FirstFocus request to the first leaf widget, and returns the location where the mouse should warp to.
func KidsFirstFocus(w *Window, self *Kid, kids []*Kid) *image.Point {
	if len(kids) == 0 {
		return nil
	}
	for _, k := range kids {
		first := k.W.FirstFocus(w, k)
		if first != nil {
			p := first.Add(k.R.Min)
			return &p
		}
	}
	return nil
}

// KidsFocus delivers the Focus request to the first leaf widget, and returns the location where the mouse should warp to.
func KidsFocus(w *Window, self *Kid, kids []*Kid, widget Widget) *image.Point {
	if len(kids) == 0 {
		return nil
	}
	for _, k := range kids {
		p := k.W.Focus(w, k, widget)
		if p != nil {
			pp := p.Add(k.R.Min)
			return &pp
		}
	}
	return nil
}

// KidsMark finds o in this UI subtree (self and kids), marks it as needing layout or draw (forLayout false), and returns whether it found and marked the UI.
func KidsMark(self *Kid, kids []*Kid, o Widget, forLayout bool) (marked bool) {
	if self.Mark(o, forLayout) {
		return true
	}
	for _, k := range kids {
		marked = k.W.Mark(k, o, forLayout)
		if !marked {
			continue
		}
		if forLayout {
			if self.Layout == Clean {
				self.Layout = DirtyKid
			}
		} else {
			if self.Draw == Clean {
				self.Draw = DirtyKid
			}
		}
		return true
	}
	return false
}
