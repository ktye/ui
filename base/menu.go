package base

type Menu struct {
	*List
}

func NewMenu() *Menu {
	return nil // TODO
}
func (m Menu) Show() int {
	println("show-menu / nyi")
	return 0
}
