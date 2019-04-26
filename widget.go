package ui

// Widget is a part of the user interface.
// It can be implemented in an external package without depending on package ui.
//
// Draw asks the widget to draw isself on img, which is a subimage of the window.
// The widget should only redraw, if force is true, or if it requested the drawing from Mouse or Key.
//
// Mouse is sent for a mouse event.
// The return value is 0 if it is ignored, >0 if the widget wants to redraw and <0 if the complete
// ui should be redrawn.
//
// Key is send for a key press or release.
