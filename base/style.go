package base

import (
	"image"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
)

var Font = fnt{TTF: nil, size: 13, metrics: basicfont.Face7x13.Metrics(), Face: basicfont.Face7x13}

type fnt struct {
	font.Face
	TTF     []byte
	size    int
	metrics font.Metrics
}

func (f fnt) Size() int {
	return f.size
}

func SetFont(ttf []byte, size int) {
	Font.TTF = ttf
	Font.Face = LoadFace(ttf, size)
	Font.size = size
	Font.metrics = Font.Face.Metrics()
}

func LoadFace(ttf []byte, size int) font.Face {
	f, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	opt := truetype.Options{
		Size: float64(size),
		DPI:  72.0,
	}
	return truetype.NewFace(f, &opt)
}

func StringSize(s string) image.Point {
	dx := int(font.MeasureString(Font.Face, s)+32) / 64
	dy := Font.size
	return image.Point{dx, dy}
}

func String(dst *image.RGBA, at image.Point, s string) int {
	f := Font.Face
	fg := Colors[0]
	dot := fixed.P(dst.Rect.Min.X, dst.Rect.Min.Y).Add(fixed.P(at.X, at.Y)).Add(fixed.Point26_6{0, Font.metrics.Ascent})
	start := dot.X
	b := rune(-1)
	for _, c := range s {
		if c == '\n' || c == '\t' {
			c = ' '
		}
		if b >= 0 {
			dot.X += f.Kern(b, c)
		}
		if dot.X.Ceil() > dst.Rect.Max.X {
			break
		}
		dr, m, p, a, ok := f.Glyph(dot, c)
		if !ok {
			continue
		}
		draw.DrawMask(dst, dr, fg, image.Point{}, m, p, draw.Over)
		dot.X += a
		b = c
	}
	return (start - dot.X).Ceil()
}

type Colorset [2]*image.Uniform

var Colors Colorset = [2]*image.Uniform{Black.Uniform(), White.Uniform()}
var ListCurCol Colorset = [2]*image.Uniform{Black.Uniform(), LightGrey.Uniform()}
var ListCurSel Colorset = [2]*image.Uniform{White.Uniform(), Blue.Uniform()}
var ListSelect Colorset = [2]*image.Uniform{White.Uniform(), LightBlue.Uniform()}
var ButtonCol Colorset = [2]*image.Uniform{White.Uniform(), LightBlue.Uniform()}
var ButtonPrim Colorset = [2]*image.Uniform{White.Uniform(), DeepBlue.Uniform()}
var SplitGutter, SplitGutterActive = Amber, Orange

type RGB uint32

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = (0xFF0000 & uint32(c)) >> 16
	r |= r << 8
	g = (0xFF00 & uint32(c)) >> 8
	g |= g << 8
	b = 0xFF & uint32(c)
	b |= b << 8
	a = 0xFFFF
	return
}
func (c RGB) Uniform() *image.Uniform {
	return &image.Uniform{c}
}

func NewColors(f, b RGB) Colorset {
	return [2]*image.Uniform{&image.Uniform{f}, &image.Uniform{b}}
}

const ( //materialuicolors.co/ 400
	White, Black, Grey, LightGrey     RGB = 0xFFFFFF, 0x000000, 0x78909C, 0xBDBDBD
	Red, Green, Blue, DeepBlue        RGB = 0xEF5350, 0x66BB6A, 0x42A5F5, 0x1E88E5
	LightBlue, Cyan, Indigo, Teal     RGB = 0x29B6F6, 0x26C6DA, 0x5C6BC0, 0x26A69A
	Yellow, Amber, Orange, DeepOrange RGB = 0xFFEE58, 0xFFCA28, 0xFFA726, 0xFF7043
	LightGreen, Lime, Brown           RGB = 0x9CCC65, 0xD4E157, 0x8D6E63
	Pink, Purple, DeepPurple          RGB = 0xEC407A, 0xAB47BC, 0x7E57C2
)
