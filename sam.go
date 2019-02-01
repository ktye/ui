package ui

import (
	"image/color"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/text"
)

// Sam is a vertically split widget that contains a command widget (sam editing commands)
// connected with a text editor widget below.
type Sam struct {
	Split
	Quit func() (e Event) // executed for q command.
	cmd  *Repl
	edt  *Edit
	w    *Window
}

// NewSam returns a Sam widget.
func NewSam(w *Window) *Sam {
	sam := &Sam{w: w}
	sam.cmd = &Repl{
		Reply: false,
	}
	sam.cmd.Edit.styles = sam.styles(argb(0xeaffffff), argb(0x9eeeeeff))
	sam.edt = &Edit{}
	sam.edt.styles = sam.styles(argb(0xffffeaff), argb(0xeeee9eff))
	sam.Split = Split{
		Vertical: true,
		Gutter:   1,
		Kids:     NewKids(sam.cmd, sam.edt),
	}
	sam.cmd.Interp = sam
	return sam
}

func (sam *Sam) SetTexts(cmd, edt rope.Rope) {
	sam.cmd.SetText(cmd)
	sam.edt.SetText(edt)
}

func (sam *Sam) styles(bg1, bg2 color.Color) []text.Style {
	return []text.Style{
		text.Style{FG: color.Black, BG: bg1, Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: bg2, Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: argb(0x99994CFF), Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: bg2, Face: sam.w.font.Face},
	}
}

func (sam *Sam) Eval(t string) {
	if t == "q" {
		if sam.Quit != nil {
			sam.Quit() // result is ignored
		}
		return
	}
	_, err := sam.edt.Edit(t)
	if e, ok := err.(edit.NoCommandError); ok {
		sam.edt.SetDot(e.At)
	} else if err != nil {
		sam.cmd.Write([]byte("\n" + err.Error()))
	}
	sam.w.MarkDraw(sam.edt)
}

func (sam *Sam) Cancel() {}
