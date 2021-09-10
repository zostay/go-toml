package document

import (
	"fmt"

	"github.com/pelletier/go-toml/v2/internal/ast"
	"github.com/pelletier/go-toml/v2/internal/parser"
)

// Document represents a TOML document.
// It is not guaranteed to be a valid document when manually modified.
// Has helper functions to perform usual operations.
type Document struct {
	Root []Node
}

type Node struct {
	Kind     ast.Kind
	Comment  string
	Value    interface{}
	Children []Node
	Flags    int
}

func Parse(b []byte) (Document, error) {
	d := Document{}

	p := parser.Parser{}
	p.Reset(b)
	p.Comments = true

	for p.NextExpression() {
		expr := p.Expression()
		node, err := exprToRootNode(expr)
		if err != nil {
			return d, err
		}
		d.Root = append(d.Root, node)
	}

	return d, p.Error()
}

func exprToRootNode(expr *ast.Node) (Node, error) {
	switch expr.Kind {
	case ast.Table:
		panic("todo")
	case ast.ArrayTable:
		panic("todo")
	case ast.KeyValue:
		panic("todo")
	default:
		// TODO: add error context
		return Node{}, fmt.Errorf("expression of type '%s' not allowed there", expr.Kind)
	}
}
