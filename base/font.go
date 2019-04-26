package base

import (
	"image"

	imfont "golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"

	"github.com/golang/freetype/truetype"
)

var Font = font{ttf: nil, size: 13, Face: basicfont.Face7x13}

type font struct {
	imfont.Face
	ttf  []byte
	size int
}

func SetFont(ttf []byte, size int) {
	f, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	Font.ttf = ttf
	opt := truetype.Options{
		Size: float64(size),
		DPI:  72.0,
	}
	if opt.Size == 0 {
		opt.Size = 18
	}
	Font.size = size
	Font.Face = truetype.NewFace(f, &opt)
}

func FontHeight() int {
	return Font.Face.Metrics().Height.Ceil()
}

func StringSize(s string) image.Point {
	dx := int(imfont.MeasureString(Font.Face, s)+32) / 64
	dy := FontHeight()
	return image.Point{dx, dy}
}
