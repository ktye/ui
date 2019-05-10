package main

func main() {}

/* TODO port to v2

import (
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"path/filepath"

	"github.com/eaburns/T/rope"
	"github.com/ktye/iv/cmd/lui/apl385"
	"github.com/ktye/ui"
	"golang.org/x/mobile/event/key"
)

var win *ui.Window

func main() {
	var root string
	flag.StringVar(&root, "root", "", "directory for file browse mode")
	flag.Parse()

	t := &ui.Tree{}

	if root == "" {
		p := Node{
			c: []Node{
				Node{"alpha", nil},
				Node{"beta", nil},
				Node{"gamma", []Node{
					Node{"a", nil},
					Node{"b", nil},
				}},
			},
		}
		t.SetRoot(p)
	} else {
		r := fs{s: root, p: root}
		r.Count()
		t.SetRoot(&r)
		t.Execute = edit
	}

	win = ui.New(nil)
	win.SetFont(apl385.TTF(), 20)
	win.Top.W = t
	win.Render()

	for {
		select {
		case e := <-win.Inputs:
			win.Input(e)

		case err, ok := <-win.Error:
			if !ok {
				return
			}
			fmt.Println("ui:", err)
		}
	}
}

// Node implements a custom tree model.
type Node struct {
	s string
	c []Node
}

func (n Node) String() string       { return n.s }
func (n Node) Leaf() bool           { return n.c == nil }
func (n Node) Count() int           { return len(n.c) }
func (n Node) Child(i int) ui.Plant { return n.c[i] }

// fs implements a tree model of the file system.
type fs struct {
	s string // base name
	p string // full path
	c []fs   // children
}

func (f *fs) String() string { return f.s }
func (f *fs) Leaf() bool     { return f.c == nil }
func (f *fs) Count() int {
	files, err := ioutil.ReadDir(f.p)
	if err != nil {
		println(err.Error())
		return 0
	}
	f.c = make([]fs, len(files))
	for i := range files {
		fi := files[i]
		name := fi.Name()
		c := fs{
			s: name,
			p: filepath.Join(f.p, name),
		}
		if fi.IsDir() != false {
			c.c = make([]fs, 0)
		}
		f.c[i] = c
	}
	return len(f.c)
}
func (f *fs) Child(i int) ui.Plant {
	return &f.c[i]
}
func edit(p ui.Plant) bool {
	if f, ok := p.(*fs); ok {
		b, err := ioutil.ReadFile(f.p)
		if err != nil {
			println(err.Error())
		} else {
			ed := editor{Edit: ui.Edit{}}
			ed.SetText(rope.New(string(b)))
			ed.save = win.Top.W
			top(&ed)
		}
	}
	return false
}

type editor struct {
	ui.Edit
	save ui.Widget
}

func (e *editor) Key(w *ui.Window, self *ui.Kid, k key.Event, m ui.Mouse, orig image.Point) (res ui.Result) {
	if k.Code == key.CodeEscape && k.Direction == key.DirRelease {
		e.quit()
	}
	return e.Edit.Key(w, self, k, m, orig)
}
func (e *editor) quit() { top(e.save) }

func top(w ui.Widget) {
	win.Top.W = w
	win.Top.Layout = ui.Dirty
	win.Top.Draw = ui.Dirty
	win.Render()
}
*/
