package main

import (
	"github.com/ktye/plot"
	"github.com/ktye/ui/base"
)

func regplot(a map[v]v) {
	p := plot.Plot{}
	setPlotFonts()
	p.Style.Dark = true
	a["plot"] = p
	a["line"] = plot.Line{}
	a["plots"] = plots
}

// plots converts a list of plots to plot.Plots which results in multiple plots next to each other as a single image.
func plots(x v) v {
	l, o := x.(l)
	if !o {
		panic("type")
	}
	plts := make(plot.Plots, len(l))
	for i := range l {
		p, o := l[i].(plot.Plot)
		if !o {
			panic("type")
		}
		plts[i] = p
	}
	return plts
}

func setPlotFonts() {
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
