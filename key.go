package ui

import (
	"image"

	"golang.org/x/mobile/event/key"
)

// Key delivers a key press event to the widget tree.
// Key is typically called by Input.
func (w *Window) Key(k key.Event) {
	if kb := w.keyboard; kb != nil {
		kb.TranslateKey(&k)
	}

	r := w.Top.W.Key(w, &w.Top, k, w.mouse, image.ZP)
	if !r.Consumed {
		switch k.Rune {
		case '\t':
			first := w.Top.W.FirstFocus(w, &w.Top)
			if first != nil {
				r.Warp = first
				r.Consumed = true
			}
		}
	}
	w.apply(r)
}

// KeyTranslator can modify a key event.
// It is used to implement different keyboard layouts.
type KeyTranslator interface {
	TranslateKey(*key.Event)
}

func (w *Window) SetKeyTranslator(kb KeyTranslator) {
	w.keyboard = kb
}

// AplKeyboard is a KeyTranslator which changes the Rune of a key event,
// if Alt-Gr[+Shift] is pressed.
// It uses AplKeymap.
type AplKeyboard struct{}

func (a AplKeyboard) TranslateKey(e *key.Event) {
	if e.Modifiers == 6 || e.Modifiers == 7 {
		if r, ok := AplKeymap[e.Code]; ok {
			if e.Modifiers == 6 {
				e.Rune = r[0]
			} else {
				e.Rune = r[1]
			}
		}
	}
}

// AplKeymap maps from a key code to two runes.
// It can be overwritten.
// The first one is used if the Alt-Gr key is used,
// The sencond if both, Alt-Gr and Shift are used.
// See github.cmd/ktye/iv/cmd/lui/keyboard.go for the keyboard layout.
var AplKeymap = map[key.Code][2]rune{
	// Top row.
	53: {'⋄', '⍨'},
	30: {'¨', '¡'},
	31: {'¯', '€'},
	32: {'<', '£'},
	33: {'≤', '⍧'},
	34: {'=', '≢'},
	35: {'≥', 'τ'},
	36: {'>', 'η'},
	37: {'≠', '⍂'},
	38: {'∨', '⍱'},
	39: {'∧', '⍲'},
	45: {'×', '≡'},
	46: {'÷', '⌹'},

	// Second row.
	20: {'?', '¿'},
	26: {'⍵', '⌽'},
	8:  {'∊', '⍷'},
	21: {'⍴', 'λ'},
	23: {'∼', '⍉'},
	28: {'↑', '¥'},
	24: {'↓', 'μ'},
	12: {'⍳', '⍸'},
	18: {'○', '⍥'},
	19: {'⋆', '⍟'},
	47: {'←', 'π'},
	48: {'→', 'Ω'},
	49: {'⍝', '⍀'},

	// Third row.
	4:  {'⍺', '⊖'},
	22: {'⌈', 'σ'},
	7:  {'⌊', 'δ'},
	9:  {'_', '⍫'},
	10: {'∇', '⍒'},
	11: {'∆', '⍋'},
	13: {'∘', '⍤'},
	14: {'\'', '⌺'},
	15: {'⎕', '⍞'},
	51: {'⊢', -1},
	52: {'⊣', -1},

	// Fourth row.
	29: {'⊂', 'ζ'},
	27: {'⊃', 'ξ'},
	6:  {'∩', '⍝'},
	25: {'∪', 'ø'},
	5:  {'⊥', '⍎'},
	17: {'⊤', '⍕'},
	16: {'|', '⌶'},
	54: {'⌷', '⍪'},
	55: {'⍎', '⍙'},
	56: {'⍕', '⌿'},
}
