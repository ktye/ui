package property

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Data links the property system application data.
// Src must be a pointer to a struct that contains structs or slices of structs.
// Example
//	type Application struct {       // application data type
//		General GeneralSettings // simple struct properties
//		Clients []Client        // slice of structs
//		Orders  []Order
//	}
//	type GeneralSettings struct {
//		Hostname string   `name:"Host"`
//		Colors   []int    `all:"colors"` // A slice to a simple type is possible
//		Speed    float64
//		Active   bool     `hidden:"true"`
//		Another  CustomType
//              Multiple []CustomType
//	}
//	var data Application
//	prop := property.New(&data)
//
// There are 3 ways to present a property dialog:
//	- a property list dialog for a single struct (e.g. General)
//	- a property table for a slice (e.g. Clients)
//		clicking on a line would show a list dialog.
//	- a property list dialog for multiple items (e.g. selected clients)
//		only custom values can be edited directly.
//		diverging values are not editable by default and must be activated
//		on Update, only changed values are overwritten
//
// Struct Tags:
//	tip(string)         tooltip or help text
//	editable:"true"     editable pulldown menu
//	textbox:"true"      show multiline string in a textbox
//	options(string)     key to the GetOptions method which returns combobox values
//      all(string)         key to the GetAll method which returns combobox values
//	browse:"file"|"dir" string input field has a helper for a file or dir browser
//      view:"ignored"      not shown in the dialog
//      hidden:"true"       shown only if non-zero, or if all values are requested (advanced button)
//
// Custom types
//	The property system is based on strings. Custom type can be used if they have a String method, or are
//	convertible to strings (string kinds). A custom type also must implement a Verify(string)error method,
//      that is used for syntax checking only.
//
// Name and Rename
//	A property key/value pair may behave as an ID that is referenced by others.
//	For this case a field with the name "Name" is treated specially if it has a string kind value of custom type.
//	If this field is updated, it triggers the Rename method of the data source, that receives the old value
//	of the custom id value as a reflection value and the new name.
//	The application has the chance to update all references to the old id to the new name.
type data struct {
	src Source
}

// Source is the interface that the application data must implement.
// It also must provide it's data in the form described in the example above.
type Source interface {
	GetAll(string) ([]string, error)
	GetOptions(string) ([]string, error)
	RenameID(reflect.Value, string) error
	DeleteID(reflect.Value) error
	PreUpdate()
	PostUpdate()
}

// List is a collection of Properties to be presented in a property dialog.
// It may present an interface to a single data struct or a selection from a slice of structs.
type list struct {
	Name          string // Name of the primary data field, such as "General" or "Clients"
	IsStruct      bool   // true for pure structs such as "General", false for slices such as "Clients"
	Fields        []property
	updateIndexes []int // Indexes of data slice to be updated.
}

// Table represents a slice of property lists.
// All items are assumed to have the same structure, to be presented in a table widget with one row per list.
type table []list

// Property is a single property key value pair.
type property struct {
	FieldName       string       // original field name
	DisplayName     string       // alternative display name (struct tag "name")
	Tip             string       // tooltip text (struct tag "tip")
	Type            reflect.Type // property type the field or field slice element
	IsSlice         bool         // value is a vector of the base type.
	IsTextbox       bool         // string value is multi-line
	IsHidden        bool         // hide property as an advanced option
	IsIgnored       bool         // can be updated but is ignored in dialog (struct tag: view:"ignored")
	Options         []string     // options for combo boxes
	EditableOptions bool         // options combobox should be editable
	IsPassword      bool         // hide input string
	IsUnique        bool         // All indexes have the same values
	Browse          string       // string value can be obtained with a "file" or "directory" browser
	IsUpdated       bool         // value is updated
	ZeroString      string       // initial value for appending new values
	Values          []string     // String represenation of value or values (if IsSlice).
}

