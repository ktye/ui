package plot

import (
	"errors"
	"image"
	"image/color"
	"strings"

	"github.com/ktye/plot"
	"github.com/ktye/ui/base"
)

func (w *Widget) setCaption(p plot.Plots) {
	if len(p) > 0 {
		p[0].SetCaptionColors()
	}
	c, err := p.MergedCaption()
	if err != nil || len(c.Columns) == 0 {
		return
	}
	w.c = &c

	s, off, err := w.CaptionStrings()
	if err != nil {
		return
	}
	white := &image.Uniform{color.White}
	var colors []base.Colorset
	colorm := make(map[color.RGBA]int)
	entries := make([]base.Stringer, len(s))
	for i := range entries {
		line := CaptionLine{s: s[i], c: 0}
		c := w.c.Color(i, off)
		if c != nil {
			rgb := color.RGBAModel.Convert(c).(color.RGBA)
			n, ok := colorm[rgb]
			if ok == false {
				n = len(colorm)
				colorm[rgb] = n
				colors = append(colors, base.Colorset{&image.Uniform{rgb}, white})
			}
			line.c = n + 1
		}
		entries[i] = line
	}
	w.List.Colorsets = colors
	w.List.List = entries
}
func (w *Widget) CaptionStrings() ([]string, int, error) {
	if w.c == nil {
		return nil, 0, errors.New("caption is empty")
	}
	var b strings.Builder
	if lineOffset, err := w.c.WriteTable(&b, plot.Numbers); err != nil {
		return nil, 0, err
	} else {
		s := strings.Split(b.String(), "\n")
		if len(s) > 0 && len(s[len(s)-1]) == 0 {
			s = s[:len(s)-1]
		}
		return s, lineOffset, nil
	}
}

// CaptionLine implements a base.Stringer with a color.
type CaptionLine struct {
	s string
	c int
}

func (l CaptionLine) String() string { return l.s }
func (l CaptionLine) Color() int     { return l.c }
