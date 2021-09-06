package document_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2/x/document"
	"github.com/stretchr/testify/require"
)

func doc(entities ...document.Entity) document.Document {
	return document.Document{
		document.Container{
			Elements: entities,
		},
	}
}

func TestDocument(t *testing.T) {
	examples := []struct {
		name string
		toml string
		doc  document.Document
		err  error
	}{
		{
			name: "assign decimal int",
			toml: `x = 42`,
			doc: doc(
				&document.KeyValue{
					Key:   []string{"x"},
					Value: &document.Integer{Value: 42},
				},
			),
			err: nil,
		},
		{
			name: "assign string",
			toml: `x = "hello"`,
			doc: doc(
				&document.KeyValue{
					Key: []string{"x"},
					Value: &document.String{
						Value: "hello",
					},
				},
			),
			err: nil,
		},
		{
			name: "assign string and int",
			toml: `a = "hello"
b = 42`,
			doc: doc(
				&document.KeyValue{
					Key:   []string{"a"},
					Value: &document.String{Value: "hello"},
				},
				&document.KeyValue{
					Key:   []string{"b"},
					Value: &document.Integer{Value: 42},
				},
			),
			err: nil,
		},
		{
			name: "table",
			toml: `[a]`,
			doc: doc(
				&document.Table{
					Key: []string{"a"},
				},
			),
			err: nil,
		},
		{
			name: "table with one assign",
			toml: `[a]
b = 1`,
			doc: doc(
				&document.Table{
					Key: []string{"a"},
				},
				&document.KeyValue{
					Key:   []string{"b"},
					Value: &document.Integer{Value: 1},
				},
			),
		},
		{
			name: "table with two assigns",
			toml: `
[a]
b = 1
c = 2`,
			doc: doc(
				&document.Table{
					Key: []string{"a"},
				},
				&document.KeyValue{
					Key:   []string{"b"},
					Value: &document.Integer{Value: 1},
				},
				&document.KeyValue{
					Key:   []string{"c"},
					Value: &document.Integer{Value: 2},
				},
			),
		},
		{
			name: "table with implicit intermediate",
			toml: `[a.b]
		c = 1`,
			doc: doc(
				&document.Table{
					Key: []string{"a", "b"},
				},
				&document.KeyValue{
					Key:   []string{"c"},
					Value: &document.Integer{Value: 1},
				},
			),
			err: nil,
		},
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			d, err := document.Parse([]byte(e.toml))
			if e.err != nil {
				require.Equal(t, e.err, err)
			} else {
				d.Viz()
				require.Equal(t, e.doc, d)
			}
		})
	}
}
