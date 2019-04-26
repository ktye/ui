package main

import (
	"fmt"

	"github.com/ktye/plot"
	"github.com/ktye/ui"
)

func main() {
	p := &ui.Plot{}
	p.SetPlots(xyplot())
	w := ui.New(nil)
	w.Top.W = p
	w.Render()

	for {
		select {
		case e := <-w.Inputs:
			w.Input(e)

		case err, ok := <-w.Error:
			if !ok {
				return
			}
			fmt.Println("ui:", err)
		}
	}
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
