// Package document provides tools for manipulating the structure of TOML
// documents.
//
// While github.com/pelletier/go-toml provides efficient functions to transform
// TOML documents to and from usual Go types, this package allows you to create
// and modify the structure of a TOML document.
//
// Design decisions
//
// Document does not represent white space. When parsing a document from bytes,
// all white space is discarded. The only control over white space is provided
// by the document encoder in the form of general rules. One use case that is
// not covered is modifying an existing TOML document, while keeping the
// non-modified part of the document exactly the same byte-for-byte. However it
// simplifies the API and parsing significantly.
//
// It is a design goal to be able to write literal Documents and modify them
// without too much assistance. For example, instead of providing dozens of
// Create / Modify / Delete functions for all kinds of nodes, the current design
// provides allows the user to manipulate pointers and slices like any other Go
// data structure. The drawback is that the operations performed on the Document
// cannot be validated immediately. A certain amount of constraint is added in
// the form of typing, but ultimately it is the responsibility of the user to
// call Valid() after reading or before writing a Document, if they wishes to
// only deal with valid documents.
//
// While many operations would feel natural on maps, this Document structure
// actually only contains slices of elements to represent parent / children
// relationships. This allows the user to completely control the ordering of
// their document, as well as its exact shape. For example, the following valid
// documents can all be represented:
//
//  a.b.c = 42
//
//  [a]
//  b.c = 42
//
//  [a.b]
//  c = 42
//
//  [a]
//  b = { c = 42 }
//
//  [a]
//  [a.b]
//  c = 42
//
//  [a.b]
//  c = 42
//  [a]
//
// Comments are a first class object in this model. An often requested feature
// is to preserve and manipulate comments in TOML documents. By embedding them
// in the core of every node, full control is provided to the user on how they
// want to comment their document.
//
// See the Examples for examples of classic Document usages.
package document

import (
	"strconv"
)

// Document represents a TOML document.
type Document struct {
	KeyValues []*KeyValue
	Tables    []*Table

	// Optional last comment of the document.
	TrailerComment Comment
}

// GetAt traverses the document to return a pointer to the Value stored at the
// path represented by parts. Returns nil if no such document exists.
//
// Even though part/s is of type interface{}, each of them should be either a
// string or an int. If it is a string, it is interpreted as a table or
// key-value key part. If it is an integer, it is interpreted as an array index.
// -1 is used to denote the last element of the array, if it exists. Any other
// type panics.
//
// This function operates on the structure of the document. If the path is not
// explicitly defined in the document this function returns nil.
func (d Document) GetAt(part interface{}, parts ...interface{}) Value {
	// TODO
	return nil
}

// ParentOf returns the immediate parent of a given Value. Panics if the parent
// does not exist. A classic use-case is to first call GetAt to retrieve a
// specific element, then call ParentOf to get the parent and possibly reorder
// or delete the element.
func (d Document) ParentOf(v Value) Value {
	// TODO
	return nil
}

// Valid verifies that the document is fully compliant with the TOML
// specification. It returns nil if it is valid, or a list of errors otherwise.
// While this function tries to find all errors, it does not guarantee to find
// them all if at least one error is found.
func (d Document) Valid() []error {
	// TODO
	return nil
}

type Key []KeyPart

type KeyKind int

const (
	BareKey KeyKind = iota
	LiteralKey
	QuotedKey
)

// StringKey is a convenience function to generate a Key from strings. It is
// mostly useful when expressing documents as literals.
// The kind precedence of each part is BareKey > LiteralKey > QuotedKey.
func StringKey(part1 string, parts ...string) Key {
	// TODO
	return Key{}
}

type StringKind int

const (
	BasicString StringKind = 1 << iota
	LiteralString
	MultilineString
)

type KeyPart struct {
	Value string
	Kind  KeyKind
}

type KeyValue struct {
	Comment   Comment
	Commented bool

	Key   Key
	Value Value
}

// Value is an interface supported by all the terminal types of a TOML document.
type Value interface {
	isValue()
}

type String struct {
	Value string
	Kind  StringKind
}

func (s *String) isValue() {}

type Integer struct {
	V string
}

func (i *Integer) isValue() {}

func (i *Integer) Set(v int64) {
	i.V = strconv.FormatInt(v, 10)
}

func (i *Integer) FromString(v string) {
	i.V = v
}

func (i Integer) Value() int64 {
	v, err := strconv.ParseInt(i.V, 10, 64)
	if err != nil {
		panic("document should not let an invalid integer be stored")
	}
	return v
}

func (i Integer) String() string {
	return i.V
}

type Boolean struct {
	V bool
}

func (b *Boolean) isValue() {}

// TODO: Float should be the same as Integer
// TODO: different types of dates should follow the same model.

type Array struct {
	Comment   Comment
	Commented bool

	// Should each element of the array be on its own line.
	Multiline bool
	Elements  []Value
}

func (a *Array) isValue() {}

type InlineTable struct {
	Elements []*KeyValue
}

type Table struct {
	Comment   Comment
	Commented bool
	Array     bool

	Key      Key
	Elements []*KeyValue
}

// Comment is usually a member of an element of the TOML document.
type Comment struct {
	// Can be multiline. An empty value means no comment.
	Value string
	// Inline puts the comment on the same line as the element it annotates
	// (by default comments are above their element).
	Inline bool
}

// Notes / TODOs:
// - Array Table may be difficult to manipulate?
// - Use *As* functions to easily write literals (and panic if construction is
// - incorrect).
// - Use *To* functions to convert to/from standard types, returning errors if
//   needed.
// - How to discover the structure?
