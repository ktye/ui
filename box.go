package ui

import (
	"image"
	"image/draw"

	"golang.org/x/mobile/event/key"
)

// NewBox returns a box containing all widgets in its Kids field.
func NewBox(ws ...Widget) *Box {
	kids := make([]*Kid, len(ws))
	for i, w := range ws {
		kids[i] = &Kid{W: w}
	}
	return &Box{Kids: kids}
}

// NewReverseBox returns a box containing all widgets in original order in its Kids field, with the Reverse field set.
func NewReverseBox(ws ...Widget) *Box {
	kids := make([]*Kid, len(ws))
	for i, w := range ws {
		kids[i] = &Kid{W: w}
	}
	return &Box{Kids: kids, Reverse: true}
}

// Box keeps elements on a line as long as they fit, then moves on to the next line.
type Box struct {
	Kids       []*Kid      // Kids and UIs in this box.
	Reverse    bool        // Lay out children from bottom to top. First kid will be at the bottom.
	Margin     image.Point // In lowDPI pixels, will be adjusted for highDPI screens.
	Padding    Space       // Padding inside box, so children don't touch the sides; in lowDPI pixels, also adjusted for highDPI screens.
	Valign     Valign      // How to align children on a line.
	Width      int         // 0 means dynamic (as much as needed), -1 means full width, >0 means that exact amount of lowDPI pixels.
	Height     int         // 0 means dynamic (as much as needed), -1 means full height, >0 means that exact amount of lowDPI pixels.
	MaxWidth   int         // if >0, the max number of lowDPI pixels that will be used.
	Background image.Image
	size       image.Point // of entire box, including padding
}

func (b *Box) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	if KidsLayout(w, self, b.Kids, force) {
		return
	}

	if b.Width < 0 && b.MaxWidth > 0 {
		panic("combination ui.Width < 0 and ui.MaxWidth > 0 invalid")
	}

	osize := sizeAvail
	if b.Width > 0 && b.Width < sizeAvail.X {
		sizeAvail.X = b.Width
	} else if b.MaxWidth > 0 && b.MaxWidth < sizeAvail.X {
		// note: b.Width is currently the same as MaxWidth, but that might change when we don't mind extending beyong given X, eg with horizontal scroll
		sizeAvail.X = b.MaxWidth
	}
	if b.Height > 0 {
		sizeAvail.Y = b.Height
	}
	padding := b.Padding
	margin := b.Margin
	sizeAvail = sizeAvail.Sub(padding.Size())
	nx := 0 // number on current line

	// variables below are about box contents not offset for padding
	cur := image.ZP
	xmax := 0  // max x seen so far
	lineY := 0 // max y of current line

	fixValign := func(kids []*Kid) {
		if len(kids) < 2 {
			return
		}
		for _, k := range kids {
			switch b.Valign {
			case ValignTop:
			case ValignMiddle:
				k.R = k.R.Add(image.Pt(0, (lineY-k.R.Dy())/2))
			case ValignBottom:
				k.R = k.R.Add(image.Pt(0, lineY-k.R.Dy()))
			}
		}
	}

	for i, k := range b.Kids {
		k.W.Layout(w, k, sizeAvail.Sub(image.Pt(0, cur.Y+lineY)), true)
		childSize := k.R.Size()
		var kr image.Rectangle
		if nx == 0 || cur.X+childSize.X <= sizeAvail.X {
			kr = rect(childSize).Add(cur).Add(padding.Topleft())
			cur.X += childSize.X + margin.X
			lineY = maximum(lineY, childSize.Y)
			nx += 1
		} else {
			if nx > 0 {
				fixValign(b.Kids[i-nx : i])
				cur.X = 0
				cur.Y += lineY + margin.Y
			}
			kr = rect(childSize).Add(cur).Add(padding.Topleft())
			nx = 1
			cur.X = childSize.X + margin.X
			lineY = childSize.Y
		}
		k.R = kr
		if xmax < cur.X {
			xmax = cur.X
		}
	}
	fixValign(b.Kids[len(b.Kids)-nx : len(b.Kids)])
	cur.Y += lineY

	if b.Reverse {
		bottomY := cur.Y + padding.Dy()
		for _, k := range b.Kids {
			y1 := bottomY - k.R.Min.Y
			y0 := y1 - k.R.Dy()
			k.R = image.Rect(k.R.Min.X, y0, k.R.Max.X, y1)
		}
	}

	b.size = image.Pt(xmax-margin.X, cur.Y).Add(padding.Size())
	if b.Width < 0 {
		b.size.X = osize.X
	}
	if b.Height < 0 && b.size.Y < osize.Y {
		b.size.Y = osize.Y
	}
	self.R = rect(b.size)
}

func (b *Box) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	KidsDraw(w, self, b.Kids, b.size, b.Background, img, orig, m, force)
}

func (b *Box) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	return KidsMouse(w, self, b.Kids, m, origM, orig)
}

func (b *Box) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (r Result) {
	return KidsKey(w, self, b.orderedKids(), k, m, orig)
}

func (b *Box) orderedKids() []*Kid {
	if !b.Reverse {
		return b.Kids
	}
	n := len(b.Kids)
	kids := make([]*Kid, n)
	for i := range b.Kids {
		kids[i] = b.Kids[n-1-i]
	}
	return kids
}

func (b *Box) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return KidsMark(self, b.Kids, o, forLayout)
}