// List returns a property List from the data source's struct field with the given name.
// If copy indexes is not nil, it copies the values from the first copyIndexes.
// For each Property it checks if the values from all copyIndexes are the same which sets IsUnique.
// It also sets the updateIndexes: If copyIndexes contains multiple values,
// these are copied to the update indexes.
// If isNew is true, the Property with the field called "Name" gets a unique value and
// it's update index is the the slice length.
// Otherwise (isNew is false, and copyIndexes has length 1) the single index is updated.
func (d data) list(name string, copyIndexes []int, newName bool) (list, error) {
	ps := list{Name: name}
	v, sliceLen, isStruct := d.newDataStructValue(name)
	if sliceLen == -1 {
		return ps, fmt.Errorf("property: unknown id: '%s'", name)
	}
	ps.IsStruct = isStruct

	for _, i := range copyIndexes {
		if i < 0 || i >= sliceLen {
			return ps, fmt.Errorf("property: %s: cannot copy from element %d (max: %d)", name, i, sliceLen)
		}
	}
	if newName == false && copyIndexes == nil {
		copyIndexes = []int{0}
	}
	if newName == true && len(copyIndexes) > 1 {
		return ps, fmt.Errorf("property: %s: newName is not allowed for multiple copy indexes", name)
	}
	if newName == true {
		ps.updateIndexes = []int{sliceLen}
	} else if len(copyIndexes) == 1 {
		ps.updateIndexes = []int{copyIndexes[0]}
	} else {
		ps.updateIndexes = make([]int, len(copyIndexes))
		copy(ps.updateIndexes, copyIndexes)
	}

	fieldnames := dataStructFieldNames(v)
	if len(fieldnames) == 0 {
		return ps, fmt.Errorf("property: %s has no fields", name)
	}

	ps.Fields = make([]property, len(fieldnames))
	n := 0
	for _, fieldname := range fieldnames {
		p, err := d.propertyTypeInfo(v, fieldname)
		if err != nil {
			return ps, err
		}
		p.IsUnique = true
		p.ZeroString = p.zeroString()
		if len(copyIndexes) == 0 {
			// Initialize with empty values.
			if p.IsSlice == false {
				p.Values = []string{p.zeroString()}
			}
		} else {
			// Copy values from the first copyIndex
			if values, err := d.getPropertyValues(name, copyIndexes[0], fieldname); err != nil {
				return ps, fmt.Errorf("property: %s", err)
			} else {
				p.Values = values
			}
			// Check if all copyIndexes have idential values.
			if len(copyIndexes) > 1 {
				for i := 1; i < len(copyIndexes); i++ {
					if values, err := d.getPropertyValues(name, copyIndexes[i], fieldname); err != nil {
						return ps, fmt.Errorf("property: %s", err)
					} else {
						if len(values) != len(p.Values) {
							p.IsUnique = false
						} else if p.IsUnique {
							for i, val := range p.Values {
								if values[i] != val {
									p.IsUnique = false
								}
							}
						}
					}
				}
			}
			if p.hasZeroValues() == false {
				p.IsHidden = false
			}
		}
		// The Name field is treated specially:
		// Skip if multiple fields are edited len(copyIndexes) > 1,
		// Set to a unique name, if newName is true.
		if fieldname == "Name" {
			if len(copyIndexes) > 1 {
				continue
			}
			if newName {
				p.Values = []string{d.newPropertyName(name, fieldname)}
			}
		}
		// For a new struct, we need to Update all fields.
		if newName {
			p.IsUpdated = true
		}
		ps.Fields[n] = p
		n++
	}
	ps.Fields = ps.Fields[:n]
	return ps, nil
}

// Table returns a property table from the data source's struct field with the given name, which must be a slice value.
func (d data) table(name string) (table, error) {
	_, sliceLen, isStruct := d.newDataStructValue(name)
	if sliceLen == -1 {
		return nil, fmt.Errorf("property: not a slice: '%s'", name)
	}
	if sliceLen == 0 {
		return nil, nil
	}
	if isStruct {
		return nil, fmt.Errorf("GetAllProperties '%s' is not a slice", name)
	}
	t := make(table, sliceLen)
	for i := range t {
		if p, err := d.list(name, []int{i}, false); err != nil {
			return nil, err
		} else {
			t[i] = p
		}
	}
	return t, nil
}
func (l list) NameField() string {
	for _, p := range l.Fields {
		if p.FieldName == "Name" && len(p.Values) == 1 {
			return p.Values[0]
		}
	}
	return ""
}

