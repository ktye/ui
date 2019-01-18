// Package ui
package ui

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"log"

	uidraw "github.com/ktye/ui/draw"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/mouse"
)

// Window is the top ui element, created by New.
type Window struct {
	*uidraw.Display
	Inputs chan Input
	Error  chan error
	Call   chan func()
	Top    Kid

	Disabled, Inverse, Selection, SelectionHover, Placeholder, Striped Colors
	Regular, Primary, Secondary, Success, Danger                       Colorset

	BackgroundColor color.Color
	Gutter          color.Color
	Background      image.Image

	ScrollBGNormal, ScrollBGHover, ScrollVisibleNormal, ScrollVisibleHover color.Color

	mouse, origMouse mouse.Event
}

// New returns a new window.
func New(opt *screen.NewWindowOptions) *Window {
	var w Window
	w.Display = uidraw.New(opt)
	w.Inputs = make(chan Input)
	w.Error = make(chan error)
	w.defaultColors()

	go func() {
		for {
			select {
			case m := <-w.Display.Mouse:
				w.Inputs <- Input{Type: InputMouse, Mouse: m}
			case k := <-w.Key:
				w.Inputs <- Input{Type: InputKey, Key: k}
			case e := <-w.Err:
				if e == io.EOF {
					close(w.Error)
					return
				}
				w.Inputs <- Input{Type: InputError, Error: e}
			}
		}
	}()

	return &w
}

// Render calls Layout followed by Draw.
// This only does a layout/draw for UIs marked as needing it. If you want to force a layout/draw, mark the top window as requiring a layout/draw.
func (d *Window) Render() {
	d.Layout()
	d.Draw()
}

// Layout the entire UI tree, as necessary.
// Only UIs marked as requiring a layout are actually layed out.
// UIs that receive a layout are marked as requiring a draw.
func (d *Window) Layout() {
	if d.Top.Layout == Clean {
		return
	}
	d.Lock()
	size := d.Display.Buffer.Size()
	d.Unlock()
	d.Top.W.Layout(d, &d.Top, size, d.Top.Layout == Dirty)
	d.Top.Layout = Clean
}

// Draw the entire window, as necessary.
// Only widgets marked as requiring a draw are actually drawn, and their children.
func (d *Window) Draw() {
	if d.Top.Draw == Clean {
		return
	}
	d.Lock()
	defer d.Unlock()

	dst := d.Display.Buffer.RGBA()
	if d.Top.Draw == Dirty {
		draw.Draw(dst, dst.Bounds(), d.Background, image.ZP, draw.Src)
	}
	d.Top.W.Draw(d, &d.Top, dst, image.ZP, d.mouse, d.Top.Draw == Dirty)
	d.Top.Draw = Clean
	// d.Display.Flush() ?
}

// Input propagates the input event through the UI tree.
// Mouse and key events are delivered the right UIs.
// Resize is handled by reattaching does a layout and draw.
// Func calls the function.
// Error implies an error from shiny and terminates the program.
func (d *Window) Input(e Input) {
	logInputs := true

	switch e.Type {
	case InputMouse:
		if logInputs {
			log.Printf("ui: mouse %+v\n", e.Mouse)
		}
		d.Mouse(e.Mouse)
	case InputKey:
		if logInputs {
			log.Printf("duit: key %+v\n", e.Key)
		}
		//d.Key(e.Key)
	case InputResize:
		if logInputs {
			log.Printf("duit: resize")
		}
		//d.Resize()
	case InputFunc:
		if logInputs {
			log.Printf("duit: func")
		}
		//e.Func()
		//d.Render()
	case InputError:
		if logInputs {
			log.Printf("duit: error: %s", e.Error)
		}
	}
}

// Mouse delivers a mouse event to the UI tree.
// Mouse is typically called by Input.
func (d *Window) Mouse(m mouse.Event) {
	/* TODO
	if m.Buttons == 0 || d.origMouse.Buttons == 0 {
		d.origMouse = m
	}
	*/
	d.mouse = m
	r := d.Top.W.Mouse(d, &d.Top, m, d.origMouse, image.ZP)
	d.apply(r)
}

func (d *Window) apply(r Result) {
	// TODO
}
