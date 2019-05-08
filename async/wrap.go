// Async wraps a ui widget to draw asynchronously.
package async

import (
	"image"

	"github.com/ktye/ui"
)

type Wrap struct {
	ui.Widget
	buf  *image.RGBA
	req  image.Rectangle
	lock bool
	call chan func() int
}

func NewWrap(w ui.Widget, call chan func() int) *Wrap {
	return &Wrap{Widget: w, call: call}
}

/*
TODO: draw slow widgets in the background.
If new requestes arrive, always store the rect of the last.
When the background drawing is finished, set image to the buffer and send a call function
to request a redraw.
If the requested size has updated, draw again.
func (w *Wrap) Draw(dst *image.RGBA, force bool) {
	if w.buf == nil {
		w.buf = image.NewRGBA(image.Rectangle{Max: dst.Rect.Size})
	}
	req := dst.Bounds
	if !w.lock {
	}
}
*/

func (w *Wrap) draw() {

}