// Verify all property fields.
func (l list) Verify() error {
	for _, p := range l.Fields {
		for _, s := range p.Values {
			if err := p.Verify(s); err != nil {
				return err
			}
		}
	}
	return nil
}

// Update updates a data field that is represented by the List.
// Only property fields which are marked with IsUpdated are set.
// Each Property carries it's own UpdateIndex, which indicates which slice index should be updated.
func (d data) update(ps list) error {
	d.src.PreUpdate()
	defer d.src.PostUpdate()

	isNew := false
	v := reflect.ValueOf(d.src).Elem().FieldByName(ps.Name)
	var val reflect.Value
	for _, i := range ps.updateIndexes {
		switch v.Kind() {
		case reflect.Slice:
			if i == v.Len() {
				// Append slice element.
				isNew = true
				newElement := reflect.Zero(v.Type().Elem())
				newSlice := reflect.Append(v, newElement)
				v.Set(newSlice)
			} else if i < 0 || i > v.Len() {
				return fmt.Errorf("%s[%d]: index is out of range", ps.Name, i)
			}
			val = v.Index(i)
		case reflect.Struct:
			if i != 0 {
				return fmt.Errorf("property: %s is not a slice (element %d requested)", ps.Name, i)
			}
			val = v
		default:
			return fmt.Errorf("property: unknown kind: %s: %v", ps.Name, v.Kind())
		}

		if val.CanSet() == false {
			return fmt.Errorf("property: %s[%d] is not settable", ps.Name, i)
		}

		for _, p := range ps.Fields {
			if p.IsUpdated {
				fieldValue := val.FieldByName(p.FieldName)
				if fieldValue == (reflect.Value{}) {
					return fmt.Errorf("%s[%d].%s: field does not exist", ps.Name, i, p.FieldName)
				}
				if p.FieldName == "Name" {
					if len(p.Values) != 1 {
						return fmt.Errorf("property: Name field must have a single value")
					} else if p.Values[0] == "" {
						return fmt.Errorf("property: Name cannot be empty")
					}
					newName := p.Values[0]
					if d.isUniqueName(ps.Name, i, newName) == false {
						return fmt.Errorf("property: Name (%s) already exists", newName)
					}
					// Don't rename new fields, otherwise this renames an existing empty id.
					// We need to make a copy before renaming.
					if isNew == false {
						old := reflect.New(fieldValue.Type()).Elem()
						old.Set(fieldValue)
						if err := d.src.RenameID(old, newName); err != nil {
							return err
						}
					}
				}
				if err := p.set(fieldValue); err != nil {
					return fmt.Errorf("%s[%d].%s: %s", ps.Name, i, p.FieldName, err)
				}
			}
		}
	}
	return nil
}

func (d data) isUniqueName(dataFieldName string, sliceIndex int, newValue string) bool {
	v := reflect.ValueOf(d.src).Elem().FieldByName(dataFieldName)
	// We are only interested in slices.
	// A single struct is already unique.
	if v.Kind() == reflect.Slice {
		for i := 0; i < v.Len(); i++ {
			if i != sliceIndex {
				if s := v.Index(i).FieldByName("Name").String(); s == newValue {
					return false
				}
			}
		}
	}
	return true
}

// GetNames returns the values of the Name field for all elements of a data struct slice.
func (d data) GetNames(name string) []string {
	v := reflect.ValueOf(d.src).Elem().FieldByName(name)
	if v.Kind() != reflect.Slice {
		return nil
	}
	names := make([]string, v.Len())
	for i := range names {
		sv := v.Index(i)
		vname := sv.FieldByName("Name")
		if vname == (reflect.Value{}) {
			names[i] = "?"
		} else {
			names[i] = vname.String()
		}
	}
	return names
}

func (l list) FieldIndex(fieldname string) int {
	for i, p := range l.Fields {
		if p.FieldName == fieldname {
			return i
		}
	}
	return -1
}

