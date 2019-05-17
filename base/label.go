package base

import (
	"image"
)

// NewLabel returns a label with default colors, that is drawn to the top left corner.
func NewLabel(s string) *Label {
	return &Label{Text: s, Align: 6}
}

// NewLink returns a hyperlink label.
func NewLink(s string, target func() int) *Label {
	return &Label{Text: s, Align: 6, Target: target, Colors: &LinkColors}
}

// Label is a single line non-editable text string,
// that aligns the text in the destination area:
// 	6---5---4
// 	7tex8tex3
// 	0---1---2
type Label struct {
	Text       string
	Align      int       // 0..8
	Colors     *Colorset // custom color (nil is default Color)
	Margin     int       // border inset pixels
	Target     func() int
	rect, used image.Rectangle
}

func (l *Label) Draw(dst *image.RGBA, force bool) {
	if dst.Rect != l.rect || force {
		l.rect = dst.Rect
		if l.Colors != nil {
			save := Colors
			Colors = *l.Colors
			defer func() { Colors = save }()
		}
		clear(dst)
		if l.Margin > 0 {
			dst = dst.SubImage(dst.Rect.Inset(l.Margin)).(*image.RGBA)
		}

		// Measure width by drawing with white on dst.
		w := 0
		if a := l.Align; a != 0 && a != 7 && a != 6 {
			Colors[0], Colors[1] = Colors[1], Colors[0]
			w = String(dst, image.Point{}, l.Text)
			Colors[0], Colors[1] = Colors[1], Colors[0]
		}

		// Alignment.
		at := image.Point{}
		switch l.Align {
		case 2, 3, 4: // right
			at.X += dst.Rect.Dx() - w
		case 1, 5, 8: // x-center
			at.X += (dst.Rect.Dx() - w) / 2
		}
		switch l.Align {
		case 0, 1, 2: // y-center
			at.Y += dst.Rect.Dy() - Font.size
		case 7, 8, 3: // bottom
			at.Y += (dst.Rect.Dy() - Font.size) / 2
		}

		x := String(dst, at, l.Text)
		l.used.Min = dst.Rect.Min.Add(at)
		l.used.Max = l.used.Min.Add(image.Point{x, Font.size})
	}
}
func (l *Label) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if l.Target != nil && but == 1 && dir > 0 && pos.In(l.used) {
		return l.Target()
	}
	return 0
}
func (l *Label) Key(r rune, code uint32, dir int, mod uint32) int { return 0 }
func (l *Label) Used() image.Rectangle                            { return l.used }
