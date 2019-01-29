// Paint supplies some image drawing primitives.
package paint

import (
	"image"
	"image/color"
	"image/draw"
)

func Border(m draw.Image, r image.Rectangle, c color.Color) {
	x0 := r.Min.X
	x1 := r.Max.X - 1
	y0 := r.Min.Y
	y1 := r.Max.Y - 1
	Line(m, image.Pt(x0, y0), image.Pt(x0, y1), c)
	Line(m, image.Pt(x0, y1), image.Pt(x1, y1), c)
	Line(m, image.Pt(x1, y1), image.Pt(x1, y0), c)
	Line(m, image.Pt(x1, y0), image.Pt(x0, y0), c)
}

func RoundedBorder(m draw.Image, r image.Rectangle, c color.Color) {
	radius := 3
	x0 := r.Min.X
	x1 := r.Max.X - 1
	y0 := r.Min.Y
	y1 := r.Max.Y - 1
	tl := image.Pt(x0+radius, y0+radius)
	bl := image.Pt(x0+radius, y1-radius)
	br := image.Pt(x1-radius, y1-radius)
	tr := image.Pt(x1-radius, y0+radius)
	arc(m, tl, radius, radius, c, 90, 90)
	arc(m, bl, radius, radius, c, 180, 90)
	arc(m, br, radius, radius, c, 270, 90)
	arc(m, tr, radius, radius, c, 0, 90)
	Line(m, image.Pt(x0, y0+radius), image.Pt(x0, y1-radius), c)
	Line(m, image.Pt(x0+radius, y1), image.Pt(x1-radius, y1), c)
	Line(m, image.Pt(x1, y1-radius), image.Pt(x1, y0+radius), c)
	Line(m, image.Pt(x1-radius, y0), image.Pt(x0+radius, y0), c)
}

func arc(m draw.Image, c image.Point, a, b int, color color.Color, alpha, phi int) {
	// We only handle phi=90 and alpha: 0, 90, 180, 270.
	var p0, p1 image.Point
	switch alpha {
	case 0:
		p0 = image.Point{a, 0}
		p1 = image.Point{0, -b}
	case 90:
		p0 = image.Point{0, -b}
		p1 = image.Point{-a, 0}
	case 180:
		p0 = image.Point{-a, 0}
		p1 = image.Point{0, b}
	case 270:
		p0 = image.Point{0, b}
		p1 = image.Point{a, 0}
	}

	// Just draw a line.
	Line(m, p0.Add(c), p1.Add(c), color)
}

func Line(m draw.Image, a, b image.Point, c color.Color) {
	line(m, a.X, a.Y, b.X, b.Y, c)
}

// Line draws a line with Besenham's algorithm.
// It only uses integer pixels.
func line(m draw.Image, x0, y0, x1, y1 int, c color.Color) {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	var dx, dy, sx, sy, e, e2 int

	dx = abs(x1 - x0)
	dy = -abs(y1 - y0)
	if sx = -1; x0 < x1 {
		sx = 1
	}
	if sy = -1; y0 < y1 {
		sy = 1
	}
	e = dx + dy
	for {
		m.Set(x0, y0, c)
		if x0 == x1 && y0 == y1 {
			break
		}
		if e2 = 2 * e; e2 >= dy {
			e += dy
			x0 += sx
		} else if e2 <= dx {
			e += dx
			y0 += sy
		}
	}
}
