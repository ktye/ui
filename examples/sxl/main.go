package main

import (
	"github.com/ktye/ui/display"
	"golang.org/x/exp/shiny/screen"
)

func main() {
	opt := screen.NewWindowOptions{
		Title: "Theodor Sexl",
	}
	d := display.New6(&opt)
	<-d.Size
}
