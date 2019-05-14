package editor

import (
	"image"

	"github.com/eaburns/T/edit"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/paint"
)

// Repl is a text editor widget, but sends commands to an interpreter.
// Escape toggles hold mode, commands can be send with the middle mouse or the enter key.
type Repl struct {
	Target **Repl
	Edit
	Interp Interpreter
	Prompt string
	Reply  bool // Append command to the Repl widget, if it is not on the last line.
	hold   bool
}

type Interpreter interface {
	Eval(string) string
	// Cancel() would be nice to have...
}

func (r *Repl) Draw(dst *image.RGBA, force bool) {
	if r.Edit.Execute == nil {
		r.Edit.Execute = r.DefaultExec
	}
	if force || r.Edit.Dirty {
		if r.hold {
			paint.Border(dst, dst.Rect, base.Blue)
			dst = dst.SubImage(dst.Rect.Inset(1)).(*image.RGBA)
		}
		r.Edit.Draw(dst, force)
	} else {
		r.DrawMenu(dst)
	}
}

func (r *Repl) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return r.Edit.Mouse(pos, but, dir, mod)
}
func (r *Repl) Key(rn rune, code uint32, dir int, mod uint32) int {
	if code == 41 { // escape (toggle hold mode)
		if dir < 0 {
			r.hold = !r.hold
			return r.DrawSelf()
		}
		return 0
	}
	if r.hold == true || code != 40 { // return
		return r.Edit.Key(rn, code, dir, mod)
	}
	if dir >= 0 {
		return 0
	}

	s := r.Edit.TextBox.Selection()
	if s == "" { // Nothing is selected: execute current line.
		r.Edit.MarkAddr("-+")
		s = r.Edit.TextBox.Selection()
	}
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}
	return r.Execute(nil, s)
}

func (r *Repl) DefaultExec(_ *Edit, s string) int {
	if s == "" {
		return 0
	}
	if r.Reply {
		r.reply(s)
	}
	s = r.Interp.Eval(s)
	b := append([]byte{'\n'}, []byte(s)...) // TODO don't add \n if we are already at the start of a new line
	if b[len(b)-1] != '\n' {
		b = append(b, '\n')
	}
	b = append(b, []byte(r.Prompt)...)
	r.Edit.Write(b)
	r.Edit.MarkAddr("$")
	return 1
}
func (r *Repl) reply(s string) { // append s to the end of the text box, if dot is not on the last line.
	_, err := r.Edit.Edit("+")
	if _, ok := err.(edit.NoCommandError); ok {
		r.Edit.Write([]byte("\n" + s))
	}
	// If we are on the last line, the error is address not found.
}
