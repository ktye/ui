package ui

import (
	"image"
	"image/color"
)

func (w *Window) defaultColors() {
	w.Disabled = Colors{
		Text:       argb(0x888888ff),
		Background: argb(0xf0f0f0ff),
		Border:     argb(0xe0e0e0ff),
	}
	w.Inverse = Colors{
		Text:       argb(0xeeeeeeff),
		Background: argb(0x3272dcff),
		Border:     argb(0x666666ff),
	}
	w.Selection = Colors{
		Text:       argb(0xeeeeeeff),
		Background: argb(0xbbbbbbff),
		Border:     argb(0x666666ff),
	}
	w.SelectionHover = Colors{
		Text:       argb(0xeeeeeeff),
		Background: argb(0x3272dcff),
		Border:     argb(0x666666ff),
	}
	w.Placeholder = Colors{
		Text:       argb(0xaaaaaaff),
		Background: argb(0xf8f8f8ff),
		Border:     argb(0xbbbbbbff),
	}
	w.Striped = Colors{
		Text:       argb(0x333333ff),
		Background: argb(0xf2f2f2ff),
		Border:     argb(0xbbbbbbff),
	}

	w.Regular = Colorset{
		Normal: Colors{
			Text:       argb(0x333333ff),
			Background: argb(0xf8f8f8ff),
			Border:     argb(0xbbbbbbff),
		},
		Hover: Colors{
			Text:       argb(0x222222ff),
			Background: argb(0xfafafaff),
			Border:     argb(0x3272dcff),
		},
	}
	w.Primary = Colorset{
		Normal: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x007bffff),
			Border:     argb(0x007bffff),
		},
		Hover: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x0062ccff),
			Border:     argb(0x0062ccff),
		},
	}
	w.Secondary = Colorset{
		Normal: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x868e96ff),
			Border:     argb(0x868e96ff),
		},
		Hover: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x727b84ff),
			Border:     argb(0x6c757dff),
		},
	}
	w.Success = Colorset{
		Normal: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x28a745ff),
			Border:     argb(0x28a745ff),
		},
		Hover: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0x218838ff),
			Border:     argb(0x1e7e34ff),
		},
	}
	w.Danger = Colorset{
		Normal: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0xdc3545ff),
			Border:     argb(0xdc3545ff),
		},
		Hover: Colors{
			Text:       argb(0xffffffff),
			Background: argb(0xc82333ff),
			Border:     argb(0xbd2130ff),
		},
	}

	w.BackgroundColor = argb(0xfcfcfcff)
	w.Background = image.NewUniform(argb(0xfcfcfcff))
	w.Gutter = argb(0xbbbbbbff)

	w.ScrollBGNormal = argb(0xf4f4f4ff)
	w.ScrollBGHover = argb(0xf0f0f0ff)
	w.ScrollVisibleNormal = argb(0xbbbbbbff)
	w.ScrollVisibleHover = argb(0x999999ff)
}

func argb(c uint32) color.Color {
	return color.RGBA{
		R: uint8(c >> 24),
		G: uint8(c >> 16),
		B: uint8(c >> 8),
		A: uint8(c),
	}
}
