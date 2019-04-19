package ui

import (
	"image"
	"image/draw"

	"github.com/ktye/plot"
	"golang.org/x/mobile/event/key"
)

type Plot struct {
	Target **Plot
	p      plot.Plots
}

func (p *Plot) SetPlots(plts plot.Plots) { p.p = plts }

func (p *Plot) Layout(w *Window, self *Kid, sizeAvail image.Point, force bool) {
	self.R = rect(sizeAvail)
	if p.Target != nil {
		*p.Target = p
	}
}

func (p *Plot) Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool) {
	if p.p == nil {
		return
	}
	width, height := self.R.Dx(), self.R.Dy()
	ip, err := p.p.IPlots(width, height)
	if err != nil {
		println("plot.Draw:", err.Error())
		return
	}
	src := plot.Image(ip, nil, width, height)
	draw.Draw(img, image.Rectangle{orig, orig.Add(image.Point{width, height})}, src, image.ZP, draw.Src)
}

func (p *Plot) Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result) {
	return r
}
func (p *Plot) Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (r Result) { return r }
func (p *Plot) FirstFocus(w *Window, self *Kid) (warp *image.Point)                         { return nil }
func (p *Plot) Focus(w *Window, self *Kid, o Widget) (warp *image.Point) {
	if o != p {
		return nil
	}
	return p.FirstFocus(w, self)
}
func (p *Plot) Mark(k *Kid, o Widget, forLayout bool) (marked bool) { return k.Mark(o, forLayout) }
