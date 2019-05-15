package property

import (
	"image"
	"strings"
	"unicode/utf8"

	"github.com/ktye/ui/base"
)

// NewList returns a list widget from the data source for the given field name and index (if it is a slice).
func NewList(data Source, name string, index int) *List {
	return &List{Data: data, name: name, copyIndexes: []int{index}}
}

// AddList returns a new list widget from the data source for the given field name.
// Values are copied from the given slice index. On Update, the list is appended to the data field.
func AddList(data Source, name string, index int) *List {
	return &List{Data: data, name: name, copyIndexes: []int{index}, add: true}
}

// MultiList returns a list widget that updates muliple entries of the data field slice value.
func MultiList(data Source, name string, index []int) *List {
	return &List{Data: data, name: name, copyIndexes: index}
}

// List is a widget for a property list.
type List struct {
	Data    Source // Data source must be present before first draw.
	Trigger func() // executed after update to modify other ui widgets
	*base.List
	err *base.Label
	list
	name        string
	add         bool
	copyIndexes []int
	keysize     int // max rune count for display name
}

func (l *List) init() {
	var err error
	l.list, err = data{l.Data}.list(l.name, l.copyIndexes, l.add)
	if err != nil {
		l.err = base.NewLabel(err.Error())
		return
	}
	var entries []base.Stringer
	for i, p := range l.list.Fields {
		if p.IsIgnored || p.IsHidden {
			continue
		}
		n := i
		prop := p
		entries = append(entries, entry{
			i: n,
			p: prop,
			k: prop.DisplayName,
			v: strings.Join(prop.nablas(), "∘"),
			l: l,
		})
	}
	max := 0
	for i := range entries {
		e := entries[i].(entry)
		if n := utf8.RuneCountInString(e.k); n > max {
			max = n
		}
	}
	for i := range entries {
		e := entries[i].(entry)
		if n := utf8.RuneCountInString(e.k); n < max {
			e.k = e.k + strings.Repeat(" ", max-n)
			entries[i] = e
		}
	}

	l.List = &base.List{
		List:   entries,
		Single: true,
		Colorsets: []base.Colorset{
			base.Colorset{base.LightGrey.Uniform(), base.White.Uniform()}, // disabled
			base.Colorset{base.Black.Uniform(), base.White.Uniform()},     // normal
			base.Colorset{base.Red.Uniform(), base.Red.Uniform()},         // password
		},
		// TODO: Execute: enable disabled, fields, edit values?
		// TODO: Delete: clear slices, ZeroString for values?
	}
}

func (p property) nablas() []string {
	if p.Options == nil {
		return p.Values
	}
	s := make([]string, len(p.Values))
	for i, t := range p.Values {
		s[i] = t + "∇"
	}
	return s
}

// base.List entry (a subset of the property list entries)
type entry struct {
	i    int // index into property list
	p    property
	k, v string
	l    *List
}

func (e entry) String() string {
	return e.k + "│" + e.v // TODO: add …+- maybe dynamic depending on dst width
}
func (e entry) Color() int {
	if !e.p.IsUnique {
		return 0 // grey (disabled)
	} else if e.p.IsPassword {
		return 2 // white on white
	}
	return 1
}

func (l *List) Draw(dst *image.RGBA, force bool) {
	if l.err != nil {
		l.err.Draw(dst, force)
		return
	}
	if l.List == nil {
		l.init()
	}
	l.List.Draw(dst, force)
}
func (l *List) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if l.err != nil {
		return 0
	}
	return l.List.Mouse(pos, but, dir, mod)
}
func (l *List) Key(r rune, code uint32, dir int, mod uint32) int {
	if l.err != nil {
		return 0
	}
	return l.List.Key(r, code, dir, mod)
}
