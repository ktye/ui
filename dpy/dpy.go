// Dpy is the shiny backend of ui.
package dpy

// TODO: this has been observed (once)
// panic: windriver: Buffer.Upload called after Buffer.Release

import (
	"image"
	"io"
	"sync"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

// Display contains a single window and a screen.
// UI Applications are restricted to single-window applications.
// Use New to create the display.
type Display struct {
	sync.Mutex
	mouse       chan mt
	key         chan kt
	image       chan *image.RGBA
	err         chan error
	PixelsPerPt float32
	Buffer      screen.Buffer
	screen      screen.Screen
	window      screen.Window
	opt         screen.NewWindowOptions
	sxl         bool
}

func (d *Display) Image() chan *image.RGBA { return d.image }
func (d *Display) Mouse() chan mt          { return d.mouse }
func (d *Display) Key() chan kt            { return d.key }
func (d *Display) Err() chan error         { return d.err }

type mt = struct {
	Pos image.Point
	But int
	Dir int
	Mod uint32
}
type kt = struct {
	Rune rune
	Code uint32
	Dir  int
	Mod  uint32
}

// New returns a new Display.
func New(opt *screen.NewWindowOptions) *Display {
	if opt == nil {
		opt = &screen.NewWindowOptions{
			Width:  800,
			Height: 600,
			Title:  "UI",
		}
	}
	d := &Display{
		mouse: make(chan mt),
		key:   make(chan kt),
		image: make(chan *image.RGBA),
		err:   make(chan error),
		opt:   *opt,
	}
	go driver.Main(func(s screen.Screen) {
		// todo: this is called twice and creates two windows, no idea why.
		// early returning does not work.
		d.screen = s
		w, err := s.NewWindow(opt)
		if err != nil {
			d.err <- err
			return
		}
		d.window = w
		defer w.Release()

		var b screen.Buffer
		defer func() {
			if b != nil {
				b.Release()
			}
		}()
		d.Buffer = b
		d.loop()
	})
	return d
}

func (d *Display) Flush() {
	d.Lock()
	d.window.Upload(image.Point{}, d.Buffer, d.Buffer.Bounds())
	d.window.Publish()
	d.Unlock()
}

func (d *Display) loop() {
	w := d.window
	b := d.Buffer

	var err error
	resizeFunc := func(e size.Event) {
		d.Lock()
		if b != nil {
			b.Release()
		}
		b, err = d.screen.NewBuffer(e.Size())
		if err != nil {
			d.err <- err
			d.Unlock()
			return
		}
		d.Buffer = b
		d.PixelsPerPt = e.PixelsPerPt
		d.Unlock()
		d.image <- b.RGBA()
	}

	// Delay and filter resize events.
	resize := make(chan size.Event)
	go func(se chan size.Event) {
		var cur size.Event
		delay := 100 * time.Millisecond
		t := time.NewTimer(delay)
		for {
			select {
			case e := <-se:
				cur = e
				t.Reset(delay) // Is that safe?
			case <-t.C:
				resizeFunc(cur)
			}
		}
	}(resize)

	var lastX, lastY float32
	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				d.err <- io.EOF
				return
			}

		case paint.Event:
			if b != nil {
				w.Upload(image.Point{}, b, b.Bounds())
				w.Publish()
			}

		case size.Event:
			// When minimizing a window, it receives a size.Event with zero size.
			if e.WidthPx == 0 {
				continue
			}
			resize <- e

		case mouse.Event:
			dir := int(e.Direction)
			if e.Direction == mouse.DirRelease {
				dir = -1
			}
			// On win32 when not on the primary monitor, wheel event position is wrong.
			// No idea how to fix that in shiny. ScreenToClient seems not to be enough.
			// We store the last non-wheel position instead. Some widgets rely on this to work correctly (e.g. split).
			if e.Button < 0 {
				e.X, e.Y = lastX, lastY
			} else {
				lastX, lastY = e.X, e.Y
			}
			d.mouse <- mt{
				Pos: image.Point{X: int(e.X), Y: int(e.Y)},
				But: int(e.Button),
				Dir: dir,
				Mod: uint32(e.Modifiers),
			}

		case key.Event:
			dir := int(e.Direction)
			if e.Direction == key.DirRelease {
				dir = -1
			}
			d.key <- kt{
				Rune: e.Rune,
				Code: uint32(e.Code),
				Dir:  dir,
				Mod:  uint32(e.Modifiers),
			}

		case error:
			d.err <- e
		}
	}
}
