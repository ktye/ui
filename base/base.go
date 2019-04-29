package base

type Base struct {
	Menu *Menu
	draw bool
}

func (l *List) DrawSelf() int {
	l.draw = true
	return 1
}
