// Sam is a graphical text editor
//
// This is a demonstration of the widget.
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

	exec := func(sam *ui.Sam, s string) {
		fmt.Printf("sam: exec %s\n", s)
	}

	sam := ui.NewSam(w)
	sam.Cmd.SetText(rope.New("command window\n+/is/"))
	sam.Edt.SetText(rope.New("this is line one\nthis is line 2"))
	sam.Commands = map[string]func(*ui.Sam, string){
		"q": func(s *ui.Sam, args string) { os.Exit(0) },
	}
	sam.SetExec(exec)
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
