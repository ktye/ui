package main

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/eaburns/T/rope"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/editor"
)

// ed(itor) device variable
//
// set
//	ed:"text"               /set text
//	ed:("line1";"line2"…)   /set lines
//	ed:5                    /mark line 5
//	ed:4 6                  /mark byte range
// get
//	ed → dict:
//	text|content
//	sel |selected text
//	dot |(start,end) selection byte index
type ed struct {
	*editor.Edit
}

func (e ed) call(x, idx v) v {
	if x == nil {
		tb := e.Edit.TextBox
		dot := tb.Dot()
		return [2]l{l{"text", "sel", "dot"}, l{tb.Text().String(), tb.Selection(), dot[:]}}
	}
	switch t := x.(type) {
	case string:
		e.Edit.SetText(rope.New(t))
	case []string:
		e.Edit.SetText(rope.New(strings.Join(t, "\n") + "\n"))
	case int:
		e.Edit.MarkAddr(strconv.Itoa(t))
	case complex128:
		e.Edit.MarkAddr(strconv.Itoa(int(real(t))))
	case []complex128:
		if len(t) == 2 {
			e.Edit.TextBox.SetDot([2]int64{int64(real(t[0])), int64(real(t[1]))})
		}
	default:
		println(reflect.TypeOf(x).String())
	}
	return nil
}

func newed(cb *base.Button) devvar {
	e := ed{Edit: editor.New("")}
	e.Edit.Nowrap = true
	e.Edit.Menu = e.Edit.StandardMenu()
	evlbutton(e.Edit)
	e.Edit.Menu.Buttons = append(e.Edit.Menu.Buttons, cb)
	return e
}

func evlbutton(e *editor.Edit) { // add a run menu entry to the editor
	b := base.NewButton("eval", "", func() int {
		rpl.Execute(nil, e.Selection())
		return -1
	})
	e.Menu.Buttons = append([]*base.Button{b}, e.Menu.Buttons...)
}

func edit(t s, line int) {
	ed := nextdev("ed")
	c := &base.Container{}
	makedev(ed, c, newed)()
	s := base.NewSplit(tiler.Widget, c)
	s.Vertical = true
	tiler.Widget = s
	keval(l{":", ed, l{"`", t}}) // ed:s…
	if line > 0 {
		go func() { // call ed:line after the textbox is created during draw.
			win.Call <- func() int {
				keval(l{":", ed, line}) // TODO?
				return -1
			}
		}()
	}
}

/*
func edit(t s, line int) {
	edt.SetText(rope.New(t))
	if sp1.Ratio > 0.95 {
		sp1.Ratio = 0
	}
	if line > 0 {
		edt.MarkAddr(strconv.Itoa(line))
	}
}

*/
