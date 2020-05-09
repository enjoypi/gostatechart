package gostatechart_test

import (
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

func logf(t *testing.T, format string, args ...interface{}) {
	if t != nil {
		t.Logf(format, args...)
	}
}

type EvStartStop struct {
}

type EvReset struct {
}

type EvClose struct {
}

type Active struct {
	sc.SimpleState
	*testing.T
}

func (s *Active) Begin(context interface{}, event sc.Event) sc.Event {
	if context != nil {
		s.T = context.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	return nil
}

func (s *Active) End(event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Active) GetTransitions() sc.Transitions {
	trans := sc.NewTranstions()
	trans.RegisterTransition((*EvReset)(nil), (*Active)(nil))
	return trans
}

func (s *Active) InitialChildState() sc.State {
	return (*Stopped)(nil)
}

type Stopped struct {
	sc.SimpleState
	*testing.T
}

func (s *Stopped) Begin(context interface{}, event sc.Event) sc.Event {
	if context != nil {
		s.T = context.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	return nil
}

func (s *Stopped) End(event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Stopped) GetTransitions() sc.Transitions {
	trans := sc.NewTranstions()
	trans.RegisterTransition((*EvStartStop)(nil), (*Running)(nil))
	return trans
}

type Running struct {
	sc.SimpleState
	*testing.T
}

func (s *Running) Begin(context interface{}, event sc.Event) sc.Event {
	if context != nil {
		s.T = context.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	return nil
}

func (s *Running) End(event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Running) GetTransitions() sc.Transitions {
	trans := sc.NewTranstions()
	trans.RegisterTransition((*EvStartStop)(nil), (*Stopped)(nil))
	return trans
}

func TestStopWatch(t *testing.T) {
	stopWatch := sc.NewStateMachine((*Active)(nil), t)
	require.Nil(t, stopWatch.CurrentState())
	require.NoError(t, stopWatch.Initiate(nil))
	defer func() {
		t.Log("Close")
		stopWatch.Close(&EvClose{})
	}()

	active := stopWatch.CurrentState().(*Active)
	require.NotNil(t, active)
	require.IsType(t, (*Active)(nil), active)

	require.NotNil(t, active.CurrentState())
	require.IsType(t, (*Stopped)(nil), active.CurrentState())

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(&EvStartStop{})
	require.IsType(t, (*Running)(nil), active.CurrentState())

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(&EvStartStop{})
	require.IsType(t, (*Stopped)(nil), active.CurrentState())

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(&EvStartStop{})
	require.IsType(t, (*Running)(nil), active.CurrentState())

	t.Logf("EvReset")
	stopWatch.ProcessEvent(&EvReset{})
	require.IsType(t, (*Active)(nil), stopWatch.CurrentState())
	require.NotEqual(t, active, stopWatch.CurrentState())
}

func BenchmarkStopWatch(b *testing.B) {
	stopWatch := sc.NewStateMachine((*Active)(nil), nil)
	stopWatch.CurrentState()
	_ = stopWatch.Initiate(nil)
	defer func() {
		stopWatch.Close(&EvClose{})
	}()

	e := &EvStartStop{}
	for i := 0; i < b.N; i++ {
		stopWatch.ProcessEvent(e)
	}
}
