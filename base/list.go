package base

import (
	"image"
	"image/draw"
)

type List struct {
	Target    **List
	Menu      *Menu
	Execute   func() int // double-click or Enter callback
	Delete    func() int // callback for del key
	List      []Stringer // list entries
	Colorsets []Colorset // custom colorsets
	Colors    []int      // index into Colorsets (#List), based-1
	Sel       []bool     // #List
	Single    bool       // single selection

	top, cur int
	draw     bool
	pagesize int // number of visible lines
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
			if i < len(l.Colors) && l.Colors[i] != 0 {
				Colors = l.Colorsets[l.Colors[i]-1]
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
	if but != 0 {
		println("mouse", but, dir)
	}
	if but == -1 {
		l.top -= 5
		return l.setcur()
	} else if but == -2 {
		l.top += 5
		return l.setcur()
	} else if but == 3 && dir == 2 { // right-click or long-press: menu
		if l.Menu != nil {
			return l.Menu.Show()
		}
	}
	return 0
}
func (l *List) Key(r rune, code uint32, dir int, mod uint32) int {
	println("key", code, dir, r)
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
			l.cur -= l.pagesize - 1
			return l.setview()
		case 78: // ↓↓
			l.cur += l.pagesize - 1
			return l.setview()
		case 74: // ↑↑↑
			l.cur = 0
			return l.setview()
		case 77: // ↓↓↓
			l.cur = len(l.List) - 1
			return l.setview()
		case 40, 44: // space, enter
			l.toggle()
			return l.drawself()
		case 4: // a
			if mod&2 != 0 { // cntrl-a
				if l.Single == false {
					for i := range l.Sel {
						l.Sel[i] = true
					}
					return l.drawself()
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
	return l.drawself()
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
	return l.drawself()
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
func (l *List) drawself() int {
	l.draw = true
	return 1
}

type Stringer interface {
	String() string
}
type str string

func (s str) String() string { return string(s) }
