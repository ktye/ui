package base

import (
	"image"
	"image/color"
	"image/draw"
)

type Menu struct {
	*ButtonBar
	active bool
	point  image.Point
	rect   image.Rectangle
}

func NewMenu(items []*Button) *Menu {
	m := Menu{
		ButtonBar: NewButtonBar(nil, items),
	}
	m.ButtonBar.Vertical = true
	return &m
}

func (m *Menu) drawMenu(dst *image.RGBA) {
	if m.active == false {
		return
	}
	//m.active = false
	s := m.ButtonBar.MenuSize()
	dr := dst.Bounds()
	r := image.Rectangle{Min: m.point, Max: m.point.Add(s)}.Sub(image.Point{Font.size, Font.size})
	if d := r.Max.Y - dr.Max.Y; d > 0 {
		r = r.Sub(image.Point{0, d})
	}
	if d := r.Min.Y - dr.Min.Y; d < 0 {
		r = r.Sub(image.Point{0, d})
	}
	if d := r.Max.X - dr.Max.X; d > 0 {
		r = r.Sub(image.Point{d, 0})
	}
	if d := r.Min.X - dr.Min.X; d < 0 {
		r = r.Sub(image.Point{d, 0})
	}
	u := &image.Uniform{color.RGBA{128, 128, 128, 128}}
	draw.Draw(dst, dst.Bounds(), u, image.ZP, draw.Over)
	m.ButtonBar.Draw(dst.SubImage(r).(*image.RGBA), true)
}
func (m *Menu) mouse(pos image.Point, but int, dir int, mod uint32) int {
	if !pos.In(m.ButtonBar.Base.Rect) && but == 1 && dir < 0 { // but-1-click outside → close menu
		m.active = false
		return 1
	}
	r := m.ButtonBar.Mouse(pos, but, dir, mod)
	if r < 0 {
		m.active = false
	}
	return r
}
func (m *Menu) key(r rune, code uint32, dir int, mod uint32) int {
	if code == 41 { // Escape
		if dir < 0 { // release Escape → close menu
			m.active = false
			return 1
		}
		return 0
	}
	ret := m.ButtonBar.Key(r, code, dir, mod)
	if ret < 0 {
		m.active = false
	}
	return ret
}
