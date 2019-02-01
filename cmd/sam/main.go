// Sam is a graphical text editor
//
// This is an example of the edit widget with an additional sam command widget on top.
//
// TODO: The idea is that any edit widget can be converted to a split sam+edit widget
// with a context-menu command. Context menus still need to be done.
package main

import (
	"fmt"
	"os"

	"github.com/eaburns/T/rope"
	"github.com/ktye/ui"
)

func main() {
	w := ui.New(nil)
	w.SetKeyTranslator(ui.AplKeyboard{})

	sam := ui.NewSam(w)
	sam.Cmd.SetText(rope.New("command window\n+/is/"))
	sam.Edt.SetText(rope.New("this is line one\nthis is line 2"))
	sam.Commands = map[string]func(*ui.Sam, string){
		"q": func(s *ui.Sam, args string) { os.Exit(0) },
	}
	w.Top.W = sam
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
