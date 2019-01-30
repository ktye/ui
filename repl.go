package ui

import (
	"image"
	"image/draw"

	"github.com/eaburns/T/edit"
	"github.com/ktye/ui/paint"

	"golang.org/x/mobile/event/key"
)

// Repl is like a text editor widget, but sends commands to an interpreter.
type Repl struct {
	Target **Repl
	Edit
	Interp Interpreter
	Reply  bool // Append command to the Repl widget, if it is not on the last line.
	hold   bool
}

type Interpreter interface {
	Eval(string)
	Cancel()
}

func (r *Repl) Write(p []byte) (n int, err error) {
	return r.Edit.Write(p)
}

func (r *Repl) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	if r.Target != nil {
		*r.Target = r
	}
	r.Edit.Layout(w, self, sizeAvail.Sub(pt(2)), force)
	self.R = rect(sizeAvail)
}

func (r *Repl) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	hover := m.In(self.R)
	cs := &w.Primary
	colors := cs.Normal
	if hover {
		colors = cs.Hover
	}
	c := colors.Text
	if r.hold {
		c = colors.Border
	}
	paint.Border(img, self.R, c)

	sr := self.R
	self.R.Max = sr.Max.Sub(pt(2))
	r.Edit.Draw(w, self, img, orig.Add(pt(1)), m, force)
	self.R = sr
}

// TODO: Mouse Execute selection on button 3

// Keys are forwarded to the underlying edit widget, except for special handling of ESC and ENTER.
// On ENTER the current selected text is interpreted.
// If nothing is selected, the current line is used.
// Pressing ESC toggles hold mode.
// On hold, the interception of ENTER key events is prevented which allows to type multiline commands.
func (r *Repl) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	// Toggle hold mode on escape.
	if k.Code == key.CodeEscape {
		res.Consumed = true
		if k.Direction != key.DirRelease {
			r.hold = !r.hold
			self.Draw = Dirty
		}
		return res
	}

	if r.hold == true || k.Code != key.CodeReturnEnter {
		return r.Edit.Key(w, self, k, m, orig)
	}

	res.Consumed = true
	if k.Direction != key.DirRelease {
		return res
	}

	if r.Interp == nil {
		return res
	}

	t := r.Edit.TextBox.Selection()
	if t == "" {
		// Nothing is selected: execute current line.
		r.Edit.MarkAddr("-+")
		t = r.Edit.TextBox.Selection()
	}
	if len(t) > 0 && t[len(t)-1] == '\n' {
		t = t[:len(t)-1]
	}
	if t != "" {
		if r.Reply {
			r.reply(t)
		}
		r.Interp.Eval(t)
	}
	self.Draw = Dirty
	return res
}

// Reply appends s to the end of the text box, if the dot is not on the last line.
func (r *Repl) reply(s string) {
	_, err := r.Edit.Edit(`+/\n/`)
	if e, ok := err.(edit.NoCommandError); ok {
		dot := r.Edit.Dot()
		if e.At[0] > dot[1] {
			r.Edit.Write([]byte("\n" + s))
		}
	}
	// If there is no newline, err is "no match" and we don't have to reply the line.
}
