package gostatechart_test

import (
	"context"
	"reflect"
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

type Greeting struct {
	sc.SimpleState
	*testing.T
}

type EvGreetingBegun struct {
}

// entry
func (s *Greeting) Begin(ctx context.Context, event sc.Event) sc.Event {
	s.RegisterReaction((*EvGreetingBegun)(nil), s.OnBegun)
	s.T = ctx.Value("testing.T").(*testing.T)
	return &EvGreetingBegun{}
}

func (s *Greeting) GetTransitions() sc.Transitions {
	return nil
}

func (s *Greeting) OnBegun(ctx context.Context, event sc.Event) sc.Event {
	s.T.Logf("%#v", event)
	return nil
}

func TestGreeting(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing.T", t)
	sm := sc.NewStateMachine(&Greeting{}, ctx)
	defer func() {
		sm.Terminate(nil)
	}()
	require.NoError(t, sm.Initiate(nil))
	require.IsType(t, (*Greeting)(nil), sm.CurrentState())
}

func BenchmarkTypeOf(b *testing.B) {
	g := &Greeting{}
	for i := 0; i < b.N; i++ {
		_ = reflect.TypeOf(g)
	}
}

func BenchmarkTypeOfName(b *testing.B) {
	g := &Greeting{}
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(g).Elem().Name()
	}
}
