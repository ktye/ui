package ui

import (
	"fmt"
	"image"

	"golang.org/x/mobile/event/key"
)

// Repl is like a text editor widget, but sends commands to an interpreter.
type Repl struct {
	Target **Repl
	Edit
	Interp      Interpreter
	Reply       bool // Append command to the Repl widget, if it is not on the last line.
	Currentline bool // Use current line if nothing is selected.
}

type Interpreter interface {
	Eval(string)
	Cancel()
}

// TODO: On ENTER execute current line if nothing is selected.
// Otherwise execute selected (possibly multiline) text.
// TODO: cancel with ESC?
func (r *Repl) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	res.Consumed = true
	if k.Code == key.CodeReturnEnter {
		if k.Direction != key.DirRelease {
			if r.Interp == nil {
				fmt.Println("Repl: Interpeter is nil")
				return res
			}
			t := r.Edit.TextBox.Selection()
			if t != "" {
				r.Interp.Eval(t)
				return res
			}
			if r.Currentline == false {
				goto def
			}
			fmt.Println("TODO: repl: exec current line")
		}
		return res
	}
def:
	e := r.Edit
	e.mods = e.keyEvent(e.mods, k)
	res.Consumed = true
	self.Draw = Dirty
	return res
}
