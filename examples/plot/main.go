package main

import (
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/fonts/apl385"
)

func main() {
	base.SetFont(apl385.TTF(), 20)
	setPlotFont()

	scl := base.NewScale(plot.NewUI(xyplot()))
	scl.Funcs = append(scl.Funcs, setPlotFont)

	win := ui.New(dpy.New(nil))
	win.Top = scl
	done := win.Run()
	<-done
}

func xyplot() plot.Plots {
	l := func(y []float64) plot.Line {
		x := make([]float64, len(y))
		for i := range x {
			x[i] = float64(i)
		}
		return plot.Line{X: x, Y: y}
	}
	return plot.Plots{plot.Plot{
		Type:   plot.XY,
		Title:  "xy-plot",
		Xlabel: "x",
		Ylabel: "y",
		Lines: []plot.Line{
			l([]float64{0, 3, 1, 2, 4}),
			l([]float64{2, 1, 1, 4, 2}),
		},
	}}
}

// The plot widget is implemented externally.
// It does it's own font handling and is notified when the ui font changes size (mouse wheel).
func setPlotFont() {
	s1 := base.Font.Size()
	s2 := (s1 * 8) / 10
	if s1 < 6 {
		s1 = 6
	}
	if s2 < 6 {
		s2 = 6
	}
	f1 := base.LoadFace(base.Font.TTF, s1)
	f2 := base.LoadFace(base.Font.TTF, s2)
	plot.SetFonts(f1, f2)
}