// Set sets the field value v.
// V may be a slice.
func (p *property) set(v reflect.Value) error {
	setElement := func(element reflect.Value, s string) error {
		switch element.Kind() {
		case reflect.Bool:
			if v, err := strconv.ParseBool(s); err != nil {
				return err
			} else {
				element.Set(reflect.ValueOf(v))
			}
		case reflect.Int:
			if v, err := strconv.Atoi(s); err != nil {
				return fmt.Errorf("%s: not an integer", s)
			} else {
				element.Set(reflect.ValueOf(v))
			}
		case reflect.Float64:
			s = strings.Replace(s, ",", ".", 1) // allow both , and . as decimal delimiters.
			if v, err := strconv.ParseFloat(s, 64); err != nil {
				return fmt.Errorf("%s: not a number", s)
			} else {
				element.Set(reflect.ValueOf(v))
			}
		default:
			/* TODO: is this used? */
			// Call FromString method if it exists.
			// FromString must return an error as a string type.
			if _, ok := element.Addr().Type().MethodByName("FromString"); ok {
				errStrings := element.Addr().MethodByName("FromString").Call([]reflect.Value{reflect.ValueOf(s)})
				errstr := errStrings[0].String()
				if errstr != "" {
					return fmt.Errorf("%s", errstr)
				}
				return nil
			}

			if reflect.TypeOf(s).ConvertibleTo(element.Type()) == false {
				return fmt.Errorf("cannot convert string to %v", element.Type())
			}
			v := reflect.ValueOf(s).Convert(element.Type())
			element.Set(v)
		}
		return nil
	}

	if v.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(v.Type(), len(p.Values), len(p.Values))
		for i, s := range p.Values {
			if err := setElement(slice.Index(i), s); err != nil {
				return err
			}
		}
		v.Set(slice)
	} else {
		if len(p.Values) != 1 {
			return fmt.Errorf("cannot set %d values to a non-slice value", len(p.Values))
		}
		return setElement(v, p.Values[0])
	}
	return nil
}

// zeroString returns the zeroValue for a given type.
// A custom type may have a ZeroString method that is used as a default value, instead of "".
func (p *property) zeroString() string {
	if p.Type == reflect.TypeOf(false) {
		return "false"
	} else if p.Type == reflect.TypeOf(0) {
		return "0"
	} else if p.Type == reflect.TypeOf(0.0) {
		return "0"
	} else if p.Type == reflect.TypeOf("") {
		return ""
	} else {
		v := reflect.New(p.Type)
		if f := v.MethodByName("ZeroString"); f == (reflect.Value{}) {
			return ""
		} else {
			ret := f.Call(nil)
			if len(ret) != 1 {
				return ""
			} else {
				return ret[0].String()
			}
		}
	}
}

// Verify verifies all string Values for the given types.
func (p *property) Verify(s string) error {
	boolType := reflect.TypeOf(false)
	intType := reflect.TypeOf(0)
	float64Type := reflect.TypeOf(0.0)
	stringType := reflect.TypeOf("")
	switch p.Type {
	case boolType:
		if _, err := strconv.ParseBool(s); err != nil {
			return fmt.Errorf("%s: %s: value must be true or false", p.DisplayName, s)
		}
		return nil
	case intType:
		if _, err := strconv.Atoi(s); err != nil {
			return fmt.Errorf("%s: %s: not an integer", p.DisplayName, s)
		}
		return nil
	case float64Type:
		s = strings.Replace(s, ",", ".", 1)
		if _, err := strconv.ParseFloat(s, 64); err != nil {
			return fmt.Errorf("%s: %s is not a float", p.DisplayName, s)
		}
		return nil
	case stringType:
		return nil
	default:
		v := reflect.New(p.Type)
		if f := v.MethodByName("Verify"); f == (reflect.Value{}) {
			return fmt.Errorf("Type has no verifier: %v", p.Type)
		} else {
			ret := f.Call([]reflect.Value{reflect.ValueOf(s)})
			if len(ret) != 1 {
				return fmt.Errorf("wrong verifier for type %v: returns wrong number of arugments", p.Type)
			} else {
				iface := ret[0].Interface()
				if e, ok := iface.(error); ok {
					return fmt.Errorf("%s: %s", p.DisplayName, e)
				} else {
					return nil
				}
			}
		}
	}
	return nil
}

