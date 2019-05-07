package editor

import (
	"fmt"

	"github.com/eaburns/T/rope"
	"github.com/ktye/ui/base"
)

// StandardMenu adds a context menu to the editor with standard entries.
func (e *Edit) StandardMenu() *base.Menu {
	entries := []struct {
		s string
		f func() int
	}{
		{"find", func() int { e.Find(); return 1 }},
		{"(un)wrap", func() int { e.Nowrap = !e.Nowrap; e.TextBox.SetNowrap(e.Nowrap); return 1 }},
		{"clear", func() int { e.SetText(rope.New("")); return 1 }},
		{"add numbers", e.addNumbers},
	}
	buttons := make([]*base.Button, len(entries))
	for i, e := range entries {
		buttons[i] = base.NewButton(e.s, "", e.f)
	}
	return base.NewMenu(buttons)
}

func (e *Edit) addNumbers() int { // this adds numbers to the edit text. use only for display.
	src := e.TextBox.Text()
	n := 1
	dst := rope.New("")
	line := rope.New("")
	last := false
	for last == false {
		idx := rope.IndexRune(src, '\n')
		if idx < 0 {
			last, idx = true, src.Len()
			line = src
		} else {
			line, src = rope.Split(src, idx+1)
		}
		dst = rope.Append(rope.Append(dst, rope.New(fmt.Sprintf("%2d|", n))), line)
		n++
	}
	e.TextBox.SetText(dst)
	return 1
}
