package editor

import (
	"image"

	"github.com/eaburns/T/edit"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/paint"
)

// Repl is a text editor widget, but sends commands to an interpreter.
type Repl struct {
	Target **Repl
	Edit
	Interp Interpreter
	Reply  bool // Append command to the Repl widget, if it is not on the last line.
	hold   bool
}

type Interpreter interface {
	Eval(string)
	Cancel() // Would be nice to have...
}

/* rm
func (r *Repl) Write(p []byte) (n int, err error) {
	return r.Edit.Write(p)
}
*/

func (r *Repl) Draw(dst *image.RGBA, force bool) {
	if force || r.Edit.Dirty {
		paint.Border(dst, dst.Rect, base.Yellow) // TODO
		r.Edit.Draw(dst.SubImage(dst.Rect.Inset(1)).(*image.RGBA), force)
	}
}

func (r *Repl) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return r.Edit.Mouse(pos, but, dir, mod)
}
func (r *Repl) Key(rn rune, code uint32, dir int, mod uint32) int {
	return r.Edit.Key(rn, code, dir, mod)
}

/*
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


// Mouse button 3 executes the current selection or current line.
// Other mouse events are propagated to the edit widget.
func (r *Repl) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (res Result) {
	if m.Button == 3 {
		res.Consumed = true
		if m.Direction == mouse.DirPress {
			r.exec()
			self.Draw = Dirty
			return res
		}
		return res
	}
	return r.Edit.Mouse(w, self, m, origM, orig)
}


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

	r.exec()
	self.Draw = Dirty
	return res
}
*/

func (r *Repl) exec() {
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
}

// Reply appends s to the end of the text box, if the dot is not on the last line.
func (r *Repl) reply(s string) {
	_, err := r.Edit.Edit("+")
	if _, ok := err.(edit.NoCommandError); ok {
		r.Edit.Write([]byte("\n" + s))
	}
	// If we are on the last line, the error is address not found.
}
