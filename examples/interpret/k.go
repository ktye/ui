package main

import (
	"runtime/debug"
	"strings"

	"github.com/ktye/i"
)

type s = string
type v = interface{}
type l = []v

var fmt func(v) v
var jon func(v, v) v
var num func(v) v

var ipr interp
var cnt func(v) v
var atx func(v, v) v
var lnx func(v) int
var til func(v) v
var cst func(v, v) v

type k map[v]v

func (a k) Eval(s string) string {
	s = plumb(s)
	x := run(s)
	if x != nil {
		s, o := x.(string)
		if !o {
			/*
				if p, o := isplot(x); o {
					setplot(p)
					return ""
				} else {
			*/
			s = fmt(x).(string)
			//}
		}
		return s
	}
	return ""
}

func run(t string) (r v) {
	defer func() {
		if c := recover(); c != nil {
			r = stack(c)
		}
	}()
	pr := i.P(t)
	return i.E(pr, ik)
}
func keval(l l) (r v) {
	defer func() {
		if c := recover(); c != nil {
			r = stack(c)
			println(r.(string))
		}
	}()
	return i.E(l, ik)
}

func stack(c interface{}) string {
	rs := ""
	for _, s := range strings.Split(string(debug.Stack()), "\n") {
		if strings.HasPrefix(s, "\t") {
			rs += "\n" + s[1:]
		}
	}
	if s, o := c.(string); o {
		rs += "\n" + s
	} else if e, o := c.(error); o {
		rs += "\n" + e.Error()
	}
	return rs
}

func kinit() k {
	a := make(map[v]v)
	i.E(l{}, a)
	fmt = a["$:"].(func(x v) v)
	jon = a["jon"].(func(x, y v) v)
	num = a["num"].(func(x v) v)
	cnt = a["#:"].(func(v) v)
	atx = a["@@"].(func(v, v) v)
	lnx = a["ln"].(func(v) int)
	til = a["!:"].(func(v) v)
	cst = a["$$"].(func(v, v) v)
	a["devvars"] = make(map[string]func(v, v) v)
	return k(a)
}
