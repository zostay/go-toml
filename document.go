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
		e, err := entityFromRootExpression(expr)
		if err != nil {
			return d, err
		}
		d.Root = append(d.Root, e)
	}

	return d, p.Error()
}

// TODO: need more finesse than this to conserve the different kind of keys.
func keyIteratorToStrings(it ast.Iterator) []string {
	key := []string{}
	for it.Next() {
		key = append(key, string(it.Node().Data))
	}
	return key
}

func entityFromRootExpression(e *ast.Node) (Entity, error) {
	switch e.Kind {
	case ast.Table:
		key := keyIteratorToStrings(e.Key())
		return &Table{
			Key: key,
		}, nil
	case ast.KeyValue:
		key := keyIteratorToStrings(e.Key())

		v, err := entityFromExpression(e.Value())
		if err != nil {
			return nil, err
		}

		return &KeyValue{
			Key:   key,
			Value: v,
		}, nil
	default:
		panic(fmt.Errorf("unhandled root expression kind %s", e.Kind))
	}
}

func entityFromExpression(e *ast.Node) (Entity, error) {
	switch e.Kind {
	case ast.Integer:
		v, err := parseInteger(e.Data)
		if err != nil {
			return nil, err
		}
		return &Integer{
			Value: v,
		}, nil
	case ast.String:
		return &String{
			Value: string(e.Data),
		}, nil
	default:
		panic(fmt.Errorf("unhandled expression kind %s", e.Kind))
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
	Inline bool
	Key    []string
}

type KeyValue struct {
	C *Comment

	Key   []string
	Value Entity
}

type String struct {
	Multiline bool
	Literal   bool
	Value     string
}

type Integer struct {
	Value int64
}
