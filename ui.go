// Package ui
package ui

import (
	"io"

	"github.com/eaburns/T/clipboard"
	"github.com/ktye/ui/dpy"
	"golang.org/x/exp/shiny/screen"
)

// Window is the top ui element, created by New.
type Window struct {
	*dpy.Display
	Quit      chan bool
	Call      chan func()
	Top       Widget
	font      font
	clipboard clipboard.Clipboard
}

// New returns a new window.
func New(opt *screen.NewWindowOptions) *Window {
	var w Window
	w.Display = dpy.New(opt)
	w.Quit = make(chan bool)   // send here to quit
	w.Call = make(chan func()) // send a function here that modifies the ui
	<-w.Display.Size           // wait for the window to be mapped
	w.defaultFont()
	w.clipboard = clipboard.New()

	return &w
}
func (w *Window) Run() chan bool {
	w.Draw(true)
	done := make(chan bool)
	go func() {
		for {
			draw := 0
			select {
			case m := <-w.Display.Mouse:
				draw = w.mouseScale(m)
				if draw == 0 {
					draw = w.Top.Mouse(m.Pos, m.But, m.Dir, m.Mod)
				}
			case k := <-w.Display.Key:
				draw = w.Top.Key(k.Rune, k.Code, k.Press, k.Mod)
			case <-w.Display.Size:
				draw = -1
			case f := <-w.Call:
				f()
			case e := <-w.Display.Err:
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
	dst := w.Display.Buffer.RGBA()
	w.Top.Draw(dst, force)
	w.Unlock()
	w.Display.Flush()
}
