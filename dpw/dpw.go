// +build !js

// dpw is a html-5 canvas backend for ui compiled to webassembly
package dpw

func New() {
	panic("this package only runs on js/wasm")
}
