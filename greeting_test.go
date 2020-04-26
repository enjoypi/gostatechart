package gostatechart

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

type Greeting struct {
	SimpleState
}

type eGreetingBegun struct {
}

// entry
func (s *Greeting) Begin(context interface{}, event Event) Event {
	return &eGreetingBegun{}
}

// exit
func (s *Greeting) End(event Event) Event {
	return nil
}

func TestGreeting(t *testing.T) {
	sm := NewStateMachine(&Greeting{})
	defer func() {
		sm.Close()
	}()
	require.NoError(t, sm.Initiate(t))
	assert.IsType(t, reflect.TypeOf((*Greeting)(nil)), sm.initialState)
}

func BenchmarkTypeOf(b *testing.B) {
	g := &Greeting{}
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(g)
	}
}

func BenchmarkTypeOfName(b *testing.B) {
	g := &Greeting{}
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(g).Elem().Name()
	}
}
