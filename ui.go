// Package ui
package ui

import (
	"image"
	"image/draw"
	"io"
)

// Window is the top ui element, created by New.
type Window struct {
	Display
	Top   Widget
	Quit  chan bool
	Call  chan func()
	image <-chan draw.Image
	mouse <-chan mt
	key   <-chan kt
	err   <-chan error
	dst   draw.Image
}
type Display interface {
	Image() chan draw.Image // Resize
	Mouse() chan struct {
		Pos image.Point
		But int    //  0(none),   1(left),    2(middle), 3(right), -1(wheel up), -2(wheel down)
		Dir int    // +1(press), -1(release), 0(motion)
		Mod uint32 // see godoc.org/golang.org/x/mobile/event/key#Modifiers
	}
	Key() chan struct {
		Rune  rune
		Code  uint32
		Press bool
		Mod   uint32
	}
	Err() chan error
	Flush()
	Lock()
	Unlock()
}
type Widget interface {
	Draw(dst draw.Image, force bool)
	Mouse(pos image.Point, but int, dir int, mod uint32) int
	Key(r rune, code uint32, press bool, mod uint32) int
}

// New returns a new window.
func New(d Display) *Window {
	var w Window
	w.Display = d
	w.Quit = make(chan bool)   // send here to quit
	w.Call = make(chan func()) // send a function here that modifies the ui
	w.image = w.Display.Image()
	w.mouse = w.Display.Mouse()
	w.key = w.Display.Key()
	w.err = w.Err()
	w.dst = <-w.image // wait for the window to be mapped
	return &w
}
func (w *Window) Run() chan bool {
	w.Draw(true)
	done := make(chan bool)
	go func() {
		for {
			draw := 0
			select {
			case m := <-w.mouse:
				draw = mouseScale(m.But, m.Mod)
				if draw == 0 {
					draw = w.Top.Mouse(m.Pos, m.But, m.Dir, m.Mod)
				}
			case k := <-w.key:
				draw = w.Top.Key(k.Rune, k.Code, k.Press, k.Mod)
			case w.dst = <-w.image:
				draw = -1
			case f := <-w.Call:
				f()
			case e := <-w.err:
				if e == io.EOF {
					done <- true
					return
				}
				println("ui error:", e.Error())
			case <-w.Quit:
				done <- true
				return
			}
			if draw != 0 {
				w.Draw(draw < 0)
			}
		}
	}()
	return done
}
func (w *Window) Draw(force bool) {
	w.Lock()
	w.Top.Draw(w.dst, force)
	w.Unlock()
	w.Flush()
}

type mt = struct {
	Pos image.Point
	But int
	Dir int
	Mod uint32
}
type kt = struct {
	Rune  rune
	Code  uint32
	Press bool
	Mod   uint32
}
