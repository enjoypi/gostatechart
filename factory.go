package gostatechart

import (
	"fmt"
	"reflect"
)

type NewState func() State

type Factory struct {
	constructors map[string]NewState
}

var DefaultFactory = &Factory{
	constructors: make(map[string]NewState),
}

func New(name string) State {
	return DefaultFactory.New(name)
}

func RegisterState(state State, newState NewState) {
	if err := DefaultFactory.RegisterState(state, newState); err != nil {
		// duplicate state will cause complicated error, it must be exposed early
		panic(err)
	}
}

func (f *Factory) New(name string) State {
	if newState, ok := f.constructors[name]; ok {
		return newState()
	}
	return nil
}

func (f *Factory) NewState(state State) State {
	return New(typename(state))
}

func (f *Factory) RegisterState(state State, newState NewState) error {
	name := typename(state)
	if _, ok := f.constructors[name]; ok {
		return fmt.Errorf("%s ", name)
	}

	f.constructors[name] = newState
	return nil
}

func typename(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).Elem().Name()
}
