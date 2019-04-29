package base

import (
	"image"
	"image/draw"
)

type List struct {
	Base
	Target    **List
	Execute   func() int // double-click or Enter callback
	Delete    func() int // callback for del key
	List      []Stringer // list entries
	Colorsets []Colorset // custom colorsets, list entries may have a Color() int method
	Sel       []bool     // #List
	Single    bool       // single selection

	top, cur  int
	pagesize  int // number of visible lines
	clickLine int
	drag      bool
}

func (l *List) SetList(x []string) {
	l.List = make([]Stringer, len(x))
	for i := range x {
		l.List[i] = str(x[i])
	}
}
func (l *List) Selection() (r []int) {
	for i, s := range l.Sel {
		if s {
			r = append(r, i)
		}
	}
	return
}

func (l *List) Draw(dst *image.RGBA, force bool) {
	if l.Target != nil {
		*l.Target = l
	}
	if force == false && l.draw == false {
		return
	}
	if len(l.List) > 0 && len(l.Sel) != len(l.List) {
		l.Sel = make([]bool, len(l.List))
	}
	save := Colors
	defer func() {
		Colors = save
	}()
	clear(dst)
	p, max := image.Point{}, dst.Bounds().Dy()
	l.pagesize = max / Font.size
	l.setview() // bring current into view after resize
	for i := l.top; i < len(l.List); i++ {
		switch {
		case i == l.cur && l.Sel[i]:
			Colors = ListCurSel
		case i == l.cur:
			Colors = ListCurCol
		case l.Sel[i]:
			Colors = ListSelect
		default:
			if cr, ok := l.List[i].(colored); ok {
				if c := cr.Color() - 1; c >= 0 && c < len(l.Colorsets) {
					Colors = l.Colorsets[c]
				}
			}
		}
		if Colors[1] != save[0] {
			r := image.Rectangle{Min: dst.Rect.Min.Add(p), Max: dst.Rect.Max}
			r.Max.Y = r.Min.Y + Font.size
			draw.Draw(dst, r, Colors[1], image.ZP, draw.Src)
		}
		s := l.List[i].String()
		String(dst, p, s)
		Colors = save
		p.Y += Font.size
		if i-l.top == l.pagesize-1 {
			break
		}
	}
}
func (l *List) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	switch but {
	case 0: // move; also if a button is down
		if l.clickLine == 0 {
			return 0
		}
		if n := 1 + l.top + pos.Y/Font.size - l.clickLine; n != 0 && dir == 0 {
			l.drag = true
			if mod&1 != 0 { // but1+shift+move  toggle
				l.toggleNext(n)
				l.clickLine += n
				return l.DrawSelf()
			} else { // but1+move → scroll
				l.top -= n
				return l.setcur()
			}
		}
	case 1:
		if dir > 0 { // drag start line number + 1
			l.clickLine = 1 + l.top + pos.Y/Font.size
			if mod&1 != 0 { // but1+shift → toggle start
				if n := l.clickLine - 1; n >= 0 && n < len(l.List) {
					l.cur = n
					l.toggle()
				}
			}
		} else if dir < 0 && l.drag == false { // but1 release (non-drag) → toggle
			if n := l.top + pos.Y/Font.size; n+1 == l.clickLine {
				if n >= 0 && n < len(l.List) {
					l.cur = n
					l.toggle()
				}
			}
			l.clickLine = 0
		} else if dir < 0 && l.drag {
			l.clickLine = 0
		}
		l.drag = false
		return l.DrawSelf()
	case -1:
		l.top -= 5
		return l.setcur()
	case -2:
		l.top += 5
		return l.setcur()
	case 3:
		if dir > 0 {
			if mod&1 != 0 { // shift+right-click → unselect all
				l.Sel = make([]bool, len(l.List))
				return l.DrawSelf()
			} else if l.Menu != nil { // right-click or long-press → menu
				return l.Menu.Show()
			}
		}
	}
	return 0
}
func (l *List) Key(r rune, code uint32, dir int, mod uint32) int {
	//println("key", code, dir, r)
	if len(l.List) == 0 {
		return 0
	}
	if dir >= 0 {
		switch code {
		case 81: // ↓
			if mod&1 != 0 {
				l.toggle()
			}
			l.cur++
			return l.setview()
		case 82: // ↑
			if mod&1 != 0 {
				l.toggle()
			}
			l.cur--
			return l.setview()
		case 75: // ↑↑
			if mod&1 != 0 {
				l.toggleNext(-l.pagesize - 1)
			} else {
				l.cur -= l.pagesize - 1
			}
			return l.setview()
		case 78: // ↓↓
			if mod&1 != 0 {
				l.toggleNext(l.pagesize - 1)
			} else {
				l.cur += l.pagesize - 1
			}
			return l.setview()
		case 74: // ↑↑↑
			if mod&1 != 0 && l.cur > 0 {
				for l.cur > 0 {
					l.cur--
					l.toggle()
				}
			}
			l.cur = 0
			return l.setview()
		case 77: // ↓↓↓
			if mod&1 != 0 && l.cur < len(l.List)-1 {
				for l.cur < len(l.List)-1 {
					l.cur++
					l.toggle()
				}
			}
			l.cur = len(l.List) - 1
			return l.setview()
		case 40, 44: // space, enter
			l.toggle()
			return l.DrawSelf()
		case 4: // a
			if mod&2 != 0 { // cntrl-a
				if l.Single == false {
					for i := range l.Sel {
						l.Sel[i] = true
					}
					return l.DrawSelf()
				}
			}
		}
	}
	return 0
}
func (l *List) setview() int { // set view after changing cur
	if l.cur >= len(l.List) {
		l.cur = len(l.List) - 1
	}
	if l.cur < 0 {
		l.cur = 0
	}
	if l.cur < l.top {
		l.top = l.cur
	}
	if l.cur >= l.top+l.pagesize {
		l.top = l.cur - l.pagesize + 1
	}
	return l.DrawSelf()
}
func (l *List) setcur() int { // set cur after changing top
	if l.top >= len(l.List) {
		l.top = len(l.List) - 1
	}
	if l.top < 0 {
		l.top = 0
	}
	if l.cur < l.top {
		l.cur = l.top
	} else if l.cur >= l.top+l.pagesize {
		l.cur = l.top + l.pagesize - 1
	}
	return l.DrawSelf()
}

func (l *List) toggle() {
	if l.cur < len(l.Sel) {
		if l.Single && l.Sel[l.cur] == false {
			l.Sel = make([]bool, len(l.List))
			l.Sel[l.cur] = true
			return
		}
		l.Sel[l.cur] = !l.Sel[l.cur]
	}
}
func (l *List) toggleNext(n int) {
	dec := false
	if n < 0 {
		dec = true
		n = -n
	}
	for i := 0; i < n; i++ {
		c := l.cur + 1
		if dec {
			c = l.cur - 1
		}
		if c >= 0 && c < len(l.List) {
			l.cur = c
			l.toggle()
		}
	}
}

type colored interface{ Color() int }
type Stringer interface {
	String() string
}
type str string

func (s str) String() string { return string(s) }
