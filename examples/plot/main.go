package main

import (
	"github.com/ktye/plot"
	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/dpy"
	"github.com/ktye/ui/fonts/apl385"
	eplot "github.com/ktye/ui/plot"
)

func main() {
	// A button bar switches between a simple and an extended example.
	// The simple example shows a single interactive plot widget, that is provided by ktye/plot.UI.
	// The extended examples links a caption table to the plot. This is provided by ktye/ui/plot.Widget.
	var simple, extended ui.Widget
	simple = simplePlot()
	extended = extendedPlot()
	var b *base.ButtonBar
	but := base.NewButtonBar(simple, []*base.Button{
		base.NewButton("simple", "", func() int { b.Kid.Widget = simple; return -1 }),
		base.NewButton("extended", "", func() int { b.Kid.Widget = extended; return -1 }),
	})
	but.Target = &b

	base.SetFont(apl385.TTF(), 20)
	setPlotFont()
	scl := base.NewScale(but)
	scl.Funcs = append(scl.Funcs, setPlotFont)

	win := ui.New(dpy.New(nil))
	win.Top = scl
	done := win.Run()
	<-done
}

func simplePlot() *plot.UI {
	return plot.NewUI(plot.Plots{xyplot()})
}
func extendedPlot() *eplot.Widget {
	p1 := xyplot()
	p1.Caption = caption("channel A", []float64{2.3, 1.23, 5.21})
	p2 := xyplot()
	p2.Caption = caption("channel B", []float64{6.28, 2.78, 5.32})
	return eplot.New(plot.Plots{p1, p2})
}

func xyplot() plot.Plot {
	l := func(id int, y []float64) plot.Line {
		x := make([]float64, len(y))
		for i := range x {
			x[i] = float64(i)
		}
		ln := plot.Line{X: x, Y: y, Id: id}
		ln.Style.Line.Width = 2
		return ln
	}
	return plot.Plot{
		Type:   plot.XY,
		Title:  "xy-plot",
		Xlabel: "x",
		Ylabel: "y",
		Lines: []plot.Line{
			l(0, []float64{0, 3, 1, 2, 4}),
			l(1, []float64{2, 1, 1, 4, 2}),
			l(2, []float64{1, 2, 4, 2, 3}),
		},
	}
}
func caption(name string, values []float64) *plot.Caption {
	return &plot.Caption{
		Columns: []plot.CaptionColumn{
			plot.CaptionColumn{
				Name: name,
				Data: values,
				Unit: "mV",
			},
			plot.CaptionColumn{
				Name: "File",
				Data: []string{"file1", "file2", "file3"},
			},
			plot.CaptionColumn{
				Name: "Tags",
				Data: []string{"abc", "def", "ghi"},
			},
		},
	}
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
