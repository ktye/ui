package ui

import (
	"image"
	"image/color"

	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/size"
)

const (
	BorderSize    = 1
	ScrollbarSize = 10
)

// State represents the layout/draw state of the widget of a Kid.
type State byte

const (
	Dirty    = State(iota) // UI itself needs layout/draw;  kids will also get a UI.Layout/UI.Draw call, with force set.
	DirtyKid               // UI itself does not need layout/draw, but one of its children does, so pass the call on.
	Clean                  // UI and its children do not need layout/draw.
	// note: order is important, Dirty is the default, Clean is highest and means least amount of work
)

// Input is an input event that is typically passed into UI through Input().
type Input struct {
	Type  InputType
	Mouse mouse.Event
	Key   key.Event
	Size  size.Event
	Func  func()
	Error error
}

// InputType presents the type of an input event.
type InputType byte

const (
	InputMouse InputType = iota // Mouse movement and/or button changes.
	InputKey                    // Key typed.
	InputFunc                   // Call the function.
	InputSize                   // window was resized, reattach; does not have/need a field in Input.
	InputError                  // An error occurred that may be recovered from.
)

// Halign represents horizontal align of elements in a Grid.
type Halign byte

const (
	HalignLeft Halign = iota // Align to the left by default, for example in a grid.
	HalignMiddle
	HalignRight
)

// Valign represents vertical align of elements in a Grid, or in a Box.
type Valign byte

const (
	ValignMiddle Valign = iota // Align vertically in the middle by default, for example in a box (line) or grid.
	ValignTop
	ValignBottom
)

// Event is returned by handlers, such as click or key handlers.
type Event struct {
	Consumed   bool // Whether event was consumed, and should not be further handled by upper UI's.  Container UIs can handle some mouse/key events and decide whether they want to pass them on, or first pass them on and only consume them when a child UI hasn't done so yet.
	NeedLayout bool // Whether UI now needs a layout. Only the UI generating the event will be marked. If you another UI needs to be marked, call MarkLayout.
	NeedDraw   bool // Like NeedLayout, but for draw.
}

// Result holds the effects of a mouse/key event, as implement by UIs.
type Result struct {
	Hit      Widget       // the UI where the event ended up
	Consumed bool         // whether event was consumed, and should not be further handled by upper UI's
	Warp     *image.Point // if set, mouse will warp to location
}

// Colors represents the style in one state of the UI.
type Colors struct {
	Text       color.Color
	Background image.Image
	Border     color.Color
}

// Colorset is typically used to style buttons. Duit provides some builtin colorsets like Primary, Danger, Success.
type Colorset struct {
	Normal, Hover Colors
}
