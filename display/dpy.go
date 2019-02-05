// Display is the shiny backend of ui.
package display

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
	Mouse       chan mouse.Event
	Key         chan key.Event
	Size        chan size.Event
	Err         chan error
	PixelsPerPt float32
	Buffer      screen.Buffer
	screen      screen.Screen
	window      screen.Window
	opt         screen.NewWindowOptions
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
		Mouse: make(chan mouse.Event),
		Key:   make(chan key.Event),
		Size:  make(chan size.Event),
		Err:   make(chan error),
		opt:   *opt,
	}
	go driver.Main(func(s screen.Screen) {
		d.screen = s
		w, err := s.NewWindow(opt)
		if err != nil {
			d.Err <- err
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
			d.Err <- err
			d.Unlock()
			return
		}
		d.Buffer = b
		d.PixelsPerPt = e.PixelsPerPt
		d.Unlock()
		d.Size <- e
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

	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				d.Err <- io.EOF
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
			d.Mouse <- e

		case key.Event:
			d.Key <- e

		case error:
			d.Err <- e
		}
	}
}
