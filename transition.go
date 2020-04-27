package gostatechart

import (
	"fmt"
	"reflect"
)

type Transitions map[reflect.Type]reflect.Type

func NewTranstions() Transitions {
	return make(Transitions)
}

func (trans Transitions) RegisterTransition(event Event, next State) {
	eventType := reflect.TypeOf(event)
	nextType := reflect.TypeOf(next)

	n, ok := trans[eventType]
	if ok && n != nextType {
		panic(fmt.Errorf("duplicate event %s", eventType.Elem().Name()))
	}

	trans[eventType] = nextType
}
