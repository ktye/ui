package main

import (
	"reflect"
	"strconv"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
)

func newwin() ui.Widget {
	c := base.Container{}
	ed := nextdev("ed")
	pl := nextdev("pl")
	b := []*base.Button{
		base.NewButton(ed, "", makedev(ed, &c, newed)),
		base.NewButton(pl, "", makedev(pl, &c, newpl)),
	}
	w := base.NewButtonBar(&base.Blank{}, b)
	w.Vertical = true
	w.Flip = true
	c.Widget = w
	return &c
}
func nextdev(d string) string { // return an unused "ed", "ed1", …
	dv := ik["devvars"].(map[string]func(v, v) v)
	if dv[d] == nil {
		return d
	}
	for i := 1; ; i++ {
		s := d + strconv.Itoa(i)
		if dv[s] == nil {
			return s
		}
	}
	return "not reached"
}
func closebutton(c *base.Container, dev string) *base.Button {
	return base.NewButton("close", "", func() int {
		d := ik["devvars"].(map[string]func(v, v) v)
		d[dev] = nil
		return closewin(c, tiler, nil, 0)
	})
}
func closewin(c *base.Container, current, parent ui.Widget, child int) int { // find c and close it
	switch t := current.(type) {
	case *base.Tiler:
		return closewin(c, t.Widget, t, 0)
	case *base.Split:
		if t.Kids[0].Widget == c {
			setchild(parent, t.Kids[1].Widget, child)
			return -1
		} else if t.Kids[1].Widget == c {
			setchild(parent, t.Kids[0].Widget, child)
			return -1
		}
		r := closewin(c, t.Kids[0].Widget, t, 0)
		if r == -1 {
			return -1
		}
		return closewin(c, t.Kids[1].Widget, t, 1)
	default:
		println("close unknown:", reflect.TypeOf(current).String())
	}
	return 0
}
func setchild(parent ui.Widget, w ui.Widget, child int) {
	switch t := parent.(type) {
	case *base.Split:
		t.Kids[child].Widget = w
	case *base.Tiler:
		t.Widget = w
	default:
		println("setchild: unknown parent type:", reflect.TypeOf(parent).String())
	}
}

type devvar interface {
	ui.Widget
	call(v, v) v
}

func makedev(dev string, c *base.Container, mk func(cb *base.Button) devvar) func() int {
	return func() int {
		dv := mk(closebutton(c, dev))
		d := ik["devvars"].(map[string]func(v, v) v)
		d[dev] = dv.call
		c.Widget = dv
		return -1
	}
}
