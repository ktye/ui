package ui

import (
	"fmt"
	"image"
	"io"

	"golang.org/x/mobile/event/key"
)

// Repl is like a text editor widget, but sends commands to an interpreter.
type Repl struct {
	Target **Repl
	Edit
	Interp Interpreter
}

type Interpreter interface {
	Eval(string, io.Writer)
	Cancel()
}

// TODO: On ENTER execute current line if nothing is selected.
// Otherwise execute selected (multiline) text.
// TODO: cancel with ESC?
func (r *Repl) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (res Result) {
	res.Consumed = true
	if k.Code == key.CodeReturnEnter {
		if k.Direction != key.DirRelease {
			fmt.Println("Eval...") // TODO
			return res
		}
		return res
	}
	e := r.Edit
	e.mods = e.keyEvent(e.mods, k)
	res.Consumed = true
	self.Draw = Dirty
	return res
}
