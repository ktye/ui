package ui

/* TODO port to v2

import (
	"github.com/eaburns/T/rope"
)

type Tree struct {
	Target **Tree
	List
	Execute func(Plant) bool
	plant   []plant
}

var TreeRunes = [6]rune{'+', '-', '∘', '│', '├', '└'} // {'⊞', '⊟', '─', '│', '├', '└'}

type Plant interface {
	String() string
	Leaf() bool
	Count() int
	Child(i int) Plant
}

type plant struct {
	t Plant
	l int // level
	s string
	u int    // state: -1 closed, 1 open, 0 leaf
	e bool   // last child
	p *plant // parent
}

func (p plant) String() string {
	b := make([]rune, p.l+1)
	if n := p.l - 1; n > 0 {
		for i := 0; i < n; i++ {
			b[i] = TreeRunes[3] // │
		}
	}
	var e rune = TreeRunes[4] // ├
	if p.e {
		e = TreeRunes[5] // └
	}
	if p.l > 0 {
		b[p.l-1] = e
	}
	var u rune = TreeRunes[2] // ∘, ─
	if p.u < 0 {
		u = TreeRunes[1] // -, ⊟
	} else if p.u > 0 {
		u = TreeRunes[0] // +, ⊞
	}
	b[p.l] = u
	return string(b) + p.s + "\n"
}

func (t *Tree) SetRoot(p Plant) {
	if p.Leaf() == false {
		n := p.Count()
		t.plant = make([]plant, n)
		for i := 0; i < n; i++ {
			c := p.Child(i)
			t.plant[i].t = c
			t.plant[i].s = c.String()
			if c.Leaf() == false {
				t.plant[i].u = -1
			}
			if i == n-1 {
				t.plant[i].e = true
			}
		}
	}
	t.List.Execute = t.execute
	t.setText(0)
}

func (t *Tree) setText(current int) {
	at := int64(0)
	tb := t.List.TextBox
	if tb != nil {
		at = tb.At()
	}

	r := rope.New("")
	for _, p := range t.plant {
		r = rope.Append(r, rope.New(p.String()))
	}
	t.List.SetText(r)
	if tb != nil {
		a := t.List.addr[current]
		tb.SetDot([2]int64{a, a})
		tb.ScrollAt(at)
	}
}

func (t *Tree) execute() bool {
	d := t.TextBox.Dot()
	n := t.List.lineNum(d[0])
	if n < len(t.plant) {
		p := t.plant[n]
		if p.u == 0 && t.Execute != nil {
			return t.Execute(p.t)
		} else if p.u != 0 {
			return t.toggle(p, n)
		}
	}
	return false
}
func (t *Tree) toggle(p plant, k int) bool {
	if p.u < 0 {
		n := p.t.Count()
		plants := make([]plant, n)
		for i := 0; i < n; i++ {
			c := p.t.Child(i)
			u := plant{
				t: c,
				l: p.l + 1,
				s: c.String(),
				u: 0,
				e: i == n-1,
				p: &p,
			}
			if c.Leaf() == false {
				u.u = -1
			}
			plants[i] = u
		}
		t.plant = append(t.plant[:k+1], append(plants, t.plant[k+1:]...)...)
		t.plant[k].u = 1
	} else if p.u > 0 {
		l := p.l
		c := 0
		for i := k + 1; ; i++ {
			if i >= len(t.plant) || t.plant[i].l <= l {
				break
			}
			c++
		}
		t.plant[k].u = -1
		copy(t.plant[k+1:], t.plant[k+1+c:])
		t.plant = t.plant[:len(t.plant)-c]
	} else {
		return false
	}
	t.setText(k)
	return true
}
*/

/*
type plant struct {
	t Plant
	l int // level
	s string
	u int   // state: -1 closed, 1 open, 0 leaf
	e bool  // last child
	p *plant // parent
}
*/

/*           t l s     u    e   p
⊞Alpha       T 0 Alpha open
├⊞Alpha      T 1 Alpha open
│├─Beta      T 2 Beta leaf
│├─Beta      T 2 Beta leaf
│└⊟Werner    T 2 Werner closed
├─Beta       T 1 Beta leaf last *T
*/
