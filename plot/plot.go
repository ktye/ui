// plot provides a widget that links github.com/ktye/plot.UI with a caption table.
package plot

import (
	"image"

	plt "github.com/ktye/plot"
	"github.com/ktye/ui/base"
)

// Widget combines plot.UI with a caption table in a list box.
type Widget struct {
	Target **Widget
	*base.Split
	*plt.UI
	*base.List
	c *plt.Caption
}

func New(p plt.Plots) *Widget {
	w := Widget{
		UI:    &plt.UI{},
		Split: &base.Split{},
		List:  &base.List{},
	}
	w.Split.Kids[0] = base.Kid{Widget: w.UI}
	w.Split.Kids[1] = base.Kid{Widget: w.List}
	w.SetPlots(p)
	return &w
}

func (w *Widget) SetPlots(p plt.Plots) {
	w.UI.SetPlots(p)
	w.c = nil
	w.setCaption(p)
	w.Split.Dirty = true
}

func (w *Widget) Draw(dst *image.RGBA, force bool) {
	if w.Target != nil {
		*w.Target = w
	}
	if w.c == nil {
		w.UI.Draw(dst, force)
		return
	}
	w.Split.Draw(dst, force)
}
func (w *Widget) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if w.c == nil {
		return w.UI.Mouse(pos, but, dir, mod)
	}
	return w.Split.Mouse(pos, but, dir, mod)
}
func (w *Widget) Key(r rune, code uint32, dir int, mod uint32) int {
	if w.c == nil {
		return w.UI.Key(r, code, dir, mod)
	}
	return w.Split.Key(r, code, dir, mod)
}
