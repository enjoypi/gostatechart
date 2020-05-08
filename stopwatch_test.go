package gostatechart_test

import (
	"reflect"
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

type EvStartStop struct {
}

type EvReset struct {
}

type Active struct {
	sc.SimpleState
	*testing.T
}

func (s *Active) Begin(context interface{}, event sc.Event) sc.Event {
	t := context.(*testing.T)
	s.T = t
	t.Logf("%s Begin %s", typename(s), typename(event))
	return nil
}

func (s *Active) End(event sc.Event) sc.Event {
	s.T.Logf("%s End %s", typename(s), typename(event))
	return nil
}

func (s *Active) GetTransitions() sc.Transitions {
	trans := sc.NewTranstions()
	trans.RegisterTransition((*EvReset)(nil), (*Active)(nil))
	return trans
}

func (s *Active) InitialChildState() sc.State {
	s.T.Logf("%s InitialChildState", typename(s))
	return (*Stopped)(nil)
}

type Stopped struct {
	sc.SimpleState
	*testing.T
}

func (s *Stopped) Begin(context interface{}, event sc.Event) sc.Event {
	t := context.(*testing.T)
	s.T = t
	t.Logf("%s Begin %s", typename(s), typename(event))
	return nil
}

func (s *Stopped) End(event sc.Event) sc.Event {
	s.T.Logf("%s End %s", typename(s), typename(event))
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
	t := context.(*testing.T)
	s.T = t
	t.Logf("%s Begin %s", typename(s), typename(event))
	return nil
}

func (s *Running) End(event sc.Event) sc.Event {
	s.T.Logf("%s End %s", typename(s), typename(event))
	return nil
}

func (s *Running) GetTransitions() sc.Transitions {
	trans := sc.NewTranstions()
	trans.RegisterTransition((*EvStartStop)(nil), (*Stopped)(nil))
	return trans
}

func typename(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).Elem().Name()
}

func TestStopWatch(t *testing.T) {
	stopWatch := sc.NewStateMachine((*Active)(nil), t)
	require.Nil(t, nil, stopWatch.CurrentState())
	require.NoError(t, stopWatch.Initiate())
	stopWatch.Run()
	defer stopWatch.Close()

	active := stopWatch.CurrentState().(*Active)
	require.IsType(t, (*Active)(nil), active)

	require.IsType(t, (*Stopped)(nil), active.CurrentState())
	stopWatch.ProcessEvent((*EvStartStop)(nil))
	stopWatch.Run()

	require.IsType(t, (*Running)(nil), active.CurrentState())
	stopWatch.ProcessEvent((*EvStartStop)(nil))
	stopWatch.Run()

	require.IsType(t, (*Stopped)(nil), active.CurrentState())
}
