package paint

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func String(m draw.Image, pt image.Point, c color.Color, face font.Face, s string) image.Point {
	ascent := face.Metrics().Ascent
	dot := fixed.P(pt.X, pt.Y).Add(fixed.Point26_6{Y: ascent})

	drawer := font.Drawer{
		Dst:  m,
		Src:  &image.Uniform{c},
		Face: face,
		Dot:  dot,
	}
	drawer.DrawString(s)
	dx := int(drawer.Dot.Sub(dot).X / 64)
	return pt.Add(image.Point{dx, 0})
}
