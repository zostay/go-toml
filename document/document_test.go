package document_test

import (
	"fmt"

	"github.com/pelletier/go-toml/v2/document"
)

func ExampleDocument_walk() {
	// TODO (https://golang.org/src/go/ast/walk.go)
}

func ExampleDocument_getAt() {
	doc := document.Document{
		KeyValues: []*document.KeyValue{
			{
				Key: document.StringKey("array"),
				Value: &document.Array{
					Elements: []document.ArrayElement{
						{Value: &document.String{Value: "zero"}},
						{Value: &document.String{Value: "one"}},
						{Value: &document.String{Value: "two"}},
					},
				},
			},
		},
		Tables: []*document.Table{
			{
				Key: document.StringKey("a", "b"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("c"),
						Value: &document.String{Value: "value"},
					},
				},
			},
		},
	}

	// Can retrieve explicit tables.
	fmt.Println("table:", doc.GetAt("a", "b"))
	// Can retrieve leaf nodes.
	fmt.Println("leaf:", doc.GetAt("a", "b", "c"))
	// Returns nil for nonexistent nodes.
	fmt.Println("nonexistent:", doc.GetAt("doesnotexist"))
	// Does not retrieve implicit tables.
	fmt.Println("implicit:", doc.GetAt("a"))
	// Can use index to get inside an array.
	fmt.Println("index:", doc.GetAt("array", 1))
	// Index outside of the range of an array returns nil.
	fmt.Println("oob:", doc.GetAt("array", 42))
	// Index can be -1 to mean the last element of the array.
	fmt.Println("last:", doc.GetAt("array", -1))

	// Output:
	// table: {{"a", "b"}, {"c", "value}}
	// leaf: {"value"}
	// nonexistent: nil
	// implicit: nil
	// index: {"one"}
	// oob: nil
	// last: {"two"}

}

func ExampleDocument_arrayTable() {
	doc := document.Document{
		Tables: []*document.Table{
			{
				Array: true,
				Key:   document.StringKey("products"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("name"),
						Value: &document.String{Value: "Hammer"},
					},
					{
						Key:   document.StringKey("sku"),
						Value: &document.Integer{V: "738594937"},
					},
				},
			},
			{
				Array: true,
				Key:   document.StringKey("products"),
				Comment: document.Comment{
					Value:  "empty table within the array",
					Inline: true,
				},
			},
			{
				Array: true,
				Key:   document.StringKey("products"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("name"),
						Value: &document.String{Value: "Nail"},
					},
					{
						Key:   document.StringKey("sku"),
						Value: &document.Integer{V: "284758393"},
					},
					{
						Key:   document.StringKey("color"),
						Value: &document.String{Value: "gray"},
					},
				},
			},
		},
	}

	fmt.Printf("%+v", doc)

	// Output:
	// [[products]]
	// name = "Hammer"
	// sku = 738594937
	//
	// [[products]] # empty table within the array
	//
	// [[products]]
	// name = "Nail"
	// sku = 284758393
	// color = "gray"
}

func ExampleDocument_reference() {
	doc := document.Document{
		KeyValues: []*document.KeyValue{
			{
				Key:   document.StringKey("title"),
				Value: &document.String{Value: "TOML Example"},
			},
		},
		Tables: []*document.Table{
			{
				Key: document.StringKey("owner"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("name"),
						Value: &document.String{Value: "Tom Preston-Werner"},
					},
					// TODO: dob
				},
			},
			{
				Key: document.StringKey("database"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("enabled"),
						Value: &document.Boolean{V: true},
					},
					{
						Key: document.StringKey("ports"),
						Value: &document.Array{
							Elements: []document.ArrayElement{
								{Value: &document.Integer{V: "8000"}},
								{Value: &document.Integer{V: "8001"}},
								{Value: &document.Integer{V: "8002"}},
							},
						},
					},
					{
						Key: document.StringKey("data"),
						Value: &document.Array{
							Elements: []document.ArrayElement{
								{Value: &document.Array{
									Elements: []document.ArrayElement{
										{Value: &document.String{Value: "delta"}},
										{Value: &document.String{Value: "phi"}},
									},
								}},
								{Value: &document.Array{
									Elements: []document.ArrayElement{
										// TODO floats
										// document.Float{V: "3.14"},
									},
								}},
							},
						},
					},
				},
			},
			{
				Key: document.StringKey("servers"),
			},
			{
				Key: document.StringKey("servers", "alpha"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("ip"),
						Value: &document.String{Value: "127.0.0.1"},
					},
					{
						Key:   document.StringKey("role"),
						Value: &document.String{Value: "frontend"},
					},
				},
			},
			{
				Key: document.StringKey("servers", "beta"),
				Elements: []*document.KeyValue{
					{
						Key:   document.StringKey("ip"),
						Value: &document.String{Value: "127.0.0.2"},
					},
					{
						Key:   document.StringKey("role"),
						Value: &document.String{Value: "backend"},
					},
				},
			},
		},
	}

	fmt.Println(doc)

	// Output:
	// title = "TOML Example"
	//
	// [owner]
	// name = "Tom Preston-Werner"
	// dob = 1979-05-27T07:32:00-08:00
	//
	// [database]
	// enabled = true
	// ports = [ 8000, 8001, 8002 ]
	// data = [ ["delta", "phi"], [3.14] ]
	// temp_targets = { cpu = 79.5, case = 72.0 }
	//
	// [servers]
	//
	// [servers.alpha]
	// ip = "10.0.0.1"
	// role = "frontend"
	//
	// [servers.beta]
	// ip = "10.0.0.2"
	// role = "backend"
}
