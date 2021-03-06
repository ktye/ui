package base

import (
	"image"

	"github.com/ktye/ui"
)

// NewStack returns a stacked widget container.
// Widgets can be pushed on top, or the top widget can be removed, which restores
// the underlying widget.
// If the esc parameter is set for a pushed widget, it calls pop when the esc key is pressed.
func NewStack(w ui.Widget) *stack {
	return &stack{b: []ui.Widget{w}}
}

type stack struct {
	b []ui.Widget
	e []bool
}

func (s *stack) Push(w ui.Widget, esc bool) {
	s.b = append(s.b, w)
	s.e = append(s.e, esc)
}
func (s *stack) Pop() {
	if len(s.b) > 1 {
		s.b = s.b[:len(s.b)-1]
		s.e = s.e[:len(s.e)-1]
	}
}
func (s *stack) Draw(dst *image.RGBA, force bool) {
	s.b[len(s.b)-1].Draw(dst, force)
}
func (s *stack) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	return s.b[len(s.b)-1].Mouse(pos, but, dir, mod)
}
func (s *stack) Key(r rune, code uint32, dir int, mod uint32) int {
	if s.e[len(s.e)-1] && code == 41 && dir < 0 {
		s.Pop()
		return -1
	}
	return s.b[len(s.b)-1].Key(r, code, dir, mod)
}
