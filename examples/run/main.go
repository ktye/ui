// Run is a test program for ui.
package main

import (
	"fmt"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/ktye/ui"
)

func main() {

	w := ui.New(nil)
	ui.RegisterIcon("HardwareComputer", icons.HardwareComputer)
	ui.RegisterIcon("ActionTimeline", icons.ActionTimeline)
	ui.RegisterIcon("FileAttachment", icons.FileAttachment)
	ui.RegisterIcon("HardwareKeyboard", icons.HardwareKeyboard)

	aplButton := &ui.Button{Icon: "HardwareComputer", Colorset: &w.Primary}
	imgButton := &ui.Button{Icon: "ActionTimeline", Colorset: &w.Primary}
	edtButton := &ui.Button{Icon: "FileAttachment", Colorset: &w.Primary}
	keyButton := &ui.Button{Icon: "HardwareKeyboard", Colorset: &w.Primary}
	buttons := ui.NewBox(aplButton, imgButton, edtButton, keyButton)

	repl := &ui.Repl{}
	// edit := &ui.Edit{}

	w.SetKeyTranslator(ui.AplKeyboard{})
	w.Top.W = &ui.Split{
		Kids: ui.NewKids(
			buttons,
			repl,
		),
	}

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

/*

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
*/
