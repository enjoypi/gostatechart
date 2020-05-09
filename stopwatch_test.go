package gostatechart_test

import (
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

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
	s.T = context.(*testing.T)
	s.T.Logf("%T Begin %#v", s, event)
	return nil
}

func (s *Active) End(event sc.Event) sc.Event {
	s.T.Logf("%T End %#v", s, event)
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
	s.T = context.(*testing.T)
	s.T.Logf("%T Begin %#v", s, event)
	return nil
}

func (s *Stopped) End(event sc.Event) sc.Event {
	s.T.Logf("%T End %#v", s, event)
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
	s.T = context.(*testing.T)
	s.T.Logf("%T Begin %#v", s, event)
	return nil
}

func (s *Running) End(event sc.Event) sc.Event {
	s.T.Logf("%T End %#v", s, event)
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
