package gostatechart

import "fmt"

type Event interface{}

type EventError struct {
	error
}

func NewEventError(err error) Event {
	return &EventError{error: err}
}

func NewEventErrorf(format string, a ...interface{}) Event {
	return NewEventError(fmt.Errorf(format, a...))
}
