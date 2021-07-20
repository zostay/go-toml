package toml

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/bradleyjkemp/memviz"
	"github.com/pelletier/go-toml/v2/internal/ast"
)

type Document struct {
	Root []Entity
}

func Parse(b []byte) (Document, error) {
	d := Document{}

	p := parser{}
	p.Reset(b)
	p.Comments = true

	for p.NextExpression() {
		expr := p.Expression()
		switch expr.Kind {
		case ast.KeyValue:
			e := &KeyValue{}

			// Handle key
			k := expr.Key()
			// TODO: this assumes just one element in the key
			for k.Next() {
				e.K = string(k.Node().Data)
			}

			// Handle value
			v := expr.Value()
			var err error
			e.V, err = entityFromTerminalNode(v)
			if err != nil {
				return d, err
			}

			d.Root = append(d.Root, e)
		default:
			panic(fmt.Errorf("unhandled node kind %s", expr.Kind))
		}
	}

	return d, p.Error()
}

func entityFromTerminalNode(n *ast.Node) (Entity, error) {
	switch n.Kind {
	case ast.Integer:
		v, err := parseInteger(n.Data)
		if err != nil {
			return nil, err
		}
		return &Integer{
			Value: v,
		}, nil
	case ast.String:
		return &String{
			Value: string(n.Data),
		}, nil
	default:
		panic(fmt.Errorf("unhandled node kind %s", n.Kind))
	}
}

// TODO: remove me
func (d *Document) Viz() {
	buf := &bytes.Buffer{}
	memviz.Map(buf, d)
	err := ioutil.WriteFile("go-toml-document-dump.dot", buf.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}

type Entity interface {
}

type Comment struct {
	Text   string
	Inline bool
}

type Table struct {
	Implicit bool
	K        string
}

type KeyValue struct {
	C *Comment

	K string
	V Entity // V is one of the terminal type
}

type String struct {
	Multiline bool
	Literal   bool
	Value     string
}

type Integer struct {
	Value int64
}
