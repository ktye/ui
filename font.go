package ui

import (
	"image"

	imfont "golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/golang/freetype/truetype"
)

func (w *Window) FontHeight() int {
	return w.font.Face.Metrics().Height.Ceil()
}

func (w *Window) StringSize(s string) image.Point {
	dx := int(imfont.MeasureString(w.font.Face, s)+32) / 64
	dy := w.FontHeight()
	return image.Point{dx, dy}
}

func (w *Window) defaultFont() {
	w.font.ttf = nil
	w.font.size = 13
	w.font.Face = basicfont.Face7x13
}

func (w *Window) SetFont(ttf []byte, size int) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	w.font.ttf = ttf
	opt := truetype.Options{
		Size: float64(size),
		DPI:  float64(w.Display.PixelsPerPt) * 72.0,
	}
	if opt.Size == 0 {
		opt.Size = 18
	}
	w.font.size = size
	w.font.Face = truetype.NewFace(f, &opt)
}

type font struct {
	imfont.Face
	ttf  []byte
	size int
}

// Scale is the callback for shift-wheel up/down.
// It changes the font size.
func (w *Window) scale(up bool) {
	if up {
		w.font.size++
	} else if w.font.size < 5 {
		return
	} else {
		w.font.size--
	}
	w.SetFont(w.font.ttf, w.font.size)

	w.Top.Layout = Dirty
	w.Render()
}