// getPropertyValues returns the Property with the given name and field
// converted to a string slice.
func (d data) getPropertyValues(name string, index int, fieldname string) ([]string, error) {
	v := reflect.ValueOf(d.src).Elem().FieldByName(name)
	var valueStruct reflect.Value
	switch v.Kind() {
	case reflect.Slice:
		if index < 0 || index >= v.Len() {
			return nil, fmt.Errorf("%s[%d]: index out of range", name, index)
		}
		valueStruct = v.Index(index)
	case reflect.Struct:
		valueStruct = v
	default:
		return nil, fmt.Errorf("unknown primary field type: %v", v.Kind())
	}

	toString := func(element reflect.Value) string {
		var s string
		stringType := reflect.TypeOf(s)
		_, isStringer := element.Type().MethodByName("String")
		if element.Kind() == reflect.Bool {
			return strconv.FormatBool(element.Bool())
		} else if element.Kind() == reflect.Int {
			return strconv.FormatInt(element.Int(), 10)
		} else if element.Kind() == reflect.Float64 {
			return strconv.FormatFloat(element.Float(), 'g', -1, 64)
		} else if isStringer {
			// This panics, if the String method takes input arguments, or returns anything else than a single string.
			stringValues := element.MethodByName("String").Call(nil)
			return stringValues[0].String()
		}
		if element.Type().ConvertibleTo(stringType) {
			s = element.Convert(stringType).Interface().(string)
			return s
		}
		return "<?>"
	}

	var values []string
	var zero reflect.Value
	field := valueStruct.FieldByName(fieldname)
	if field == zero {
		return nil, fmt.Errorf("%s.%s: field does not exist", name, fieldname)
	}
	if field.Kind() == reflect.Slice {
		values = make([]string, field.Len())
		for i := range values {
			values[i] = toString(field.Index(i))
		}
	} else {
		values = []string{toString(field)}
	}
	return values, nil
}

// newPropertyName returns an ID for a Properties value, which is not taken already.
// The name is the first integer starting with the slice length + 1.
func (d data) newPropertyName(name, fieldname string) string {
	v := reflect.ValueOf(d.src).Elem().FieldByName(name)
	// We are only interested in slices.
	// A single struct is already unique.
	if v.Kind() == reflect.Slice {
		usedNames := make(map[string]bool)
		for i := 0; i < v.Len(); i++ {
			s := v.Index(i).FieldByName("Name").String()
			usedNames[s] = true
		}
		for i := v.Len(); ; i++ {
			s := strconv.Itoa(i + 1)
			if usedNames[s] == false {
				return s
			}
		}
	}
	return "1"
}

// newDataStructValue returns a reflect object for an empty primary data struct and
// the slice length, it the primary data struct is a slice.
// It returns -1 as the length in case of errors.
// The last return argument, is if the object is a struct (otherwise it's a slice).
func (d data) newDataStructValue(dataFieldName string) (reflect.Value, int, bool) {
	v := reflect.ValueOf(d.src).Elem().FieldByName(dataFieldName)
	if v == (reflect.Value{}) {
		return reflect.Value{}, -1, false
	}
	t := v.Type()
	switch v.Kind() {
	case reflect.Slice:
		return reflect.New(t.Elem()).Elem(), v.Len(), false
	case reflect.Struct:
		return reflect.New(t).Elem(), 1, true
	default:
		return reflect.Value{}, -1, true
	}
}

// swapSliceElements swaps two slice elements of a primary data struct.
func (d data) swapSliceElements(dataFieldName string, i, j int) error {
	slice := reflect.ValueOf(d.src).Elem().FieldByName(dataFieldName)
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("property: swap: struct field '%s' is not a slice", dataFieldName)
	}
	l := slice.Len()
	if i < 0 || j < 0 || i >= l || j >= l {
		return fmt.Errorf("property: swap: s '%s': index out of range", dataFieldName)
	}
	tmp := reflect.New(slice.Type().Elem()).Elem()
	tmp.Set(slice.Index(i))
	slice.Index(i).Set(slice.Index(j))
	slice.Index(j).Set(tmp)
	return nil
}

