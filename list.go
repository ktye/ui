package ui

import (
	"image"
	"image/draw"
	"sort"
	"time"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/syntax"
	"github.com/eaburns/T/text"
	"github.com/ktye/ui/tb"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type List struct {
	Target  **List
	Execute func()
	Quit    func()
	Delete  func()
	Single  bool
	text    rope.Rope
	*tb.TextBox
	mods     [4]bool
	styles   []text.Style
	sel      map[int]bool
	fontSize int
	line     map[int64]int
	addr     map[int]int64
	ct       time.Time
}

func (l *List) SetText(t rope.Rope) {
	l.text = t
	if l.TextBox != nil {
		l.TextBox.SetText(l.text)
	}
	l.line = make(map[int64]int)
	l.line[0] = 0
	end := t.Len()
	var sum int64 = 0
	ln := 1
	for {
		i := rope.IndexRune(t, '\n')
		if i < 0 {
			break
		}
		sum += i + 1
		l.line[sum] = ln
		ln++
		_, t = rope.Split(t, i+1)
	}
	l.line[end] = len(l.line)
	l.sel = make(map[int]bool)
	l.addr = make(map[int]int64)
	for i, k := range l.line {
		l.addr[k] = i
	}
}

func (l *List) Selection() []int {
	var sel []int
	for i, o := range l.sel {
		if o {
			sel = append(sel, i)
		} else {
			delete(l.sel, i)
		}
	}
	if len(sel) > 0 {
		sort.Ints(sel)
	}
	return sel
}

func (l *List) toggle(n int) {
	if l.sel[n] {
		delete(l.sel, n)
	} else if l.Single {
		if len(l.sel) > 0 {
			l.sel = make(map[int]bool)
		}
		l.sel[n] = true
	} else {
		l.sel[n] = true
	}
}

func (l *List) Update(h []syntax.Highlight, d edit.Diffs, t rope.Rope) []syntax.Highlight {
	sel := l.Selection()
	hi := make([]syntax.Highlight, len(sel))
	for i, n := range sel {
		hi[i].At[0] = l.addr[n]
		hi[i].At[1] = l.addr[n+1]
		hi[i].Style = l.styles[1]
	}
	return hi
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
			l.SetText(rope.New(""))
		} else {
			l.SetText(l.text)
		}
		l.SetHighlighter(l)
		l.TextBox.SetDoubleClick(func(x *tb.TextBox) {}) // ignore
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
	dot, dirty, sel, dbl := l.mouseEvent(m.Event)
	if sel {
		start := l.lineNum(dot[0])
		end := l.lineNum(dot[1])
		for i := start; i <= end; i++ {
			l.toggle(i)
		}
		dot[0] = dot[1]
		l.SetDot(dot)
		l.SetHighlighter(l)
	}
	if dbl {
		if l.Execute != nil {
			l.Execute()
		}
	}
	r.Consumed = true
	if dirty {
		self.Draw = Dirty
	}
	return r
}

func (l *List) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	res.Consumed = true
	if k.Direction == key.DirNone {
		k.Direction = key.DirPress
	}
	if k.Direction == key.DirPress && dirKeyCode[k.Code] {
		e := &Edit{TextBox: l.TextBox}
		e.dirKey(k)
		self.Draw = Dirty
		return res
	}
	if k.Direction == key.DirRelease {
		switch k.Code {
		case key.CodeReturnEnter:
			if l.Execute != nil {
				l.Execute()
			}
		case key.CodeSpacebar:
			d := l.TextBox.Dot()
			n := l.lineNum(d[0])
			l.toggle(n)
			l.SetHighlighter(l)
			self.Draw = Dirty
			return res

		case key.CodeEscape:
			if l.Quit != nil {
				l.Quit()
			}
		case key.CodeDeleteForward, key.CodeDeleteBackspace:
			if l.Delete != nil {
				l.Delete()
			}
		}
	}
	return res
}
func (l *List) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return self.Mark(o, forLayout)
}

func (l *List) mouseEvent(e mouse.Event) (d [2]int64, dirty, sel, dbl bool) {
	w := l.TextBox
	switch pt := image.Pt(int(e.X), int(e.Y)); {
	case e.Button == mouse.ButtonWheelUp:
		w.Wheel(pt, 0, 1)
		dirty = true
	case e.Button == mouse.ButtonWheelDown:
		w.Wheel(pt, 0, -1)
		dirty = true
	case e.Direction == mouse.DirPress:
		w.Click(pt, int(e.Button))
		dirty = true
	case e.Direction == mouse.DirRelease:
		_, d = w.Click(pt, -int(e.Button))
		if e.Button == 1 {
			dirty, sel = true, true
			if time.Now().Sub(l.ct) < 500*time.Millisecond {
				dbl = true
			}
			l.ct = time.Now()
		} else if e.Button == 3 {
			l.sel = make(map[int]bool)
			l.SetHighlighter(l)
			dirty = true
		}
	case e.Direction == mouse.DirNone:
		dirty = w.Move(pt)
	}
	return d, dirty, sel, dbl
}

func (l *List) lineNum(at int64) int {
	if l.text == nil {
		return 0
	}
	front, _ := rope.Split(l.text, at)
	start := rope.LastIndexFunc(front, func(r rune) bool { return r == '\n' })
	if start < 0 {
		start = 0
	} else {
		start++ // Don't include the \n.
	}
	return l.line[start]
}
