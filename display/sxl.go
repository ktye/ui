package display

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/geom"
)

// New6 returns a Display for a terminal session.
// Images are written in sixel encoding and mouse events are decoded from escape sequences.
func New6(opt *screen.NewWindowOptions) *Display {
	if opt.Title != "" {
		fmt.Printf("\x1b]0;%s\x07", opt.Title)
	}
	d := &Display{
		Mouse: make(chan mouse.Event),
		Key:   make(chan key.Event),
		Size:  make(chan size.Event),
		Err:   make(chan error),
		opt:   *opt,
		sxl:   true,
	}
	go sxl(d)
	return d
}

// Run terminal in raw mode
//	stty raw
// or stty -icanon min 1 time 0
// TODO: can this be done with an escape sequence?

func sxl(d *Display) {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
	}()

	w := func(s string) { os.Stdout.WriteString(s) }
	defer w("\x1b[0m") // reset
	w("\x1b[?7766l")   // hide scroll bar \x1b[?7766h to show it again
	w("\x1b[?1000h")   // request mouse tracking
	w("\x1b[?1003h")   // request motion tracking (with and without clicks)
	w("\x1b[?1006h")   // request mouse tracking (SGR_EXT_MODE_MOUSE)

	var b byte
	in := bufio.NewReader(os.Stdin)
	csi := &csi{
		r:   in,
		buf: bytes.NewBuffer(nil),
	}

	w("\x1b[18t") // report text area size in cells, response should be csi 8;HEIGHT;WIDTH;t
	w("\x1b[14t") // report text area size in pixels, response should be csi 4;HEIGHT;WIDTH;t
	for {
		b, err = in.ReadByte()
		fmt.Printf("%d 0x%d %q\n", b, b, string(b))
		if err != nil {
			return
		}

		if b == 0x1b {
			b, err = in.ReadByte()
			if err != nil {
				return
			}
			if b != '[' {
				fmt.Printf("ignore 0x%x\n", b)
			} else {
				err = csi.read()
				if err != nil {
					return
				}
				if csi.mouse() {
					fmt.Printf("mouse event: %+v\n", csi.m)
					// TODO d.Mouse <- csi.m
				} else if csi.size() {
					fmt.Printf("size event: %+v\n", csi.s)
					// TODO d.Size <- csi.s
				} else {
					fmt.Printf("unknown csi: %q\n", string(csi.buf.Bytes()))
				}
			}
		} else if b == 0x03 {
			err = fmt.Errorf("sxl: ctrl-c")
			return
		} else {
			fmt.Printf("send key %q\n", string(b))
			// d.Key <- ...
		}
	}
}

// Flush6 writes a sixel encoded frame to stdout.
func (d *Display) flush6() {
	fmt.Println("sxl: flush")
}

type csi struct {
	r             *bufio.Reader
	buf           *bytes.Buffer
	m             mouse.Event
	s             size.Event
	width, height int
	rows, cols    int
}

// Read reads a CSI sequence from the reader.
// 0x1b 0x91 are already consumed.
// The result is written to buf.
func (c *csi) read() error {
	c.buf.Reset()
	for i := 0; ; i++ {
		b, err := c.r.ReadByte()
		if err != nil {
			return err
		}
		c.buf.WriteByte(b)
		if b >= 0x40 && b <= 0x7E {
			return nil
		}
	}
}

// Size decodes a size event from the buffer in the form "4;HEIGHT;WIDTHt"
// It also scans for "8;COLS;ROWS;t" (character cells) and stores it in c.rows, c.cols.
// The function returns true, only for a text area size event.
func (c *csi) size() bool {
	var w, h int
	if n, err := fmt.Sscanf(string(c.buf.Bytes()), "8;%d;%dt", &h, &w); n == 2 && err == nil {
		c.cols = w
		c.rows = h
	}
	if n, err := fmt.Sscanf(string(c.buf.Bytes()), "4;%d;%dt", &h, &w); n != 2 || err != nil {
		return false
	}
	c.width = w
	c.height = h
	fmt.Println("width", w, "height", h)
	c.s = size.Event{
		WidthPx:     w,
		HeightPx:    h,
		WidthPt:     geom.Pt(w),
		HeightPt:    geom.Pt(h),
		PixelsPerPt: 1,
		Orientation: size.OrientationLandscape,
	}
	return true // only size events are
}

// Mouse decodes a mouse event from the buffer.
// It must be called after a successfull read.
func (c *csi) mouse() bool {
	// invisible-island.net/xterm/ctlseqs/ctlseqs.html#h2-Mouse-Tracking
	// Mouse press button 0:   <0;85;25M  at position 85 25
	// Mouse release button 3: <3;12;34m  at position 12 34
	// Wheel up                <64;41;10M
	// Wheel down              <65;1;2M
	var but, x, y int
	var dir rune
	if n, err := fmt.Sscanf(string(c.buf.Bytes()), "<%d;%d;%d%c", &but, &x, &y, &dir); n != 4 || err != nil {
		return false
	}

	if c.rows == 0 || c.cols == 0 {
		return false
	}

	// TODO modifiers.
	c.m.Button = 0
	switch a := but; {
	case a >= 0 && a < 4:
		c.m.Button = mouse.Button(a + 1)
	case a == 64:
		c.m.Button = mouse.ButtonWheelUp
	case a == 65:
		c.m.Button = mouse.ButtonWheelDown
	default:
		return false
	}

	c.m.X = float32((x * c.width) / c.cols)
	c.m.Y = float32((y * c.height) / c.rows)

	c.m.Direction = 0
	switch dir {
	case 'M':
		c.m.Direction = mouse.DirPress
	case 'm':
		c.m.Direction = mouse.DirRelease
	default:
		return false
	}
	if c.m.Button < 0 {
		c.m.Direction = mouse.DirStep
	}
	return true
}
