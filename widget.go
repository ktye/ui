package ui

import (
	"image"
	"image/draw"

	"golang.org/x/mobile/event/key"
)

// Widget is the interface implemented by a user interface element. For example Button, List, Grid, Scroll.
// Widgets must be able to layout themselves, draw themselves, handle mouse events, key presses, deal with focus requests.
// Widgets also help with propagating UI state and logging.
// For container widgets (those that mostly just contain other widgets), many of these functions can be implemented by a single call to the corresponding Kids*-function.
//
// Layout asks the widget to layout itself and its children in `availSize`.
// Layout must check `self.Layout` and `force`.
// If force is set, it must layout itself and its kids, and pass on force.
// Else, if self.Layout is DirtyKid, it only needs to call Layout on its kids (common for layout UIs).
// The UI can lay itself out beyond size.Y, not beyond size.X.
// size.Y is the amount of screen real estate that will still be visible.
// Layout must update self.Draw if it needs to be drawn after.
// Layout must update self.R with a image.ZP-origin image.Rectangle of the size it allocated.
//
// Draw asks the widget to draw itself on `img`, with `orig` as offset and `m` as the current mouse (for hover states)
// as self.Kid indicates, and pass further Draw calls on to its children as necessary.
// If `force` is set, the widget must draw itself, overriding self.Draw.
//
// Mouse tells the widget about mouse movement over it.
// Layout widgets are in charge of passing these mouse events to their children.
// `self.Layout` and `self.Draw` can be updated if the mouse event resulted in widgets needing relayout/redraw.
// Again it's layout UI's responsibility to propagate requests from self.Layout and self.Draw to its parent.
// `m` is the current mouse state, relative to this UIs zero point.
// `origM` is the mouse of first button down change, to facilitate tracking dragging. If no button is down, origM is the same as m.
// `orig` is the origin location of this UI. If you want to warp the mouse, add the origin to the UI-relative point.
// Result is used to communicate results of the event back to the top.
//
// Key tells the UI about a key press over it.
// Like in Mouse, `self.Layout` and `self.Draw` can be updated.
// `k` is the key pressed. There are no key down/up events, only keys typed.
// See the Key-constants in the draw library for use special keys like the arrow keys,
// function keys and combinations with the cmd key.
// `m` is the mouse location at the time of the key, relative to this UIs zero point.
// `orig` is the origin location of this UI. If you want to warp the mouse, add the origin to the UI-relative point.
//
// FirstFocus returns where the focus should go next when "tab" is hit, if anything.
//
// Focus returns the focus-point for `ui`.
//
// Mark looks for ui (itself or children), marks it as dirty for layout or draw (forLayout),
// and propagates whether it marked anything back to the caller.
type Widget interface {
	Layout(w *Window, self *Kid, sizeAvail image.Point, force bool)
	Draw(w *Window, self *Kid, img draw.Image, orig image.Point, m Mouse, force bool)
	Mouse(w *Window, self *Kid, m Mouse, origM Mouse, orig image.Point) (r Result)
	Key(w *Window, self *Kid, k key.Event, m Mouse, orig image.Point) (r Result)
	FirstFocus(w *Window, self *Kid) (warp *image.Point)
	Focus(w *Window, self *Kid, o Widget) (warp *image.Point)
	Mark(k *Kid, o Widget, forLayout bool) (marked bool)
}
