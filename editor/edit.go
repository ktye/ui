package editor

import (
	"fmt"
	"image"
	"math"
	"strings"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/re1"
	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/text"
	"github.com/ktye/ui/base"
)

// Edit is a text editor widget for ui.
type Edit struct {
	base.Base
	Target   **Edit
	Execute  func(*Edit, string) int
	Nowrap   bool
	text     rope.Rope
	mods     [4]bool
	styles   []text.Style
	fontSize int
	draw     bool
	*TextBox
}

// New returns a text editor widget with the content from a string.
// To edit a rope, use SetText instead.
func New(s string) *Edit {
	e := Edit{}
	e.SetText(rope.New(s))
	return &e
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

func (e *Edit) initText(size image.Point) {
	styles := [4]text.Style{
		text.Style{FG: base.Black, BG: base.White, Face: base.Font.Face},
		text.Style{FG: base.White, BG: base.Blue, Face: base.Font.Face},
		text.Style{FG: base.Black, BG: base.Green, Face: base.Font.Face},
		text.Style{FG: base.Black, BG: base.Red, Face: base.Font.Face},
	}
	if len(e.styles) > 0 && len(e.styles) <= 4 {
		for i, s := range e.styles {
			styles[i] = s
		}
	}
	e.TextBox = NewTextBox(styles, size)
	e.TextBox.SetNowrap(e.Nowrap)
	if e.text == nil {
		e.TextBox.SetText(rope.New(""))
	} else {
		e.TextBox.SetText(e.text)
	}
}

func (e *Edit) Draw(dst *image.RGBA, force bool) {
	if e.Target != nil {
		*e.Target = e
	}
	if e.TextBox == nil {
		e.initText(dst.Rect.Size())
	}
	if fs := base.Font.Size(); fs != e.fontSize {
		e.TextBox.SetFace(base.Font.Face)
		e.fontSize = fs
	}
	if force || dst.Bounds() != e.Base.Rect || e.Dirty {
		e.Base.Rect = dst.Bounds()
		e.TextBox.Draw(true, dst)
		e.Dirty = false
	}
	e.DrawMenu(dst)
}

// Mouse and Key events are adapted from eaburns/T/main.go

func (e *Edit) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if r, o := e.MenuMouse(pos, but, dir, mod); o {
		return r
	}
	if but == 3 {
		if dir >= 0 {
			return 0
		}
		if dir < 0 {
			return e.ShowMenu(pos)
		}
	}

	draw := 0
	b, dot, dirty := mouseEvent(e.TextBox, pos.Sub(e.Rect.Min), but, dir)
	if dirty {
		draw = 1
	}

	if b == -2 && e.Execute != nil {
		if dot[0] != dot[1] {
			e.TextBox.SetDot(dot)
		} else {
			dot = e.TextBox.Dot()
		}
		if d := e.Execute(e, e.TextBox.Selection()); d != 0 {
			draw = d
		}
		if draw == 0 {
			draw = 1
		}
	}
	if draw != 0 {
		e.Dirty = true
	}
	return draw
}
func (e *Edit) Find() {
	t := e.Escape(e.TextBox.Selection())
	t = "+/" + t + "/"
	if err := e.MarkAddr(t); err != nil {
		fmt.Printf("ui/edit: %s %s", t, err)
	}
}

// Escape a string that it can be inserted/appended with "a/"+string"+"/".
func (e *Edit) Escape(t string) string {
	t = re1.Escape(t)
	t = strings.Replace(t, "\n", `\n`, -1)
	t = strings.Replace(t, "/", `\/`, -1)
	return t
}

func mouseEvent(tb *TextBox, pos image.Point, but int, dir int) (int, [2]int64, bool) {
	dirty := true
	switch {
	case but == -1: // wheel up
		tb.Wheel(pos, 0, 1)
	case but == -2: // wheel down
		tb.Wheel(pos, 0, -1)
	case dir == 0:
		dirty = tb.Move(pos)
	case dir > 0:
		tb.Click(pos, but)
	case dir < 0:
		b, d := tb.Click(pos, -but)
		return b, d, true
	case dir == 3: // step/both; TODO: do we need this?
		tb.Click(pos, but)
		tb.Click(pos, -but)
	}
	return 0, [2]int64{0, 0}, dirty
}

func (e *Edit) Key(r rune, code uint32, dir int, mod uint32) int {
	if r, o := e.MenuKey(r, code, dir, mod); o {
		return r
	}
	dirKeyCode := map[uint32]bool{
		79: true, 80: true, 81: true, 82: true, 74: true, 75: true, 77: true, 78: true, // ← → ↓ ↑, ↑↑↑, ↑↑, ↓↓↓, ↓↓
	}
	if dir == 0 {
		dir = 1
	}
	if dir > 0 && dirKeyCode[code] {
		dirKey(e.TextBox, code)
		return e.DrawSelf()
	}

	// Copy/Paste/Cut: Cntrl-C/V/X
	if code == 6 && e.mods[3] { // c
		if dir > 0 {
			e.TextBox.Copy()
		}
		return e.DrawSelf()
	}
	if code == 25 && e.mods[3] { // v
		if dir > 0 {
			e.TextBox.Paste()
		}
		return e.DrawSelf()
	}
	if code == 27 && e.mods[3] { // x
		if dir > 0 {
			e.TextBox.Cut()
		}
		return e.DrawSelf()
	}

	switch {
	case code == 42: // backspace
		r = '\b'
	case code == 76: // delete
		r = 0x7f
	case r == '\r':
		r = '\n'
	}
	if r > 0 {
		if dir > 0 {
			e.TextBox.Rune(r)
		}
		return e.DrawSelf()
	}

	e.mods = modKey(e.TextBox, e.mods, mod)
	return e.DrawSelf()
}
func dirKey(tb *TextBox, code uint32) {
	switch code {
	case 79: // →
		tb.Dir(1, 0)
	case 80: // ←
		tb.Dir(-1, 0)
	case 82: // ↑
		tb.Dir(0, -1)
	case 81: // ↓
		tb.Dir(0, 1)
	case 74: // ↑↑↑
		tb.Dir(0, math.MinInt16)
	case 75: // ↑↑
		tb.Dir(0, -2)
	case 77: // ↓↓↓
		tb.Dir(0, math.MaxInt16)
	case 78: // ↓↓
		tb.Dir(0, 2)
	default:
		panic("impossible")
	}
}
func modKey(tb *TextBox, mods [4]bool, mod uint32) [4]bool {
	var newMods [4]bool
	if mod&1 != 0 { // shift
		newMods[1] = true
	}
	if mod&4 != 0 { // alt
		newMods[2] = true
	}
	if mod&8 != 0 || mod&2 != 0 { // meta|ctrl
		newMods[3] = true
	}
	for i := 0; i < len(newMods); i++ {
		if newMods[i] != mods[i] {
			m := i
			if !newMods[i] {
				m = -m
			}
			tb.Mod(m)
			mods = newMods
			break
		}
	}
	return mods
}
