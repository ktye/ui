// Package ui
package ui

import (
	"image"
	"io"
)

func New(d Display) *Window {
	var w Window
	w.Display = d
	w.Quit = make(chan bool)       // send here to quit
	w.Call = make(chan func() int) // send a function here that modifies the ui
	w.image = w.Display.Image()
	w.mouse = w.Display.Mouse()
	w.key = w.Display.Key()
	w.err = w.Err()
	w.dst = <-w.image // wait for the window to be mapped
	return &w
}

type Window struct {
	Display
	Top   Widget
	Quit  chan bool
	Call  chan func() int
	image <-chan *image.RGBA
	mouse <-chan mt
	key   <-chan kt
	err   <-chan error
	dst   *image.RGBA
}
type Display interface {
	Image() chan *image.RGBA
	Mouse() chan mt
	Key() chan kt
	Err() chan error
	Flush()
	Lock()
	Unlock()
}
type Widget interface {
	Draw(dst *image.RGBA, force bool)
	Mouse(pos image.Point, but int, dir int, mod uint32) int
	Key(r rune, code uint32, dir int, mod uint32) int
}

func (w *Window) Run() chan bool {
	w.Draw(true)
	done := make(chan bool)
	go func() {
		for {
			draw := 0
			select {
			case m := <-w.mouse:
				draw = w.Top.Mouse(m.Pos, m.But, m.Dir, m.Mod)
			case k := <-w.key:
				draw = w.Top.Key(k.Rune, k.Code, k.Dir, k.Mod)
			case w.dst = <-w.image:
				draw = -1
			case e := <-w.err:
				if e == io.EOF {
					done <- true
					return
				}
				println("ui error:", e.Error())
			case f := <-w.Call:
				draw = f()
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
	But int    // 0(none),   1(left),    2(middle), 3(right), -1(wheel up), -2(wheel down)
	Dir int    // 0(motion)  1(press),  -1(release) 3(step,both)
	Mod uint32 // 1(shift),  2(cntrl),   4(alt),    8(meta) (bitmask)
}
type kt = struct {
	Rune rune
	Code uint32 // see golang.org/x/mobile/event/key/key.go
	Dir  int    // 0(hold),  1(press), -1(release)
	Mod  uint32 // 1(shift), 2(cntrl),  4(alt),    8(meta) (bitmask)
}
