package main

import (
	"github.com/ktye/plot"
	"github.com/ktye/ui/base"
)

func newpl(cb *base.Button) devvar {
	return nil
}

/*
func isplot(x v) (plot.Plots, bool) {
	if p, o := x.(plot.Plots); o {
		return p, true
	} else if p, o := x.(plot.Plot); o {
		return plot.Plots{p}, true
	} else if p, o := x.([]plot.Plot); o {
		return plot.Plots(p), true
	} else if u, o := x.(l); o {
		for i := range u {
			if _, o := u[i].(plot.Plot); !o {
				return nil, false
			}
		}
		p = make(plot.Plots, len(u))
		for i := range u {
			p[i] = u[i].(plot.Plot)
		}
		return p, true
	}
	return nil, false
}
func setplot(p plot.Plots) {
	cnv = plot.NewUI(p)
	if sp1.Ratio > 0.95 {
		sp1.Ratio = 0
	}
	if sp2.Ratio > 0.95 {
		sp2.Ratio = 0
	}
	sp2.Kids[1].Widget = cnv
}
*/

func plotfont() {
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

func regplot(a map[v]v) {
	p := plot.Plot{}
	plotfont()
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
