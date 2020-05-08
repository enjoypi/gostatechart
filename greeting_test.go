package gostatechart_test

import (
	"reflect"
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

type Greeting struct {
	sc.SimpleState
}

type eGreetingBegun struct {
}

// entry
func (s *Greeting) Begin(context interface{}, event sc.Event) sc.Event {
	return &eGreetingBegun{}
}

func (s *Greeting) GetTransitions() sc.Transitions {
	return nil
}

func TestGreeting(t *testing.T) {
	sm := sc.NewStateMachine(&Greeting{}, t)
	defer func() {
		sm.Close()
	}()
	require.NoError(t, sm.Initiate())
	require.IsType(t, (*Greeting)(nil), sm.CurrentState())
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
