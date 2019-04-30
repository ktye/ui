// +build !js

package base

import "github.com/eaburns/T/clipboard"

var Clipboard clipboard.Clipboard

func init() {
	Clipboard = clipboard.New()
}
