# ui - go gui toolkit

## Simple widget interface (image is the only package dependency)
```go
type Widget interface {
	Draw(dst *image.RGBA, force bool)
	Mouse(pos image.Point, but int, dir int, mod uint32) int
	Key(r rune, code uint32, dir int, mod uint32) int
}
```

## The display backend needs to provide this interface
```go
type Display interface {
	Image() chan *image.RGBA
	Mouse() chan struct{Pos image.Point, But, Dir int, Mod uint32},
	Key()   chan struct{Rune rune, Code uint32, Dir int, Mod uint32},
	Err()   char error
	Flush()
	Lock()
	Unlock()
}
```
Package dpy provides an implementation based on shiny.

## Basic example with a single widget
```go
package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base" 
	"github.com/ktye/ui/dpy"  // shiny backend
)

var win *ui.Window

func main() {
	win = ui.New(dpy.New(nil))
	win.Top = base.NewScale(&basic{})
	done := win.Run()
	<-done
}

// A basic widget.
type basic struct{}

func (b *basic) Draw(dst *image.RGBA, force bool) {
	println("draw", force)
	if force {
		draw.Draw(dst, dst.Bounds(), &image.Uniform{color.RGBA{0xFF, 0x00, 0xFF, 0xFF}}, image.ZP, draw.Src)
	}
}
func (b *basic) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	println("mouse", but)
	if but == 1 && dir > 0 {
		return 1
	}
	return 0
}
func (b *basic) Key(r rune, code uint32, dir int, mod uint32) int {
	println("key", code)
	if code == 41 { // ESC
		go func() { win.Quit <- true }()
	}
	return 0
}
```
