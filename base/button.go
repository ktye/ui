package base

import (
	"image"
	"image/draw"

	"github.com/ktye/ui"
)

type Button struct {
	Base
	Target  **Button
	Execute func() int
	Text    string
	Icon    string
	Color   *image.Uniform
	Primary bool
	off     int
	fill    bool
}

func NewButton(text, icon string, f func() int) *Button {
	return &Button{Text: text, Icon: icon, Execute: f}
}
func FillButton() *Button { return &Button{fill: true} }
func (b *Button) Draw(dst *image.RGBA, force bool) {
	if force == false && b.draw == false {
		return
	}
	b.draw = false
	save := Colors
	defer func() { Colors = save }()
	Colors = ButtonCol
	if b.Primary {
		Colors = ButtonPrim
	}
	clear(dst)
	x := Font.size
	if b.Icon != "" {
		r := dst.Rect.Add(image.Point{Font.size/2 + b.off, Font.size/2 + b.off})
		ico := getIcon(icon{b.Icon, 2 * Font.size})
		draw.DrawMask(dst, r, Colors[0], image.ZP, ico, image.ZP, draw.Over)
		x += 3 * Font.size
	}
	if b.Text != "" {
		String(dst, image.Point{x + b.off, Font.size + b.off}, b.Text)
	}
}
func (b *Button) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if but == 1 && dir > 0 {
		b.off = 1
		return b.DrawSelf()
	} else if but == 1 && dir < 0 {
		b.off = 0
		if b.Execute != nil {
			if b.Execute() < 0 {
				return -1
			}
		}
		return b.DrawSelf()
	}
	return 0
}
func (b *Button) Key(r rune, code uint32, dir int, mod uint32) int {
	if dir > 0 && (code == 44 || code == 40) { // space, enter → execute
		if b.Execute != nil {
			return b.Execute()
		}
	}
	return 0
}
func (b *Button) Size() image.Point {
	H := 3 * Font.size
	if b.fill {
		return image.Point{}
	}
	if b.Text == "" {
		return image.Point{H, H}
	}
	p := StringSize(b.Text)
	p.X += 2 * Font.size
	if b.Icon != "" {
		p.X += H
	}
	p.Y = H
	return p
}

// A button bar places buttons below a kid widget.
type ButtonBar struct {
	Base
	Target   **ButtonBar
	Kid      Kid
	Buttons  []*Button
	Vertical bool
	Flip     bool // top or right
	but      []Kid
	focus    int
}

func NewButtonBar(kid ui.Widget, buttons []*Button) *ButtonBar {
	var b ButtonBar
	b.Kid = Kid{Widget: kid}
	b.Buttons = buttons
	return &b
}
func (b *ButtonBar) SetButtons(buttons []*Button) {
	b.Buttons = buttons
	b.but = nil
}
func (b *ButtonBar) Draw(dst *image.RGBA, force bool) {
	if b.Target != nil {
		*b.Target = b
	}
	if b.but == nil {
		b.but = make([]Kid, len(b.Buttons))
		for i := range b.but {
			b.but[i] = Kid{Widget: b.Buttons[i]}
		}
	}
	if force || dst.Rect != b.Rect {
		b.Rect = dst.Rect
		sum, ishor, isver, isflip, H := 0, 1, 0, 0, 3*Font.size
		if b.Vertical {
			ishor, isver = 0, 1
		}
		if b.Flip {
			isflip = 1
		}
		hsize, fills := 0, 0
		for i := range b.Buttons {
			s := b.Buttons[i].Size()
			b.but[i].Rectangle.Max = s
			sum += s.X*ishor + s.Y*isver
			if s.X > hsize {
				hsize = s.X
			}
			if b.Buttons[i].fill {
				fills++
			}
		}
		space := ishor*dst.Rect.Dx() + isver*dst.Rect.Dy() - sum
		if space < 0 {
			space = 0
		}
		if fills > 0 {
			space /= fills
		}
		p := image.Point{dst.Rect.Min.X, isver*dst.Rect.Min.Y + (1-isflip)*ishor*(dst.Rect.Max.Y-H)}
		p.X += isver * isflip * (dst.Rect.Dx() - hsize)
		for i := range b.but {
			if b.Buttons[i].fill {
				p = p.Add(image.Point{space * ishor, space * isver})
			} else {
				s := b.but[i].Rectangle.Max
				s.X = s.X*ishor + hsize*isver
				r := image.Rectangle{Min: p, Max: p.Add(s)}
				r.Min = r.Min.Add(image.Point{1, 1})
				r.Max = r.Max.Sub(image.Point{1, 1})
				b.but[i].Rectangle = r
				p = p.Add(image.Point{s.X * ishor, s.Y * isver})
			}
		}
		if b.Kid.Widget != nil {
			b.Kid.Rectangle = dst.Rect
			b.Kid.Rectangle.Min.X += hsize * isver * (1 - isflip)
			b.Kid.Rectangle.Min.Y += ishor * isflip * H
			b.Kid.Rectangle.Max.X -= isver * isflip * hsize
			b.Kid.Rectangle.Max.Y -= H * ishor * (1 - isflip)
		}
		bar := b.Rect
		bar.Min.Y = ishor * (1 - isflip) * (bar.Max.Y - H)
		draw.Draw(dst, bar, Colors[1], image.ZP, draw.Src)
	}
	if force || b.draw {
		for i := range b.Buttons {
			if b.Buttons[i].Primary {
				b.Buttons[i].draw = true
				b.Buttons[i].Primary = false
			}
			if b.focus == i {
				b.Buttons[i].draw = true
				b.Buttons[i].Primary = true
			}
		}
	}
	if b.Kid.Widget != nil {
		b.Kid.Draw(dst, force)
	}
	for i := range b.but {
		b.but[i].Draw(dst, force)
	}
}
func (b *ButtonBar) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if but == 1 && dir > 0 {
		for i := range b.but {
			if pos.In(b.but[i].Rectangle) {
				b.focus = i
			}
		}
		if b.Kid.Widget != nil && pos.In(b.Kid.Rectangle) {
			b.focus = -1
		}
	}
	if b.Kid.Widget != nil {
		if pos.In(b.Kid.Rectangle) {
			return b.Kid.Mouse(pos, but, dir, mod)
		}
	}
	for i := range b.but {
		if pos.In(b.but[i].Rectangle) {
			return b.but[i].Mouse(pos, but, dir, mod)
		}
	}
	return 0
}
func (b *ButtonBar) Key(r rune, code uint32, dir int, mod uint32) int {
	if code == 43 { // tab
		if dir > 0 {
			for {
				b.focus++
				if b.focus >= len(b.Buttons) {
					b.focus = 0
				}
				if b.Buttons[b.focus].fill == false {
					break
				}
			}
		}
		return b.DrawSelf()
	}
	if n := b.focus; n > 0 && n < len(b.but) {
		return b.but[n].Key(r, code, dir, mod)
	}
	if b.focus < 0 && b.Kid.Widget != nil {
		return b.Kid.Key(r, code, dir, mod)
	}
	return 0
}
func (b *ButtonBar) MenuSize() image.Point {
	h, count := 0, 0
	for i := range b.Buttons {
		if s := b.Buttons[i].Size().X; s > h {
			h = s
		}
		if b.Buttons[i].fill == false {
			count++
		}
	}
	return image.Point{h, 3 * Font.size * count}
}
