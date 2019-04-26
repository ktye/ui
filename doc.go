// Package ui is a lightweight gui toolkit for go.
//
// It is composed of three components:
// - the display backend (dpy/ is an implementation based on shiny)
// - a collection of widgets, (base/ provides some)
// - this package, which connects both and defines the interface
//
// The basic idea is that external packages that implement a widget, do not need to
// import anything besides image and draw.
// This is also why ui does not define types for Mouse and Key events, or depends on golang.org/x/mobile/event.
//
// All a widget has to do, is to implement:
// 	interface {
//		Draw(dst draw.Image, force bool)
//		Mouse(pos image.Point, but int, dir int, mod uint32) int
//		Key(r rune, code uint32, press bool, mod uint32) int
// 	}
// Draw asks the widget to draw itself on img, which is a subimage of the window.
// The widget should only redraw, if force is true, or if it requested the drawing itself from Mouse or Key.
//
// Mouse is sent for a mouse event.
// The return value is 0 if it is ignored, >0 if the widget wants to redraw and <0 if the complete
// ui should be redrawn.
//
// For Key rune, code, mod and Mouse but, see golang.org/x/mobile/event/{key,mouse}
// Key press is press or release, Mouse dir is >0(press), <0(release), 0(motion).
package ui
