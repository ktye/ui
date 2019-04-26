package base

/* TODO port to v2

import (
	"image"
	"image/draw"

	"github.com/ktye/ui/paint"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type Button struct {
	Text     string
	Icon     string
	Disabled bool
	Colorset *Colorset
	Click    func() (e Event)

	//m    Mouse
	size image.Point
}

func (b *Button) space(w *Window) image.Point {
	return b.padding(w).Add(pt(BorderSize))
}

func (b *Button) padding(w *Window) image.Point {
	h := w.FontHeight()
	return image.Pt(h/2, h/4)
}

func (b *Button) iconSize(w *Window) int {
	if b.Icon == "" {
		return 0
	}
	h := w.FontHeight()
	if b.Text == "" {
		return 3 * h
	}
	return h
}

func (b *Button) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	size := b.innerSize(w).Add(b.space(w).Mul(2))
	b.size = size
	self.R = rect(size)
}

func (b *Button) innerSize(w *Window) image.Point {
	if len(b.Text) > 0 && len(b.Icon) > 0 {
		h := w.FontHeight()
		return w.StringSize(b.Text).Add(image.Point{3 * h / 2, 0})
	} else if len(b.Text) > 0 {
		return w.StringSize(b.Text)
	} else {
		h := b.iconSize(w)
		return image.Point{h, h}
	}
}

func (b *Button) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	r := rect(b.innerSize(w).Add(b.space(w).Mul(2)))
	hover := m.In(r)
	var colors Colors
	if b.Disabled {
		colors = w.Disabled
	} else {
		cs := b.Colorset
		if cs == nil {
			cs = &w.Regular
		}
		colors = cs.Normal
		if hover {
			colors = cs.Hover
		}
	}

	r = r.Add(orig)
	draw.Draw(img, r.Inset(1), colors.Background, image.ZP, draw.Src)
	paint.RoundedBorder(img, r, colors.Border)

	hit := image.ZP
	if hover && !b.Disabled && m.Button == 1 && m.Direction == mouse.DirPress {
		hit = image.Pt(0, 1)
	}
	p := r.Min.Add(b.space(w)).Add(hit)
	if b.Icon != "" {
		size := b.iconSize(w)
		ico := getIcon(w, icon{b.Icon, size})
		draw.DrawMask(img, r.Bounds(), &image.Uniform{colors.Text}, image.ZP, ico, r.Min.Sub(p), draw.Over)
		p = p.Add(image.Point{3 * size / 2, 0})
	}
	if b.Text != "" {
		paint.String(img, p, colors.Text, w.font, b.Text)
	}
}

func (b *Button) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	if b.Disabled {
		return
	}
	rr := rect(b.size)
	hover := m.In(rr)
	if hover && m.Button == 1 && m.Direction == mouse.DirPress {
		self.Draw = Dirty
	}
	if hover && m.Button == 1 && m.Direction == mouse.DirRelease && b.Click != nil {
		self.Draw = Dirty
		e := b.Click()
		propagateEvent(self, &r, e)
	}
	return r
}

func (b *Button) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (r Result) {
	if !b.Disabled && (k.Rune == ' ' || k.Rune == '\n') && k.Direction == key.DirRelease {
		r.Consumed = true
		if b.Click != nil {
			e := b.Click()
			propagateEvent(self, &r, e)
		}
	}
	return
}

func (b *Button) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return self.Mark(o, forLayout)
}
*/
