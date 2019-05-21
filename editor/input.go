package editor

import (
	"image"
	"image/draw"

	"github.com/eaburns/T/rope"
	"github.com/eaburns/T/text"
	"github.com/ktye/ui/base"
)

// Input is a widget for one or more single line input editors with an optional label.
// If Options are present, they are shown in a list below the input field.
type Input struct {
	Texts     []string
	Labels    []string
	Commit    func([]string) error
	Cancel    func() int
	Margin    int
	Multi     bool
	Options   []string
	NewLabel  func(int) string
	edt       []*Edit
	cur       int // editor in focus
	lst       *base.List
	err       error
	ect       []image.Rectangle // edit rectangles
	lct, last image.Rectangle
	lstfocus  bool
}

func NewInput(label, text string) *Input {
	return &Input{Labels: []string{label}, Texts: []string{text}, Margin: base.Font.Size()}
}
func NewInputs(labels, texts []string) *Input {
	return &Input{Labels: labels, Texts: texts, Margin: base.Font.Size(), Multi: true}
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
func NewInputsWithButtons(labels, texts []string) *base.ButtonBar {
	input := NewInputs(labels, texts)
	buttons := []*base.Button{
		base.FillButton(),
		base.NewButton("ok", "ok", input.commit),
		base.NewButton("add", "add", input.add),
		base.NewButton("del", "del", input.del),
		base.NewButton("cancel", "cancel", input.cancel),
		base.FillButton(),
	}
	return base.NewButtonBar(input, buttons)
}
func (e *Input) layout(dst *image.RGBA) {
	fs := base.Font.Size()

	// Draw labels.
	max := 0
	at := image.Point{e.Margin, e.Margin}
	for i := range e.Labels {
		w := base.String(dst, at, e.Labels[i])
		if w > max {
			max = w
		}
		at.Y += 2 * fs
	}

	// Draw text edits.
	if len(e.ect) < len(e.Texts) {
		e.ect = make([]image.Rectangle, len(e.Texts))
	}
	at = image.Point{e.Margin + max + fs, e.Margin}
	w := dst.Rect.Dx() - 2*e.Margin - max - fs
	h := (3 * fs) / 2
	for i := range e.Texts {
		min := dst.Rect.Min.Add(at)
		e.ect[i] = image.Rectangle{Min: min, Max: min.Add(image.Point{w, h})}
		at.Y += 2 * fs
	}
	errpos := image.Point{0, at.Y - dst.Rect.Min.Y}
	at.Y += 2 * fs

	if len(e.Options) > 0 {
		base.String(dst, image.Point{0, at.Y}, "Options:")
		at.Y += 2 * fs
	}

	// Option list rectangle.
	e.lct = dst.Rect
	e.lct.Min.Y += at.Y

	if e.err != nil {
		save := base.Colors
		base.Colors[0] = base.Red.Uniform()
		base.String(dst, errpos, e.err.Error())
		base.Colors = save
	}
}
func (e *Input) Draw(dst *image.RGBA, force bool) {
	if force || dst.Rect != e.last {
		e.last = dst.Rect
		draw.Draw(dst, dst.Rect, base.Colors[1], image.ZP, draw.Src) // clear
		e.layout(dst)
	}
	if len(e.edt) < len(e.Texts) {
		e.edt = make([]*Edit, len(e.Texts))
		for i, t := range e.Texts {
			ed := New(t)
			ed.styles = []text.Style{
				text.Style{FG: base.Black, BG: base.LightYellow, Face: base.Font.Face},
				text.Style{FG: base.White, BG: base.Blue, Face: base.Font.Face},
				text.Style{FG: base.Black, BG: base.Green, Face: base.Font.Face},
				text.Style{FG: base.Black, BG: base.Red, Face: base.Font.Face},
			}
			ed.initText(e.ect[i].Size())
			e.edt[i] = ed
		}
	}
	if len(e.Options) > 0 && e.lst == nil {
		e.lst = base.NewList(e.Options)
		e.lst.Single = true
		e.lst.Execute = func() int {
			if idx := e.lst.Selection(); len(idx) == 0 {
				return 0
			} else if e.cur >= 0 && e.cur < len(e.edt) {
				ed := e.edt[e.cur]
				ed.TextBox.SetText(rope.New(e.lst.List[idx[0]].String()))
				ed.Dirty = true
				if len(e.Texts) == 1 {
					return e.commit()
				}
				return 1
			}
			return 0
		}
		e.setCurrent(0)
	}
	for i, ed := range e.edt {
		ed.Draw(dst.SubImage(e.ect[i]).(*image.RGBA), force)
	}
	if e.lst != nil {
		e.lst.Draw(dst.SubImage(e.lct).(*image.RGBA), force)
	}
}
func (e *Input) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if e.lst != nil && pos.In(e.lct) {
		if but == 1 {
			e.lstfocus = true
		}
		return e.lst.Mouse(pos, but, dir, mod)
	}
	ret := 0
	if but == 1 {
		for i := range e.ect {
			if pos.In(e.ect[i]) {
				e.setCurrent(i)
				e.lstfocus = false
				ret = 1
				break
			}
		}
	}
	if len(e.edt) == 0 {
		return ret
	}
	r := e.edt[e.cur].Mouse(pos, but, dir, mod)
	if r == 0 {
		return ret
	}
	return r
}
func (e *Input) Key(r rune, code uint32, dir int, mod uint32) int {
	if e.lstfocus {
		return e.lst.Key(r, code, dir, mod)
	}
	switch code {
	case 81: // down
		if dir < 0 && e.cur < len(e.edt)-1 {
			e.setCurrent(e.cur + 1)
		}
		return 1
	case 82: //up
		if dir < 0 && e.cur > 0 {
			e.setCurrent(e.cur - 1)
		}
		return 1
	case 75, 78: // page-up,down
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
	return e.edt[e.cur].Key(r, code, dir, mod)
}
func (e *Input) commit() int {
	if e.Commit == nil {
		return 0
	}
	e.last = image.Rectangle{}
	save := e.Texts
	e.Texts = make([]string, len(e.edt))
	for i, ed := range e.edt {
		e.Texts[i] = ed.TextBox.text.String()
	}
	e.err = e.Commit(e.Texts)
	if e.err != nil {
		e.Texts = save
		e.edt = nil
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
func (e *Input) setCurrent(n int) {
	e.cur = n
	if e.cur >= len(e.Texts) {
		e.cur = len(e.Texts) - 1
	}
	if e.cur < 0 {
		e.cur = 0
	}
	if e.cur >= len(e.edt) {
		return
	}
	for i := range e.edt {
		ed := e.edt[i]
		ed.Dirty = true
		ln := int64(len(ed.TextBox.text.String()))
		ed.TextBox.dots[1].At = [2]int64{ln, ln}
		e.edt[i].Dirty = true
	}
	ed := e.edt[e.cur]
	s := ed.TextBox.text.String()
	ed.TextBox.dots[1].At = [2]int64{0, int64(len(s))}
	for i, t := range e.Options {
		if t == s {
			e.lst.SelectSingle(i)
			break
		}
	}
}
func (e *Input) add() int {
	e.edt = nil
	e.last = image.Rectangle{}
	s := ""
	if len(e.Texts) > 0 {
		s = e.Texts[len(e.Texts)-1]
	}
	e.Texts = append(e.Texts, s)
	if e.NewLabel != nil {
		e.Labels = append(e.Labels, e.NewLabel(len(e.Labels)))
	}
	e.setCurrent(len(e.Texts) - 1)
	return 1
}
func (e *Input) del() int {
	if len(e.Texts) == 0 {
		return 0
	}
	e.Texts = e.Texts[:len(e.Texts)-1]
	if len(e.Labels) > len(e.Texts) {
		e.Labels = e.Labels[:len(e.Texts)]
	}
	e.edt = nil
	e.last = image.Rectangle{}
	if e.cur >= len(e.Texts) {
		e.setCurrent(len(e.Texts) - 1)
	}
	return 1
}
