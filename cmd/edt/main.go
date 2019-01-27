// edt is a line editor
//
// Usage:
//	edt  # interactive, start with empty buffer
//	edt FILE # interactive, read text from FILE
//
// TODO:
// read commands from file, operate on multiple files
//	edt -e CMD FILE
//	edt [-e CMD] FILES...
package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/eaburns/T/edit"
	"github.com/eaburns/T/rope"
)

func main() {
	var text rope.Rope
	if n := len(os.Args); n < 2 {
		text = rope.New("")
	} else if n > 2 {
		fmt.Fprintf(os.Stderr, "too many arguments")
		os.Exit(1)
	} else {
		text = readfile(os.Args[1])
	}

	var dot [2]int64
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cmd := scanner.Text()
		diffs, err := edit.Edit(dot, cmd, os.Stdout, text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "? %s\n", err)
		}
		text, _ = diffs.Apply(text)
		dot = diffs.Update(dot)
	}
}

func readfile(name string) rope.Rope {
	f, err := os.Open(os.Args[1])
	fatal(err)
	defer f.Close()
	t, err := rope.ReadFrom(bufio.NewReader(f))
	fatal(err)
	return t
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
