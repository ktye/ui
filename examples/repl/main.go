package main

import (
	"fmt"

	"github.com/eaburns/T/rope"
	"github.com/ktye/ui"
)

func main() {
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.Edit.Nowrap = true
	repl.SetText(rope.New("\t"))
	interp.repl = repl
	repl.Interp = &interp

	w := ui.New(nil)
	w.SetKeyTranslator(ui.AplKeyboard{})
	w.Top.W = repl
	w.Render()

	for {
		select {
		case e := <-w.Inputs:
			w.Input(e)

		case err, ok := <-w.Error:
			if !ok {
				return
			}
			fmt.Println("ui:", err)
		}
	}
}

type interp struct {
	repl *ui.Repl
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	i.repl.Write([]byte("> " + s + "\n"))
	i.repl.Write([]byte("\t"))
	i.repl.Edit.MarkAddr("$")
}

func (i *interp) Cancel() {}
