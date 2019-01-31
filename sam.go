package ui

import (
	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
)

// Sam is a vertically split widget that contains a command widget (sam editing commands)
// connected with a text editor widget below.
type Sam struct {
	Split
	cmd *Repl
	edt *Edit
	w   *Window
}

// NewSam returns a Sam widget.
func NewSam(w *Window) *Sam {
	sam := &Sam{w: w}
	sam.cmd = &Repl{
		Reply: false,
	}
	sam.edt = &Edit{}
	sam.Split = Split{
		Vertical: true,
		Kids:     NewKids(sam.cmd, sam.edt),
	}
	sam.cmd.Interp = sam
	return sam
}

func (sam *Sam) SetTexts(cmd, edt rope.Rope) {
	sam.cmd.SetText(cmd)
	sam.edt.SetText(edt)
}

func (sam *Sam) Eval(t string) {
	_, err := sam.edt.Edit(t)
	if e, ok := err.(edit.NoCommandError); ok {
		sam.edt.SetDot(e.At)
	} else if err != nil {
		sam.cmd.Write([]byte("\n" + err.Error()))
	}
	sam.w.MarkDraw(sam.edt)
}

func (sam *Sam) Cancel() {}
