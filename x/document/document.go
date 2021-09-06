package document

import (
	"fmt"

	"github.com/pelletier/go-toml/v2/internal/ast"
	"github.com/pelletier/go-toml/v2/internal/parser"
)

type Document struct {
	Container
}

func (d *Document) container() *Container {
	return &d.Container
}

func Parse(b []byte) (Document, error) {
	d := Document{}

	p := parser.Parser{}
	p.Reset(b)
	p.Comments = true

	var cursor Entity = &d

	for p.NextExpression() {
		expr := p.Expression()
		var err error

		switch expr.Kind {
		case ast.Table:
			cursor, err = docAddTable(&d, expr)
		case ast.ArrayTable:
			cursor = &d
			panic("not implemented")
		case ast.KeyValue:
			err = docAddKeyValue(cursor, expr)
		default:
			// TODO: add error context
			err = fmt.Errorf("expression of type '%s' not allowed there", expr.Kind)
		}

		if err != nil {
			return d, err
		}
	}

	return d, p.Error()
}

func docAddKeyValue(parent Entity, expr *ast.Node) error {

	return nil
}

// Returns the new cursor or an error.
func docAddTable(root Entity, expr *ast.Node) (Entity, error) {
	cursor := root
	key := expr.Key()

parts:
	for key.Next() {
		c, ok := cursor.(container)
		if !ok {
			return nil, fmt.Errorf("tried to use a key on a non-container element")
		}
		parent := c.container()
		name := string(key.Node().Data)

		for _, element := range parent.Elements {
			e, ok := element.(keyed)
			if !ok {
				continue
			}
			if e.key().Name == name {
				cursor = element
				continue parts
			}
		}

		newTable := Table{
			Key: Key{
				Name: name,
			},
		}
		parent.Elements = append(parent.Elements, newTable)
		cursor = parent.Elements[len(parent.Elements)-1]
	}

	return cursor, nil
}

type Entity interface {
}

// Container type is meant to be embedded in all TOML types that are contain
// other elements:
//
// Document (root), Table.
//
// It allows direct access to elements order by manipulation of the Elements
// slice.
type Container struct {
	Elements []Entity
}

// Private interface that needs to be implemented by structs that embed a
// Container. Used for runtime check and dispatch.
type container interface {
	container() *Container
}

type Key struct {
	Name string
	// TODO: merge into one.
	Quoted  bool
	Literal bool
}

type keyed interface {
	key() *Key
}

type Comment struct {
	Text   string
	Inline bool
}

type Table struct {
	Container

	Inline bool
	Key    Key
}

func (t *Table) container() *Container {
	return &t.Container
}

func (t *Table) key() *Key {
	return &t.Key
}

type KeyValue struct {
	C *Comment

	Key   Key
	Value Entity
}

func (kv *KeyValue) key() *Key {
	return &kv.Key
}

type String struct {
	Multiline bool
	Literal   bool
	Value     string
}

type Integer struct {
	Value int64
}
