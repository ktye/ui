package property

import (
	"fmt"
	"image"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ktye/ui"
	"github.com/ktye/ui/base"
	"github.com/ktye/ui/editor"
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
	ToTable func() int
	*base.List
	initerr *base.Label
	title   *base.Label
	link    *base.Label
	editor  ui.Widget
	list
	current     int
	name        string
	add         bool
	copyIndexes []int
	keysize     int // max rune count for display name
}

func (l *List) init() {
	var err error
	l.list, err = data{l.Data}.list(l.name, l.copyIndexes, l.add)
	if err != nil {
		l.initerr = base.NewLabel(err.Error())
		return
	}
	l.link = nil
	label := l.name
	if n := len(l.copyIndexes); n == 1 {
		if l.list.IsStruct == false {
			label += " " + strconv.Itoa(l.copyIndexes[0]+1)
			if l.ToTable != nil {
				l.link = base.NewLink("(show all)", l.ToTable)
			}
		}
		if s := l.list.NameField(); s != "" {
			label += ": " + s
		}
	} else if n > 0 {
		label += " (multiple)"
	}
	l.title = base.NewLabel(label)
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
			base.Colorset{base.Red.Uniform(), base.Red.Uniform()},         // password
			base.Colorset{base.Red.Uniform(), base.White.Uniform()},       // error
		},
		Execute: l.edit,
		Delete:  l.clear,
	}
	l.List.SelectSingle(l.current)
}

func (l *List) selection() int { // single selection index
	sel := l.List.Selection()
	if len(sel) < 1 {
		return -1
	}
	return sel[0]
}
func (l *List) clear() int {
	i := l.selection()
	if i < 0 {
		return 0
	}
	p := l.list.Fields[i]
	if p.IsSlice {
		l.list.Fields[i].Values = nil
	}
	d := data{l.Data}
	if err := d.update(l.list); err != nil {
		base.Errlog(err)
	}
	l.List.List = nil
	return 1
}
func (l *List) edit() int {
	// TODO: toggle bool, show menu for popups, edit strings, …
	idx := l.selection()
	if idx < 0 {
		return 0
	}
	l.current = idx
	p := l.list.Fields[idx]
	commit := func() error {
		p.IsUpdated = true
		l.list.Fields[idx] = p
		d := data{l.Data}
		if err := d.update(l.list); err != nil {
			return err
		}
		l.List = nil
		return nil
	}

	t := p.Type
	sl := p.IsSlice
	switch {
	case !sl && len(p.Options) > 0: // combobox
		println("TODO show combobox")
	case !sl && t == reflect.TypeOf(false): // toggle bool
		if len(p.Values) == 1 && p.Values[0] == "true" {
			p.Values[0] = "false"
		} else {
			p.Values[0] = "true"
		}
		if err := commit(); err != nil {
			base.Errlog(err)
		}
		return 1
	case !sl: // line edit
		l.editor = editor.NewInputWithButtons(p.DisplayName, p.Values[0])
		in := l.editor.(*base.ButtonBar).Kid.Widget.(*editor.Input)
		in.Cancel = func() int { l.editor = nil; return 1 }
		in.Commit = func() error {
			p.Values[0] = in.Text
			err := commit()
			if err == nil {
				l.editor = nil
			}
			return err
		}
		return 1
	default:
		println(fmt.Errorf("unknown type: %s", p.Type))
		return 0
	}
	return 1
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
		return 1 // grey (disabled)
	} else if e.p.IsPassword {
		return 2 // white on white
	}
	return 0
}
func (l *List) Draw(dst *image.RGBA, force bool) {
	if l.initerr == nil && l.List == nil {
		l.init()
	}
	if l.initerr != nil {
		l.initerr.Draw(dst, force)
		return
	}
	if l.title != nil {
		l.title.Draw(dst, force)
		if l.link != nil {
			u := l.title.Used()
			r := dst.Rect
			r.Min.X = u.Max.X + 3*base.Font.Size()
			l.link.Draw(dst.SubImage(r).(*image.RGBA), force)
		}
	}
	r := dst.Rect
	r.Min.Y += (3 * base.Font.Size()) / 2
	dst = dst.SubImage(r).(*image.RGBA)
	if l.editor != nil {
		l.editor.Draw(dst, force)
		return
	}
	l.List.Draw(dst, force)
}
func (l *List) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if l.initerr != nil {
		return 0
	}
	if l.link != nil && but == 1 && pos.In(l.link.Used()) {
		return l.link.Mouse(pos, but, dir, mod)
	}
	if l.editor != nil {
		return l.editor.Mouse(pos, but, dir, mod)
	}
	return l.List.Mouse(pos, but, dir, mod)
}
func (l *List) Key(r rune, code uint32, dir int, mod uint32) int {
	if l.initerr != nil {
		return 0
	}
	if l.editor != nil {
		return l.editor.Key(r, code, dir, mod)
	}
	return l.List.Key(r, code, dir, mod)
}
