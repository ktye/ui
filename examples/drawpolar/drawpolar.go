package main

import (
	"flag"
	"fmt"

	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/fonts/apl385"
)

// mouse clicks draw points into a polar diagram and print the values to stdout.

var zfmt string
var rmax float64

func main() {
	flag.Float64Var(&rmax, "rmax", 1, "radial value for outer ring")
	flag.StringVar(&zfmt, "zfmt", "%.1fa%.1f", "abs/angle format")
	flag.Parse()

	var p ui.Widget = polarPlot()

	win := ui.New(dpy.New(nil))
	base.SetFont(apl385.TTF(), 20)
	setPlotFont()

	win.Top = p
	done := win.Run()
	<-done
}

func polarPlot() ui.Widget {
	return PolarClickUI(plot.Plots{
		plot.Plot{
			Type:   plot.Polar,
			Title:  fmt.Sprintf("click to draw points (delete:right), rmax %v", rmax),
			Limits: plot.Limits{Ymax: rmax},
		}})
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
