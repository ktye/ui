package ui

import (
	"image/color"
	"strings"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/text"
)

// Sam is a vertically split widget that contains a command widget (sam editing commands)
// connected with a text editor widget below.
type Sam struct {
	Split
	Commands map[string]func(*Sam, string)
	exec     func(*Sam, string)
	Cmd      *Repl
	Edt      *Edit
	w        *Window
}

// NewSam returns a Sam widget.
func NewSam(w *Window) *Sam {
	sam := &Sam{w: w}
	sam.Cmd = &Repl{
		Reply: false,
	}
	sam.Cmd.Edit.styles = sam.styles(argb(0xeaffffff), argb(0x9eeeeeff), argb(0xff0000ff))
	sam.Edt = &Edit{}
	sam.Edt.styles = sam.styles(argb(0xffffeaff), argb(0xeeee9eff), argb(0xff0000ff))
	sam.Split = Split{
		Vertical: true,
		Gutter:   true,
		Ratio:    -0.236,
		Kids:     NewKids(sam.Cmd, sam.Edt),
	}
	sam.Cmd.Interp = sam
	return sam
}

// SetExec sets a callback for button-2 commands to both, the cmd and edt window.
func (sam *Sam) SetExec(e func(*Sam, string)) {
	sam.Cmd.Execute = func(ed *Edit, s string) {
		e(sam, s)
	}
	sam.Edt.Execute = func(ed *Edit, s string) {
		e(sam, s)
	}
}

func (sam *Sam) styles(bg1, bg2, bg3 color.Color) []text.Style {
	return []text.Style{
		text.Style{FG: color.Black, BG: bg1, Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: bg2, Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: bg3, Face: sam.w.font.Face},
		text.Style{FG: color.Black, BG: bg2, Face: sam.w.font.Face},
	}
}

func (sam *Sam) Eval(t string) {
	// Try user defined commands.
	if len(t) > 0 && sam.Commands != nil {
		for c, f := range sam.Commands {
			if strings.HasPrefix(t, c) {
				f(sam, strings.TrimPrefix(t, c))
				return
			}
		}
	}

	// Built-in editing commands.
	_, err := sam.Edt.Edit(t)
	if e, ok := err.(edit.NoCommandError); ok {
		sam.Edt.SetDot(e.At)
	} else if err != nil {
		sam.Cmd.Write([]byte("\n" + err.Error()))
	}
	sam.w.MarkDraw(sam.Edt)
}

func (sam *Sam) Cancel() {}
