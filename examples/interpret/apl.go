package main

import (
	"strings"

	ivapl "github.com/ktye/iv/apl"
	"github.com/ktye/iv/apl/numbers"
	"github.com/ktye/iv/apl/operators"
	"github.com/ktye/iv/apl/primitives"
)

func ivinit() apl {
	a := ivapl.New(nil)
	numbers.Register(a)
	primitives.Register(a)
	operators.Register(a)
	return apl{a}
}

type apl struct{ *ivapl.Apl }

func (a apl) Eval(s string) string {
	var b strings.Builder
	a.Apl.SetOutput(&b)
	err := a.Apl.ParseAndEval(s)
	if err != nil {
		return err.Error()
	}
	return b.String()
}
