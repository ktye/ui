package main

import (
	"fmt"
	"reflect"
)

// This data structure is an example of an application.
// The application could as well be in a separate package, that does not know about ui.
// All it has to do, is follow the convensions (structs, slices of structs, max nesting, struct tags...)

func NewApp() *app {
	var a app
	a.General.FontSize = 12
	a.Clients = []Client{
		Client{"John", "North", false, []OrderID{"866"}},
		Client{"Emma", "North", true, []OrderID{"123", "866"}},
		Client{"Peter", "West", true, []OrderID{"732", "980"}},
		Client{"Hans-Jürgen", "South", true, []OrderID{"923"}},
	}
	a.Orders = []Order{
		Order{"123", 1, 256.3, []Discount{"10%", "20%"}, nil, "telephone"},
		Order{"421", 1, 89.5, nil, []string{"a", "b"}, "camera"},
		Order{"732", 2, 5.98, nil, nil, "inner tube"},
		Order{"866", 300, 1.35, []Discount{"50%"}, nil, "M4x10"},
		Order{"923", 1, 105.5, nil, nil, "camera"},
		Order{"980", 10, 5.60, nil, nil, "M4x10"},
	}
	a.Products = []ProductID{"telephe", "camera", "inner tube", "M4x10"}
	a.current = map[string]int{
		"Clients": 1,
		"Orders":  3,
	}
	return &a
}

type app struct {
	General  Settings
	Clients  []Client
	Orders   []Order
	Products []ProductID
	current  map[string]int
}

type Settings struct {
	Dark     bool
	FontSize int
}

type Client struct {
	Name        string
	Region      Region `all:"region"`
	Trustworthy bool
	Orders      []OrderID `all:"orders"`
}

type Order struct {
	Name      OrderID
	Quantity  int
	Price     float64
	Discounts []Discount `all:"discount"`
	Tags      []string
	Product   ProductID `all:"product"`
}

type Region string
type OrderID string
type ProductID string
type Discount string

// Implement the property.Source interface:

func (a *app) GetAll(s string) ([]string, error) {
	switch s {
	case "product":
		v := make([]string, len(a.Products))
		for i := range a.Products {
			v[i] = string(a.Products[i])
		}
		return v, nil
	case "region":
		return []string{"North", "South", "East", "West"}, nil
	case "discount":
		return []string{"10%", "20%", "50%"}, nil
	case "orders":
		s := make([]string, len(a.Orders))
		for i := range s {
			s[i] = string(a.Orders[i].Name)
		}
		return s, nil
	}
	return nil, fmt.Errorf("getall: key does not exist: %s", s)
}
func (a *app) GetOptions(s string) ([]string, error)    { return nil, fmt.Errorf("there are no options") }
func (a *app) RenameID(v reflect.Value, s string) error { return nil } // Triggered after rename of a "Name" property
func (a *app) DeleteID(v reflect.Value) error           { return nil } // Triggered after deleting a property with the "Name" value
func (a *app) GetCurrent(s string) int                  { return a.current[s] }
func (a *app) SetCurrent(s string, n int)               { a.current[s] = n }
func (a *app) PreUpdate()                               { println("lock") }
func (a *app) PostUpdate()                              { println("unlock") }
