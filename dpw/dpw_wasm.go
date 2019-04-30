// dpw is a html-5 canvas backend for ui compiled to webassembly
//
// See ui/examples/wasm.
package dpw

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"sync"
	"syscall/js"
)

// New creates a new display server.
func New(canvasID string) *Display {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", canvasID)
	width, height := canvas.Get("width").Int(), canvas.Get("height").Int()
	ctx := canvas.Call("getContext", "2d")

	d := &Display{
		canvas: js.Global().Get("canvas"),
		ctx:    ctx,
		mouse:  make(chan mt),
		key:    make(chan kt),
		image:  make(chan *image.RGBA),
		err:    make(chan error),
	}

	mouse := func(e js.Value) mt {
		var mod uint32
		for i, s := range []string{"shiftKey", "ctrlKey", "altKey", "metaKey"} {
			println(s)
			if e.Get(s).Bool() {
				mod |= 1 << uint(i)
			}
		}
		return mt{
			But: e.Get("button").Int(),
			Pos: image.Point{e.Get("clientX").Int(), e.Get("clientY").Int()},
			Mod: mod,
		}
	}
	mouseMove := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m := mouse(args[0])
		d.mouse <- m
		return nil
	})
	mouseDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m := mouse(args[0])
		m.Dir = 1
		d.mouse <- m
		return nil
	})
	mouseUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m := mouse(args[0])
		m.Dir = -1
		d.mouse <- m
		return nil
	})
	mouseWheel := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		but := -1
		if e.Get("deltaY").Int() > 0 {
			but = -2
		}
		d.mouse <- mt{
			But: but,
			Dir: 3,
			Pos: image.Point{e.Get("clientX").Int(), e.Get("clientY").Int()},
		}
		return nil
	})
	doc.Call("addEventListener", "mousemove", mouseMove)
	doc.Call("addEventListener", "mousedown", mouseDown)
	doc.Call("addEventListener", "mouseup", mouseUp)
	doc.Call("addEventListener", "wheel", mouseWheel)

	key := func(e js.Value) kt {
		var mod uint32
		for i, s := range []string{"Shift", "Control", "Alt", "Meta"} {
			if e.Call("getModifierState", s).Bool() {
				mod |= 1 << uint(i)
			}
		}
		return kt{
			// Rune: TODO
			Code: uint32(e.Get("keyCode").Int()), // TODO convert from js code to go/mobile/key codes
			Mod:  mod,
		}
	}
	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		dir := 1
		if e.Get("repeat").Bool() {
			dir = 0
		}
		k := key(e)
		k.Dir = dir
		d.key <- k
		return nil
	})
	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		k := key(args[0])
		k.Dir = -1
		d.key <- k
		return nil
	})
	doc.Call("addEventListener", "keydown", keyDown)
	doc.Call("addEventListener", "keyup", keyUp)

	done := make(chan bool)
	go func() {
		d.img = image.NewRGBA(image.Rect(0, 0, width, height))
		println("size", width, height)
		d.image <- d.img
		<-done
	}()
	return d
}

type Display struct {
	sync.Mutex
	canvas js.Value
	ctx    js.Value
	mouse  chan mt
	key    chan kt
	image  chan *image.RGBA
	err    chan error
	img    *image.RGBA
	buf    bytes.Buffer
}

func (d *Display) Image() chan *image.RGBA { return d.image }
func (d *Display) Mouse() chan mt          { return d.mouse }
func (d *Display) Key() chan kt            { return d.key }
func (d *Display) Err() chan error         { return d.err }
func (d *Display) Flush() {
	// TODO: use the pixel buffer of a 2d drawing context directly (putImageData).
	png.Encode(&d.buf, d.img)
	s := base64.StdEncoding.EncodeToString(d.buf.Bytes())

	image := js.Global().Call("eval", "new Image()")
	image.Call("addEventListener", "load", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		d.ctx.Call("drawImage", image, 0, 0)
		return nil
	}))
	image.Set("src", pngPrefix+s)
}

const pngPrefix = "data:image/png;base64,"

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
