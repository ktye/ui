// Sam is a graphical text editor
//
// This is an example of the edit widget with an additional sam command widget on top.
//
// TODO: The idea is that any edit widget can be converted to a split sam+edit widget
// with a context-menu command. Context menus still need to be done.
package main

import (
	"fmt"

	"github.com/ktye/ui"
)

func main() {
	w := ui.New(nil)
	w.SetKeyTranslator(ui.AplKeyboard{})

	w.Top.W = ui.NewSam()
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
