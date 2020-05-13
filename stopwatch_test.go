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

type EvSth struct {
}

// no alloc for benchmark
var (
	activeTrans  = sc.NewTranstions()
	stoppedTrans = sc.NewTranstions()
	runningTrans = sc.NewTranstions()
)

func init() {
	activeTrans.RegisterTransition((*EvReset)(nil), (*Active)(nil))
	stoppedTrans.RegisterTransition((*EvStartStop)(nil), (*Running)(nil))
	runningTrans.RegisterTransition((*EvStartStop)(nil), (*Stopped)(nil))
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
	s.RegisterReaction((*EvSth)(nil), s.OnSth)
	return nil
}

func (s *Active) End(event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Active) GetTransitions() sc.Transitions {
	return activeTrans
}

func (s *Active) InitialChildState() sc.State {
	return (*Stopped)(nil)
}

func (s *Active) OnSth(event sc.Event) sc.Event {
	logf(s.T, "%T OnSth %#v", s, event)
	return nil
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
	return stoppedTrans
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
	return runningTrans
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

	t.Logf("EvSth")
	stopWatch.ProcessEvent(&EvSth{})
}

func BenchmarkMigrate(b *testing.B) {
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

func BenchmarkProcessEvent(b *testing.B) {
	stopWatch := sc.NewStateMachine((*Active)(nil), nil)
	stopWatch.CurrentState()
	_ = stopWatch.Initiate(nil)
	defer func() {
		stopWatch.Close(&EvClose{})
	}()

	e := &EvSth{}
	for i := 0; i < b.N; i++ {
		stopWatch.ProcessEvent(e)
	}
}
