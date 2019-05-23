package main

import (
	"image"

	"github.com/ktye/plot"
	"github.com/ktye/ui/base"
)

// pl(ot) device variable
//
// set
//	pl:""                      /reset
//	pl:[plot.Plot]             /set plots
//	pl:[plot.Plots]            /set plots
//	pl:(plot.Plot;plot.Plot…)  /set plots
//	pl:real array              /new xy plot
//	pl:complex array           /new polar plot
//	pl:list of array           /new xy or polar plot with multiple lines
// get
//	pl → plot.Plots or plot.Plot (if len 1).
type pl struct {
	base.Base
	plot.Plots
	*plot.UI
}

func (p *pl) call(x, idx v) v {
	if x == nil {
		if len(p.Plots) == 1 {
			return p.Plots[0]
		}
		return p.Plots
	}
	switch t := x.(type) {
	case string:
		if t == "" { // reset
			p.Plots = plot.Plots{plot.Plot{}}
		} else {
			return nil
		}
	case []complex128:
		p.Plots = plot.Plots{plotData(t)}
	case plot.Plot:
		p.Plots = plot.Plots{t}
	case plot.Plots:
		p.Plots = t
	case l:
		plots := plotList(t) // list of plot.Plot
		if plots == nil {
			plots = dataList(t) // list of []complex128
		}
		if plots == nil {
			return nil
		}
		p.Plots = plots
	default:
		return nil
	}
	p.UI.SetPlots(p.Plots)
	return nil
}
func newpl(cb *base.Button) devvar {
	plots := plot.Plots{plot.Plot{}}
	p := pl{
		Plots: plots,
		UI:    plot.NewUI(plots),
	}
	p.Menu = base.NewMenu([]*base.Button{cb})
	return &p
}

func (p *pl) Draw(dst *image.RGBA, force bool) {
	p.UI.Draw(dst, force)
	p.DrawMenu(dst)
}
func (p *pl) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if r, o := p.MenuMouse(pos, but, dir, mod); o {
		return r
	}
	return p.UI.Mouse(pos, but, dir, mod)
}
func (p *pl) Key(r rune, code uint32, dir int, mod uint32) int {
	if r, o := p.MenuKey(r, code, dir, mod); o {
		return r
	}
	return p.UI.Key(r, code, dir, mod)
}

func plotData(z []complex128) plot.Plot {
	xvec := func(n int) []float64 {
		r := make([]float64, n)
		for i := range r {
			r[i] = float64(i)
		}
		return r
	}
	isreal := true
	for _, u := range z {
		if imag(u) != 0 {
			isreal = false
			break
		}
	}
	p := plot.Plot{Type: plot.Polar}
	x := xvec(len(z))
	var y []float64
	if isreal {
		p.Type = plot.XY
		y = realvec(z)
		z = nil
	}
	p.Lines = []plot.Line{plot.Line{X: x, Y: y, C: z}}
	return p
}
func dataList(t []interface{}) plot.Plots {
	var r plot.Plot
	for i := range t {
		if z, o := t[i].([]complex128); o {
			p := plotData(z)
			if i == 0 {
				r = p
			} else if p.Type == r.Type {
				r.Lines = append(r.Lines, p.Lines[0])
			} else {
				x := p.Lines[0].X
				if p.Type == plot.XY {
					r.Lines = append(r.Lines, plot.Line{X: x, Y: realvec(z)})
				} else {
					r.Lines = append(r.Lines, plot.Line{X: x, C: z})
				}
			}
			r.Lines[i].Id = i
		} else {
			return nil
		}
	}
	return plot.Plots{r}
}
func plotList(t []interface{}) plot.Plots {
	plots := make(plot.Plots, len(t))
	for i := range t {
		if pt, o := t[i].(plot.Plot); o {
			plots[i] = pt
		} else {
			return nil
		}
	}
	return plots
}
func realvec(z []complex128) []float64 {
	r := make([]float64, len(z))
	for i := range z {
		r[i] = real(z[i])
	}
	return r
}

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
