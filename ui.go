// Package ui
package ui

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/ktye/ui/display"
	"golang.org/x/exp/shiny/screen"
)

// Window is the top ui element, created by New.
type Window struct {
	*display.Display
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

	mouse, origMouse Mouse
	font             font
	lastMouseWidget  Widget

	keyboard KeyTranslator
}

// New returns a new window.
func New(opt *screen.NewWindowOptions) *Window {
	var w Window
	w.Display = display.New(opt)
	w.Inputs = make(chan Input)
	w.Error = make(chan error)

	// Wait for the window to be mapped. This creates the draw buffer and sets PixelsPerPt.
	<-w.Display.Size

	w.defaultFont()
	w.defaultColors()

	go func() {
		for {
			select {
			case m := <-w.Display.Mouse:
				w.Inputs <- Input{Type: InputMouse, Mouse: m}
			case k := <-w.Display.Key:
				w.Inputs <- Input{Type: InputKey, Key: k}
			case s := <-w.Display.Size:
				w.Inputs <- Input{Type: InputSize, Size: s}
			case e := <-w.Display.Err:
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
	dst := d.Display.Buffer.RGBA()
	if d.Top.Draw == Dirty {
		draw.Draw(dst, dst.Bounds(), d.Background, image.ZP, draw.Src)
	}
	d.Top.W.Draw(d, &d.Top, dst, image.ZP, d.mouse, d.Top.Draw == Dirty)
	d.Top.Draw = Clean
	d.Unlock()
	d.Display.Flush()
}

// Input propagates the input event through the UI tree.
// Mouse and key events are delivered the right UIs.
// Resize is handled by reattaching does a layout and draw.
// Func calls the function.
// Error implies an error from shiny and terminates the program.
func (w *Window) Input(e Input) {
	logInputs := false

	switch e.Type {
	case InputMouse:
		if logInputs {
			fmt.Printf("ui: mouse %+v\n", e.Mouse)
		}
		w.Mouse(e.Mouse)
	case InputKey:
		if logInputs {
			fmt.Printf("duit: key %+v\n", e.Key)
		}
		w.Key(e.Key)
	case InputSize:
		if logInputs {
			fmt.Printf("duit: size")
		}
		w.Resize()
	case InputFunc:
		if logInputs {
			fmt.Printf("duit: func")
		}
		//e.Func()
		//w.Render()
	case InputError:
		if logInputs {
			fmt.Printf("duit: error: %s", e.Error)
		}
	}
}

func (w *Window) Resize() {
	w.Top.Layout = Dirty
	w.Top.Draw = Dirty
	w.Render()
}

// MarkDraw is like MarkLayout, but marks ui as requiring a draw.
func (w *Window) MarkDraw(widget Widget) {
	if widget == nil {
		w.Top.Draw = Dirty
	} else {
		if !w.Top.W.Mark(&w.Top, widget, false) {
			fmt.Printf("duit: markdraw %T: nothing marked\n", widget)
		}
	}
}

func (w *Window) apply(r Result) {
	if r.Warp != nil {
		/* TODO mouse warp
		if err := w.Display.MoveTo(*r.Warp); err != nil {
			fmt.Printf("ui: warp to %v: %s\n", r.Warp, err)
		} else {
			d.mouse.Point = *r.Warp
			d.mouse.Buttons = 0
			d.origMouse = d.mouse
			r = d.Top.UI.Mouse(d, &d.Top, d.mouse, d.origMouse, image.ZP)
		}
		*/
	}
	if r.Hit != w.lastMouseWidget {
		if r.Hit != nil {
			w.MarkDraw(r.Hit)
		}
		if w.lastMouseWidget != nil {
			w.MarkDraw(w.lastMouseWidget)
		}
	}
	w.lastMouseWidget = r.Hit
	w.Render()
}
