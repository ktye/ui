// Run is a test program for ui.
package main

import (
	"bytes"
	"fmt"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
	"github.com/ktye/ui"
)

func main() {
	var editor *ui.Edit
	var repl *ui.Repl

	ui.RegisterIcon("DeviceGPSFixed", icons.DeviceGPSFixed)
	ui.RegisterIcon("DeviceSignalWiFi1BarLock", icons.DeviceSignalWiFi1BarLock)

	w := ui.New(nil)

	w.SetKeyTranslator(ui.AplKeyboard{})
	w.Top.W = ui.NewBox(
		&ui.Button{
			Text:     "APL",
			Colorset: &w.Primary,
			Click: func() (e ui.Event) {
				runapl(editor)
				return
			},
		},
		&ui.Button{
			Text: "Edit",
			Click: func() (e ui.Event) {
				edit(editor)
				return
			},
		},
		&ui.Button{
			Icon:     "DeviceGPSFixed",
			Colorset: &w.Primary,
		},
		&ui.Button{
			Text:     "Icon+Text",
			Icon:     "DeviceSignalWiFi1BarLock",
			Colorset: &w.Primary,
		},
		&ui.Edit{
			Target: &editor,
		},
		&ui.Repl{
			Target: &repl,
		},
	)
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

func edit(e *ui.Edit) {
	a, b := e.Edit(e.Text().String())
	fmt.Println(a, b)
}

func runapl(e *ui.Edit) {
	var buf bytes.Buffer
	a := apl.New(&buf)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	err := a.ParseAndEval(e.Text().String())
	if err != nil {
		e.SetText(rope.New(err.Error()))
		return
	}
	r, err := rope.ReadFrom(&buf)
	if err != nil {
		fmt.Println(err)
	} else {
		e.SetText(r)
	}
	// TODO: w.Call needs to be implemented to update ui asynchronously
}
