package property

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
	name string
	// TODO
}

// TODO
