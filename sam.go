package ui

import (
	"fmt"
	"os"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
)

// Sam is a vertically split widget that contains a command widget (sam editing commands)
// connected with a text editor widget below.
type Sam struct {
	Split
	cmd *Repl
	edt *Edit
}

// NewSam returns a Sam widget.
func NewSam() Widget {
	sam := &Sam{}
	sam.cmd = &Repl{
		Reply:       false,
		Currentline: false,
	}
	sam.cmd.SetText(rope.New("commands"))
	sam.edt = &Edit{}
	sam.edt.SetText(rope.New("this is text\nsecond line\n"))
	sam.Split = Split{
		Vertical: true,
		Kids:     NewKids(sam.cmd, sam.edt),
	}
	sam.cmd.Interp = sam
	return sam
}

func (sam *Sam) Eval(t string) {
	fmt.Println("sam: ", t)
	_, err := sam.edt.EditWrite(t, os.Stdout) // TODO: usd cmd widget as output.
	if e, ok := err.(edit.NoCommandError); ok {
		sam.edt.SetDot(e.At)
	} else if err != nil {
		fmt.Println(err) // TODO write to cmd widget
	}
}

func (sam *Sam) Cancel() {}
