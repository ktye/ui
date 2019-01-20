package ui

import (
	"image"
	"image/draw"
	"math"

	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/text"
	tui "github.com/eaburns/T/ui"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

type Edit struct {
	Target **Edit
	*tui.TextBox
	mods [4]bool
}

func (e *Edit) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	self.R = rect(sizeAvail)

	if e.Target != nil {
		*e.Target = e
	}
}

func (e *Edit) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	if e.TextBox == nil {
		win := tui.NewWin(w.Display.PixelsPerPt * 72.0)
		styles := [4]text.Style{
			text.Style{FG: w.Regular.Normal.Text, BG: w.Regular.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
		}
		e.TextBox = tui.NewTextBox(win, styles, self.R.Size())
		e.TextBox.SetText(rope.New("+/⍳100"))
	}
	subimage := img.(*image.RGBA).SubImage(self.R)
	e.TextBox.Draw(true, subimage.(draw.Image)) // TODO rect, dirty
}

func (e *Edit) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	e.mouseEvent(m.Event)
	r.Consumed = true
	self.Draw = Dirty // Let T decide when to draw. Not sure if that draws too often.
	return r
}

func (e *Edit) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	e.mods = e.keyEvent(e.mods, k)
	res.Consumed = true
	self.Draw = Dirty
	return res
}

func (e *Edit) FirstFocus(w *Window, self *Kid) (warp *image.Point) {
	// TODO lastCursorPoint
	return nil
}

func (e *Edit) Focus(w *Window, self *Kid, o Widget) (warp *image.Point) {
	if o != e {
		return nil
	}
	return e.FirstFocus(w, self)
}

func (e *Edit) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return self.Mark(o, forLayout)
}

// From T.main:

func (edt *Edit) mouseEvent(e mouse.Event) {
	w := edt.TextBox
	switch pt := image.Pt(int(e.X), int(e.Y)); {
	case e.Button == mouse.ButtonWheelUp:
		w.Wheel(pt, 0, 1)

	case e.Button == mouse.ButtonWheelDown:
		w.Wheel(pt, 0, -1)

	case e.Button == mouse.ButtonWheelLeft:
		w.Wheel(pt, -1, 0)

	case e.Button == mouse.ButtonWheelRight:
		w.Wheel(pt, 1, 0)

	case e.Direction == mouse.DirNone:
		w.Move(pt)

	case e.Direction == mouse.DirPress:
		w.Click(pt, int(e.Button))

	case e.Direction == mouse.DirRelease:
		w.Click(pt, -int(e.Button))

	case e.Direction == mouse.DirStep:
		w.Click(pt, int(e.Button))
		w.Click(pt, -int(e.Button))
	}
}

func (edt *Edit) keyEvent(mods [4]bool, e key.Event) [4]bool {
	w := edt.TextBox
	if e.Direction == key.DirNone {
		e.Direction = key.DirPress
	}
	if e.Direction == key.DirPress && dirKeyCode[e.Code] {
		edt.dirKey(e)
		return mods
	}

	switch {
	case e.Code == key.CodeDeleteBackspace:
		e.Rune = '\b'
	case e.Code == key.CodeDeleteForward:
		e.Rune = 0x7f
	case e.Rune == '\r':
		e.Rune = '\n'
	}
	if e.Rune > 0 {
		if e.Direction == key.DirPress {
			w.Rune(e.Rune)
		}
		return mods
	}

	return edt.modKey(mods, e)
}

var dirKeyCode = map[key.Code]bool{
	key.CodeUpArrow:    true,
	key.CodeDownArrow:  true,
	key.CodeLeftArrow:  true,
	key.CodeRightArrow: true,
	key.CodePageUp:     true,
	key.CodePageDown:   true,
	key.CodeHome:       true,
	key.CodeEnd:        true,
}

func (edt *Edit) dirKey(e key.Event) {
	w := edt.TextBox
	switch e.Code {
	case key.CodeUpArrow:
		w.Dir(0, -1)

	case key.CodeDownArrow:
		w.Dir(0, 1)

	case key.CodeLeftArrow:
		w.Dir(-1, 0)

	case key.CodeRightArrow:
		w.Dir(1, 0)

	case key.CodePageUp:
		w.Dir(0, -2)

	case key.CodePageDown:
		w.Dir(0, 2)

	case key.CodeHome:
		w.Dir(0, math.MinInt16)

	case key.CodeEnd:
		w.Dir(0, math.MaxInt16)

	default:
		panic("impossible")
	}
}

func (edt *Edit) modKey(mods [4]bool, e key.Event) [4]bool {
	w := edt.TextBox
	var newMods [4]bool
	if e.Modifiers&key.ModShift != 0 {
		newMods[1] = true
	}
	if e.Modifiers&key.ModAlt != 0 {
		newMods[2] = true
	}
	if e.Modifiers&key.ModMeta != 0 ||
		e.Modifiers&key.ModControl != 0 {
		newMods[3] = true
	}
	for i := 0; i < len(newMods); i++ {
		if newMods[i] != mods[i] {
			m := i
			if !newMods[i] {
				m = -m
			}
			w.Mod(m)
			mods = newMods
			break
		}
	}
	return mods
}
