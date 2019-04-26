package main

import (
	"fmt"
	"os"

	"github.com/eaburns/T/rope"
	"github.com/ktye/ui"
)

func main() {
	l := &ui.List{}
	r := rope.New("")
	for i := 0; i < 3; i++ {
		r = rope.Append(r, rope.New(fmt.Sprintf("line %d alpha beta gamma\n", i)))
	}
	l.SetText(r)
	l.Execute = func() bool { println("exec"); return false }
	l.Delete = func() { println("delete") }
	l.Quit = func() { os.Exit(0) }
	l.Single = true
	w := ui.New(nil)
	w.Top.W = l
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
