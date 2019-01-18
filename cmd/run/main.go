// Run is a test program for ui.
package main

import (
	"fmt"

	"github.com/ktye/ui"
)

func main() {
	w := ui.New(nil)

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
