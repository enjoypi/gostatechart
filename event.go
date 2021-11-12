package gostatechart

import "reflect"

type Event interface{}

var (
	typeOfString = reflect.TypeOf("")
)

func TypeOf(event Event) reflect.Type {
	switch event.(type) {
	case string:
		return typeOfString
	default:
		return reflect.TypeOf(event)
	}
}
