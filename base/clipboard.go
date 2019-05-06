// +build !js

package base

import (
	"sync"

	"github.com/atotto/clipboard"
)

// Clipboard is the program wide clipboard.
// It is already initialized and can be used with Store and Fetch.
var Clipboard cb

func init() {
	Clipboard = newClipboard()
}

// adapted from github.com/eaburns/T/clipboard
// It uses github.com/atotto/clipboard (the system clipboard if available), or an in-memory clipboard.

type cb interface {
	Store(string) error
	Fetch() (string, error)
}

func newClipboard() cb {
	if clipboard.Unsupported {
		return newMem()
	}
	return sysClipboard{}
}
func newMem() cb {
	return &memClipboard{text: ""}
}

type sysClipboard struct{}

func (sysClipboard) Store(s string) error {
	return clipboard.WriteAll(s)
}
func (sysClipboard) Fetch() (string, error) {
	return clipboard.ReadAll()
}

type memClipboard struct {
	text string
	sync.Mutex
}

func (m *memClipboard) Store(s string) error {
	m.Lock()
	m.text = s
	m.Unlock()
	return nil
}
func (m *memClipboard) Fetch() (string, error) {
	m.Lock()
	s := m.text
	m.Unlock()
	return s, nil
}
