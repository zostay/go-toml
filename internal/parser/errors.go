package parser

import "fmt"

type Key []string

type DecodeError struct {
	Highlight []byte
	Message   string
	Key       Key // optional
}

func (de *DecodeError) Error() string {
	return de.Message
}

func NewDecodeError(highlight []byte, format string, args ...interface{}) error {
	return &DecodeError{
		Highlight: highlight,
		Message:   fmt.Errorf(format, args...).Error(),
	}
}
