package property

import (
	"bytes"
	"fmt"
	"image"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/ktye/ui/base"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

// NewTable returns a table widget for the for the data source and the given field name.
// The field value must be a slice.
func NewTable(data Source, name string) *Table {
	return &Table{Data: data, name: name}
}

// List is a widget for a property table that shows all property elements from a slice value.
// Executing a line returns a property List for the row.
// Selecting multiple lines returns a MultiList.
type Table struct {
	Data Source
	*base.List
	*base.ButtonBar
	Buttons []*base.Button
	Edit    func([]int) int
	err     *base.Label
	table
	name string
	// TODO
}

func (t *Table) init() {
	var err error
	t.table, err = data{t.Data}.table(t.name)
	if err != nil {
		t.err = base.NewLabel(err.Error())
		return
	}
	var buf bytes.Buffer
	tw := tabwriter.NewWriter(&buf, 0, 8, 1, ' ', tabwriter.Debug)
	for i, row := range t.table {
		if i == 0 { // title
			sep := ""
			for _, p := range row.Fields {
				if p.IsIgnored || p.IsHidden {
					continue
				}
				fmt.Fprintf(tw, "%s%s", sep, p.DisplayName)
				sep = "\t"
			}
			fmt.Fprintf(tw, "\n")
		}
		sep := ""
		for _, p := range row.Fields {
			if p.IsIgnored || p.IsHidden {
				continue
			}
			s := strings.Join(p.Values, "∘")
			if len(s) > 20 {
				s = s[:19] + "…"
			}
			fmt.Fprintf(tw, "%s%s", sep, s)
			sep = "\t"
		}
		fmt.Fprintf(tw, "\n")
	}
	tw.Flush()
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 { // empty last line (final newline)
		lines = lines[:len(lines)-1]
	}
	entries := make([]base.Stringer, len(lines))
	for i := range entries {
		entries[i] = tabline{i - 1, string(lines[i])}
	}
	t.List = &base.List{
		List:   entries,
		Single: false,
		Colorsets: []base.Colorset{
			base.Colorset{base.White.Uniform(), base.Black.Uniform()}, // header inverted
		},
		Delete:  t.delete,
		Execute: t.edit,
	}
	if t.Buttons != nil {
		t.ButtonBar = base.NewButtonBar(t.List, t.Buttons)
	}
}

type tabline struct {
	i int // -1: header
	s string
}

func (l tabline) String() string { return l.s }
func (l tabline) Color() int {
	if l.i < 0 {
		return 1
	}
	return 0
}

func (t *Table) Draw(dst *image.RGBA, force bool) {
	if t.err == nil && t.List == nil {
		t.init()
	}
	if t.err != nil {
		t.err.Draw(dst, force)
		return
	}
	if t.ButtonBar != nil {
		t.ButtonBar.Draw(dst, force)
		return
	}
	t.List.Draw(dst, force)
}
func (t *Table) Mouse(pos image.Point, but int, dir int, mod uint32) int {
	if t.err != nil {
		return 0
	}
	if t.ButtonBar != nil {
		return t.ButtonBar.Mouse(pos, but, dir, mod)
	}
	return t.List.Mouse(pos, but, dir, mod)
}
func (t *Table) Key(r rune, code uint32, dir int, mod uint32) int {
	if t.err != nil {
		return 0
	}
	if t.ButtonBar != nil {
		return t.ButtonBar.Key(r, code, dir, mod)
	}
	return t.List.Key(r, code, dir, mod)
}

func (t *Table) DefaultButtons() {
	t.Buttons = []*base.Button{
		base.FillButton(),
		&base.Button{Text: "Edit", Icon: "edit", Execute: t.edit},
		&base.Button{Text: "Add", Icon: "add", Execute: t.additem},
		&base.Button{Text: "Delete", Icon: "del", Execute: t.delete},
		&base.Button{Text: "Move Up", Icon: "up", Execute: func() int { t.moveitem(true); return 1 }},
		&base.Button{Text: "Move Down", Icon: "down", Execute: func() int { t.moveitem(false); return 1 }},
		base.FillButton(),
	}
}
func (t *Table) edit() int {
	sel := t.List.Selection()
	if len(sel) > 0 && sel[0] == 0 {
		sel = sel[1:]
	}
	for i := range sel {
		sel[i]--
	}
	if len(sel) == 0 {
		return 0
	}
	if t.Edit != nil {
		return t.Edit(sel)
	}
	return 0
}
func (t *Table) additem() int {
	d := data{t.Data}
	var idx []int
	if k := t.singleSelection(); k >= 0 {
		idx = []int{k}
	}
	ps, err := d.list(t.name, idx, true)
	if err != nil {
		t.err = base.NewLabel(err.Error())
		return 1
	}
	if err := d.update(ps); err != nil {
		t.err = base.NewLabel(err.Error())
	}
	t.List = nil // rebuild on next draw
	return -1
}
func (t *Table) delete() int {
	d := data{t.Data}
	sel := t.List.Selection()
	if len(sel) > 0 && sel[0] == 0 {
		sel = sel[1:]
	}
	if len(sel) == 0 {
		return 0
	}
	for i := range sel {
		sel[i] = sel[i] - 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sel)))

	for _, i := range sel {
		if err := d.deleteSliceElement(t.name, i); err != nil {
			t.err = base.NewLabel(err.Error())
		}
	}
	t.List = nil // rebuild on next draw
	return -1
}
func (t *Table) singleSelection() int { // return the first selection, substract title
	sel := t.List.Selection()
	if len(sel) > 0 && sel[0] == 0 {
		sel = sel[1:]
	}
	if len(sel) == 0 {
		return -1
	}
	return sel[0] - 1
}
func (t *Table) moveitem(up bool) {
	if len(t.List.List) < 3 { // title, firstitem, ...
		return
	}
	k := t.singleSelection() // k: -1..L list: 0..L+1 (header line)
	if (up && k < 1) || (!up && k > len(t.List.List)-1) {
		return
	}
	o := k + 1
	if up {
		o = k - 1
	}
	if err := (data{t.Data}).swapSliceElements(t.name, k-1, o-1); err != nil {
		t.err = base.NewLabel(err.Error())
		return
	}
	t.List.List[o], t.List.List[k] = t.List.List[k], t.List.List[o]
	t.List.SelectSingle(o)
}

func init() {
	base.RegisterIcon("edit", icons.ImageEdit)
	base.RegisterIcon("add", icons.ContentAdd)
	base.RegisterIcon("del", icons.ActionDelete)
	base.RegisterIcon("up", icons.HardwareKeyboardArrowUp)
	base.RegisterIcon("down", icons.HardwareKeyboardArrowDown)
}
