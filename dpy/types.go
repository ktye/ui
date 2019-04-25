package dpy

import "image"

type Mouse struct {
	Pos image.Point
	But int    // 0(none), 1(left), 2(middle), 3(right), -1(wheel up), -2(wheel down)
	Dir int    // +1(press), -1(release), 0(motion)
	Mod uint32 // see Key mod
}

// Key event, see: golang.org/x/mobile/event/key
type Key struct {
	Rune  rune
	Code  uint32
	Press bool
	Mod   uint32
}

type Size image.Point
