package ui

import "image"

// Kid holds a UI and its layout/draw state.
type Kid struct {
	W      Widget          // UI this state is about.
	R      image.Rectangle // Location and size within this UI.
	Draw   State           // Whether UI or its children need a draw.
	Layout State           // Whether UI or its children need a layout.
	ID     string          // For (re)storing settings with ReadSettings and WriteSettings. If empty, no settings for the UI will be (re)stored.
}
