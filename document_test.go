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
					K: "x",
					V: &toml.Integer{Value: 42},
				},
			),
			err: nil,
		},
		{
			name: "assign string",
			toml: `x = "hello"`,
			doc: doc(
				&toml.KeyValue{
					K: "x",
					V: &toml.String{
						Value: "hello",
					},
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
