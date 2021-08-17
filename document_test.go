package toml_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
)

func doc(entities ...toml.Entity) toml.Document {
	return toml.Document{
		Root: entities,
	}
}

func TestDocument(t *testing.T) {
	examples := []struct {
		name string
		toml string
		doc  toml.Document
		err  error
	}{
		{
			name: "assign decimal int",
			toml: `x = 42`,
			doc: doc(
				&toml.KeyValue{
					Key:   []string{"x"},
					Value: &toml.Integer{Value: 42},
				},
			),
			err: nil,
		},
		{
			name: "assign string",
			toml: `x = "hello"`,
			doc: doc(
				&toml.KeyValue{
					Key: []string{"x"},
					Value: &toml.String{
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
				&toml.KeyValue{
					Key:   []string{"a"},
					Value: &toml.String{Value: "hello"},
				},
				&toml.KeyValue{
					Key:   []string{"b"},
					Value: &toml.Integer{Value: 42},
				},
			),
			err: nil,
		},
		{
			name: "table",
			toml: `[a]`,
			doc: doc(
				&toml.Table{
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
				&toml.Table{
					Key: []string{"a"},
				},
				&toml.KeyValue{
					Key:   []string{"b"},
					Value: &toml.Integer{Value: 1},
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
				&toml.Table{
					Key: []string{"a"},
				},
				&toml.KeyValue{
					Key:   []string{"b"},
					Value: &toml.Integer{Value: 1},
				},
				&toml.KeyValue{
					Key:   []string{"c"},
					Value: &toml.Integer{Value: 2},
				},
			),
		},
		{
			name: "table with implicit intermediate",
			toml: `[a.b]
		c = 1`,
			doc: doc(
				&toml.Table{
					Key: []string{"a", "b"},
				},
				&toml.KeyValue{
					Key:   []string{"c"},
					Value: &toml.Integer{Value: 1},
				},
			),
			err: nil,
		},
	}

	for _, e := range examples {
		t.Run(e.name, func(t *testing.T) {
			d, err := toml.Parse([]byte(e.toml))
			if e.err != nil {
				require.Equal(t, e.err, err)
			} else {
				d.Viz()
				require.Equal(t, e.doc, d)
			}
		})
	}
}