// deleteSliceElements delete the element at the given index of a primary data struct.
func (d data) deleteSliceElement(dataFieldName string, i int) error {
	slice := reflect.ValueOf(d.src).Elem().FieldByName(dataFieldName)
	if slice.Kind() != reflect.Slice {
		return fmt.Errorf("property: delete %s[%d]; not a slice", dataFieldName, i)
	}
	if i < 0 {
		return fmt.Errorf("property: delete %s[%d]", dataFieldName, i)
	}
	n := slice.Len()
	if i >= n {
		return fmt.Errorf("property: delete %s[%d] >= len", dataFieldName, i)
	}

	// notify DeleteID
	e := slice.Index(i)
	if e.Kind() != reflect.Struct {
		return fmt.Errorf("property: delete %s[%d]: not a struct", dataFieldName, i)
	}
	var zv reflect.Value
	namefield := e.FieldByName("Name")
	if namefield != zv {
		if err := d.src.DeleteID(namefield); err != nil {
			return err
		}
	}

	if i == n-1 {
		slice.SetLen(n - 1)
	} else {
		newSlice := reflect.AppendSlice(slice.Slice(0, i), slice.Slice(i+1, n))
		slice.Set(newSlice)
	}
	return nil
}

func (d data) propertyTypeInfo(dataStruct reflect.Value, fieldname string) (property, error) {
	p := property{
		FieldName:   fieldname,
		DisplayName: fieldname,
	}

	sf, _ := dataStruct.Type().FieldByName(fieldname)
	t := sf.Type

	if s := sf.Tag.Get("name"); s != "" {
		p.DisplayName = s
	}

	link := regexp.MustCompile(`\[([^\]]+)\]\(([^\)]+)\)`)
	p.Tip = link.ReplaceAllString(sf.Tag.Get("tip"), `$1`)

	if t.Kind() == reflect.Slice {
		p.IsSlice = true
		t = t.Elem()
	}

	p.Type = t
	if p.Type == reflect.TypeOf(false) {
		//p.Options = []string{"false", "true"}
	}

	if sf.Tag.Get("hidden") != "" {
		p.IsHidden = true
	}

	if sf.Tag.Get("view") == "ignored" {
		p.IsIgnored = true
	}

	if key := sf.Tag.Get("options"); key != "" {
		if o, err := d.src.GetOptions(key); err != nil {
			return p, fmt.Errorf("GetProperties '%s': %s", fieldname, err)
		} else {
			p.Options = o
		}
	}
	if key := sf.Tag.Get("all"); key != "" {
		if all, err := d.src.GetAll(key); err != nil {
			return p, fmt.Errorf("GetProperties '%s': %s", fieldname, err)
		} else {
			p.Options = all
		}
	}
	if sf.Tag.Get("editable") != "" {
		p.EditableOptions = true
	}
	if sf.Tag.Get("password") != "" {
		p.IsPassword = true
	}
	if sf.Tag.Get("textbox") != "" {
		p.IsTextbox = true
	}
	p.Browse = sf.Tag.Get("browse")
	return p, nil
}

func (p property) hasZeroValues() bool {
	if p.IsSlice {
		return len(p.Values) == 0
	}
	if len(p.Values) != 1 {
		return false
	}
	boolType := reflect.TypeOf(false)
	intType := reflect.TypeOf(0)
	float64Type := reflect.TypeOf(0.0)
	s := p.Values[0]
	if p.Type == boolType {
		if s == "false" {
			return true
		}
		return false
	}
	if p.Type == intType {
		if s == "0" {
			return true
		}
		return false
	}
	if p.Type == float64Type {
		if s == "0.0" || s == "0" {
			return true
		}
		return false
	}
	if s == "" {
		return true
	}
	return false
}

// dataStructFieldNames returns all fieldnames for a Properties struct.
// Lowercase names are skipped.
func dataStructFieldNames(dataStruct reflect.Value) []string {
	fields := make([]string, dataStruct.NumField())
	structType := dataStruct.Type()
	n := 0
	for i := 0; i < len(fields); i++ {
		sf := structType.Field(i)
		name := sf.Name
		if r, _ := utf8.DecodeRuneInString(name); unicode.IsUpper(r) == false {
			continue
		}
		fields[n] = name
		n++
	}
	return fields[:n]
}
