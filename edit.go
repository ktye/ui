package ui

/* TODO port to v2

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"strings"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/re1"
	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/text"
	"github.com/ktye/ui/tb"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
)

// Edit is a text editor widget.
type Edit struct {
	Target  **Edit
	Execute func(*Edit, string)
	Nowrap  bool
	text    rope.Rope
	*tb.TextBox
	mods     [4]bool
	styles   []text.Style
	fontSize int
}

func (e *Edit) SetText(text rope.Rope) {
	e.text = text
	if e.TextBox != nil {
		e.TextBox.SetText(e.text)
	}
}

// AppendText appends the string to the end.
// If there is no newline at the end, it first inserts a newline.
func (e *Edit) AppendText(s string) {
	s = strings.Replace(s, `\`, `\\`, -1)
	s = strings.Replace(s, "\n", `\n`, -1)
	s = strings.Replace(s, "\t", `\t`, -1)
	s = strings.Replace(s, "/", `\/`, -1)
	_, err := e.Edit(`/\n$/`)
	if _, ok := err.(edit.NoCommandError); ok == false {
		s = "\\n" + s // No newline at EOF
	}
	t := "$a/" + s + "/"
	_, err = e.Edit(t)
	if err != nil {
		fmt.Printf("ui/edit: %q %s\n", t, err)
	}
}

// MarkAddr sets the current dot to the address given as an edit command.
func (e *Edit) MarkAddr(addr string) error {
	_, err := e.TextBox.Edit(addr)
	if ae, ok := err.(edit.NoCommandError); ok {
		e.TextBox.SetDot(ae.At)
		return nil
	} else {
		if err == nil {
			return fmt.Errorf("addr is more than an address")
		}
		return err
	}
}

func (e *Edit) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	self.R = rect(sizeAvail)
	if e.Target != nil {
		*e.Target = e
	}
}

func (e *Edit) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {

	if e.TextBox == nil {
		styles := [4]text.Style{
			text.Style{FG: w.Regular.Normal.Text, BG: w.Regular.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
			text.Style{FG: w.Primary.Normal.Text, BG: w.Primary.Normal.Background.At(0, 0), Face: w.font.Face},
		}
		if len(e.styles) > 0 && len(e.styles) <= 4 {
			for i, s := range e.styles {
				styles[i] = s
			}
		}

		e.TextBox = tb.NewTextBox(styles, self.R.Size())
		e.TextBox.SetNowrap(e.Nowrap)
		if e.text == nil {
			e.TextBox.SetText(rope.New(""))
		} else {
			e.TextBox.SetText(e.text)
		}
	}

	// We track the last font size update in e.fontSize.
	if w.font.size != e.fontSize {
		e.TextBox.SetFace(w.font.Face)
		e.fontSize = w.font.size
	}

	rect := self.R.Add(orig).Sub(self.R.Min)
	//rect.Max = rect.Max.Add(image.Point{1000, 0}) //
	subimage := img.(*image.RGBA).SubImage(rect)
	e.TextBox.Draw(true, subimage.(draw.Image))
}

func (e *Edit) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	b, dot, dirty := e.mouseEvent(m.Event)

	// Search on Button-3 release.
	// If Button-3 marked text, it is used, otherwise the current dot is used.
	if b == -2 || b == -3 {
		if dot[0] != dot[1] {
			e.TextBox.SetDot(dot)
		} else {
			dot = e.TextBox.Dot()
		}
		sel := e.TextBox.Selection()
		if len(sel) > 0 {
			if b == -3 {
				// Search on Button 3
				t := e.Escape(sel)
				t = "+/" + t + "/"
				if err := e.MarkAddr(t); err != nil {
					fmt.Printf("ui/edit: %s %s", t, err)
				}
			} else if b == -2 && e.Execute != nil {
				// Execute on Button 2
				e.Execute(e, sel)
			}
		}
		dirty = true
	}

	r.Consumed = true
	if dirty {
		self.Draw = Dirty
	}
	return r
}

// Escape a string that it can be inserted/appended with "a/"+string"+"/".
func (e *Edit) Escape(t string) string {
	t = re1.Escape(t)
	t = strings.Replace(t, "\n", `\n`, -1)
	t = strings.Replace(t, "/", `\/`, -1)
	return t
}

func (e *Edit) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	e.mods = e.keyEvent(w, e.mods, k)
	res.Consumed = true
	self.Draw = Dirty
	return res
}

func (e *Edit) Mark(self *Kid, o Widget, forLayout bool) (marked bool) {
	return self.Mark(o, forLayout)
}

// From T.main:

func (edt *Edit) mouseEvent(e mouse.Event) (int, [2]int64, bool) {
	w := edt.TextBox
	draw := true
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
		// Only draw if there was a change.
		draw = w.Move(pt)

	case e.Direction == mouse.DirPress:
		w.Click(pt, int(e.Button))

	case e.Direction == mouse.DirRelease:
		b, d := w.Click(pt, -int(e.Button))
		return b, d, true

	case e.Direction == mouse.DirStep:
		w.Click(pt, int(e.Button))
		w.Click(pt, -int(e.Button))
	}
	return 0, [2]int64{0, 0}, draw
}

func (edt *Edit) keyEvent(w *Window, mods [4]bool, e key.Event) [4]bool {
	if e.Direction == key.DirNone {
		e.Direction = key.DirPress
	}
	if e.Direction == key.DirPress && dirKeyCode[e.Code] {
		edt.dirKey(e)
		return mods
	}

	// Copy/Paste/Cut: Cntrl-C/V/X
	if e.Code == key.CodeC && mods[3] {
		if e.Direction == key.DirPress {
			edt.TextBox.Copy(w.clipboard)
		}
		return mods
	}
	if e.Code == key.CodeV && mods[3] {
		if e.Direction == key.DirPress {
			edt.TextBox.Paste(w.clipboard)
		}
		return mods
	}
	if e.Code == key.CodeX && mods[3] {
		if e.Direction == key.DirPress {
			edt.TextBox.Cut(w.clipboard)
		}
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
			edt.TextBox.Rune(e.Rune)
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
*/
