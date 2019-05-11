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

func run(t string, a map[v]v) (r interface{}) {
	defer func() {
		if c := recover(); c != nil {
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
			r = rs
		}
	}()
	pr := i.P(t)
	return i.E(pr, a)
}

func kinit() map[v]v {
	a := make(map[v]v)
	i.E(l{}, a)
	fmt = a["$:"].(func(x v) v)
	jon = a["jon"].(func(x, y v) v)
	num = a["num"].(func(x v) v)

	regplot(a)
	return a
}
