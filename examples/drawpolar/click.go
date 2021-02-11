package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/cmplx"

	"github.com/ktye/plot"
)

// The ClickPlot implements a simpler version than plot.UI
// It treats all clicks as double clicks with shift pressed.
type UI struct {
	p       plot.Plots
	iplots  []plot.IPlotter
	rect    image.Rectangle
	save    *image.RGBA
	invalid bool
}

func PolarClickUI(p plot.Plots) *UI {
	var u UI
	u.SetPlots(p)
	return &u
}
func (u *UI) SetPlots(p plot.Plots) {
	u.p = p
	u.rect = image.Rectangle{}
	u.save = nil
}
func (u *UI) Draw(dst *image.RGBA, force bool) {
	if force && u.save != nil && u.save.Rect == dst.Rect {
		// Forced draw, but unchanged plot.
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
		return
	}
	if u.rect != dst.Rect || u.invalid {
		size := dst.Rect.Size()
		if u.rect != dst.Rect || u.save == nil {
			u.rect = dst.Rect
			h, err := u.p.IPlots(size.X, size.Y, 0)
			if err != nil {
				return
			}
			u.iplots = h
		}
		u.save = plot.Image(u.iplots, nil, size.X, size.Y, 0).(*image.RGBA)
		if u.save == nil {
			u.save = image.NewRGBA(dst.Rect)
		} else {
			u.save.Rect = u.save.Rect.Add(dst.Rect.Min)
		}
		draw.Draw(dst, dst.Rect, u.save, dst.Rect.Min, draw.Src)
	}
}
func (u *UI) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	switch {
	case but == -1 || but == -2:
		return 0 // wheel  -1→1 (up), -2→-1 (down)
	case dir > 0:
		return 0 // mouse down
	case dir < 0:
		return u.mouseUp(pos.Sub(u.rect.Min), but, mod)
	default:
		return 0 // move
	}
}
func (u *UI) Key(r rune, code uint32, dir int, mod uint32) int { return 0 }

func (u *UI) mouseUp(pos image.Point, but int, mod uint32) int {
	if but == 3 {
		return u.removeLast()
	}
	bounds := u.rect
	snapToPoint := false
	if cb, ok := plot.ClickIPlotters(u.iplots, pos.X, pos.Y, bounds.Dx(), bounds.Dy(), 0, snapToPoint); ok {
		if cb.Type == plot.MeasurePoint {
			z := cb.PointInfo.C
			abs := cmplx.Abs(z)
			ang := cmplx.Phase(z) * 180 / math.Pi
			if ang < 0 {
				ang += 360
			}
			fmt.Printf(zfmt+"\n", abs, ang)
			u.setPointID(z)
			return u.invalidate()
		}
	}
	return u.invalidate()
}
func (u *UI) invalidate() int { // plot needs update from iplots due to interactive action
	u.invalid = true
	u.save = nil
	return 1
}
func (u *UI) removeLast() int {
	fmt.Println("remove")
	n := len(u.p[0].Lines)
	if n > 0 {
		u.p[0].Lines = u.p[0].Lines[:n-1]
	}
	return u.invalidate()
}

func (u *UI) setPointID(z complex128) {
	n := len(u.p[0].Lines)
	u.p[0].Lines[n-1].Id = n - 1
	u.p[0].Lines[n-1].Style.Marker.Color = 0
}
