package editor

import (
	"image"
	"image/draw"

	"github.com/eaburns/T/text"
	"github.com/ktye/ui/base"
)

// Input is a single line input with an optional label.
type Input struct {
	*Edit
	Text   string
	Label  string
	Commit func() error
	Cancel func() int
	Margin int
	err    error
	rect   image.Rectangle
	est    image.Rectangle
}

func NewInputWithButtons(label, text string) *base.ButtonBar {
	input := NewInput(label, text)
	buttons := []*base.Button{
		base.FillButton(),
		base.NewButton("ok", "ok", input.commit),
		base.NewButton("cancel", "cancel", input.cancel),
		base.FillButton(),
	}
	return base.NewButtonBar(input, buttons)
}
func NewInput(label, text string) *Input {
	return &Input{Label: label, Text: text, Margin: base.Font.Size()}
}
func (e *Input) Draw(dst *image.RGBA, force bool) {
	if force || dst.Rect != e.rect {
		e.rect = dst.Rect
		draw.Draw(dst, dst.Rect, base.Colors[1], image.ZP, draw.Src) // clear
		e.est = e.rect.Inset(e.Margin)                               // margin
		if e.Label != "" {                                           // label
			x := base.String(dst.SubImage(e.est).(*image.RGBA), image.Point{}, e.Label)
			x += base.Font.Size()
			e.est.Min.X += x
		}
		line := e.est
		line.Max.Y = line.Min.Y + (3*base.Font.Size())/2
		e.est = e.est.Intersect(line) // single line editor rectangle
		if e.err != nil {
			save := base.Colors
			base.Colors[0] = base.Red.Uniform()
			base.String(dst, image.Point{e.Margin, e.Margin + 3*base.Font.Size()}, e.err.Error())
			base.Colors = save
		}
	}
	if e.Edit == nil {
		e.Edit = New(e.Text)
		e.Edit.styles = []text.Style{
			text.Style{FG: base.Black, BG: base.LightYellow, Face: base.Font.Face},
			text.Style{FG: base.White, BG: base.Blue, Face: base.Font.Face},
			text.Style{FG: base.Black, BG: base.Green, Face: base.Font.Face},
			text.Style{FG: base.Black, BG: base.Red, Face: base.Font.Face},
		}
		e.Edit.initText(e.est.Size())
		e.Edit.TextBox.dots[1].At = [2]int64{0, int64(len(e.Text))}
	}
	e.Edit.Draw(dst.SubImage(e.est).(*image.RGBA), force)
}
func (e *Input) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return e.Edit.Mouse(pos, but, dir, mod)
}
func (e *Input) Key(r rune, code uint32, dir int, mod uint32) int {
	switch code {
	case 75, 78, 81, 82: // [page]up,down
		return 0
	case 40: // enter
		if dir > 0 {
			return e.commit()
		}
		return 0
	case 41: // esc
		if dir > 0 {
			return e.cancel()
		}
		return 0
	}
	return e.Edit.Key(r, code, dir, mod)
}
func (e *Input) commit() int {
	e.rect = image.Rectangle{}
	save := e.Text
	if e.Edit != nil && e.Edit.TextBox != nil {
		e.Text = e.Edit.TextBox.text.String()
	}
	if e.Commit == nil {
		return 0
	}
	e.err = e.Commit()
	if e.err != nil {
		e.Text = save
		e.Edit = nil
		return 1
	}
	return -1
}
func (e *Input) cancel() int {
	if e.Cancel == nil {
		return 0
	}
	return e.Cancel()
}
