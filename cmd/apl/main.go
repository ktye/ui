package main

import (
	"fmt"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/ui"
)

func main() {
	var interp interp
	repl := &ui.Repl{Reply: true}
	repl.Edit.Nowrap = true
	repl.SetText(rope.New("\t"))
	interp.repl = repl
	repl.Interp = &interp

	a := apl.New(repl)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	interp.apl = a

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
	apl  *apl.Apl
	repl *ui.Repl
}

func (i *interp) Eval(s string) {
	i.repl.Write([]byte{'\n'})
	if err := i.apl.ParseAndEval(s); err != nil {
		i.repl.Write([]byte(err.Error()))
	}
	i.repl.Write([]byte("\t"))
	i.repl.Edit.MarkAddr("$")
}

func (i *interp) Cancel() {}
