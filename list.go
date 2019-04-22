package ui

import (
	"image"
	"image/draw"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/syntax"
	"github.com/eaburns/T/text"
	"github.com/ktye/ui/tb"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type List struct {
	Target **List
	text   rope.Rope
	*tb.TextBox
	mods     [4]bool
	styles   []text.Style
	sel      []int
	fontSize int
}

func (l *List) SetText(t rope.Rope) {
	l.text = t
	if l.TextBox != nil {
		l.TextBox.SetText(l.text)
	}

}

func (l *List) Update(h []syntax.Highlight, d edit.Diffs, t rope.Rope) []syntax.Highlight {
	println("hilight", len(h), t.String())
	return []syntax.Highlight{
		syntax.Highlight{
			At:    [2]int64{2, 4},
			Style: l.styles[1],
		},
	}
}

func (l *List) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	self.R = rect(sizeAvail)
	if l.Target != nil {
		*l.Target = l
	}
}

func (l *List) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	if l.TextBox == nil {
		styles := [4]text.Style{
			text.Style{FG: w.Regular.Normal.Text, BG: w.Regular.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
		}
		l.styles = styles[:]
		l.TextBox = tb.NewTextBox(styles, self.R.Size())
		l.TextBox.SetNowrap(true)
		if l.text == nil {
			l.TextBox.SetText(rope.New(""))
		} else {
			l.TextBox.SetText(l.text)
		}
		l.SetHighlighter(l)
	}
	if w.font.size != l.fontSize {
		l.TextBox.SetFace(w.font.Face)
		l.fontSize = w.font.size
	}

	rect := self.R.Add(orig).Sub(self.R.Min)
	subimage := img.(*image.RGBA).SubImage(rect)
	l.TextBox.Draw(true, subimage.(draw.Image))
}

func (l *List) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	dot, ok := l.mouseEvent(m.Event)
	if ok {
		println("dot", dot[0], dot[1])
		span0 := l.lineSpan(dot[0])
		span1 := l.lineSpan(dot[1])
		println("span", span0[0], span0[1], span1[0], span1[1])
	}
	r.Consumed = true
	return r
}

func (l *List) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	//e.mods = e.keyEvent(w, e.mods, k)
	//res.Consumed = true
	//self.Draw = Dirty
	return res
}
func (l *List) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return self.Mark(o, forLayout)
}

func (l *List) mouseEvent(e mouse.Event) (d [2]int64, ok bool) {
	w := l.TextBox
	switch pt := image.Pt(int(e.X), int(e.Y)); {
	case e.Button == mouse.ButtonWheelUp:
		w.Wheel(pt, 0, 1)
	case e.Button == mouse.ButtonWheelDown:
		w.Wheel(pt, 0, -1)
	case e.Direction == mouse.DirPress:
		w.Click(pt, int(e.Button))
		return d, false
	case e.Direction == mouse.DirRelease:
		_, dot := w.Click(pt, -int(e.Button))
		return dot, true
	case e.Direction == mouse.DirNone:
		w.Move(pt)
	}
	return d, false
}

func (l *List) lineSpan(at int64) [2]int64 { // line start and end at addr.
	if l.text == nil {
		return [2]int64{0, 0}
	}
	front, back := rope.Split(l.text, at)
	start := rope.LastIndexFunc(front, func(r rune) bool { return r == '\n' })
	if start < 0 {
		start = 0
	} else {
		start++ // Don't include the \n.
	}
	end := rope.IndexFunc(back, func(r rune) bool { return r == '\n' })
	if end < 0 {
		end = l.text.Len()
	} else {
		end += at + 1 // Do include the \n.
	}
	return [2]int64{start, end}
}
