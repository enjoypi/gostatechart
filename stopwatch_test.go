package gostatechart_test

import (
	"context"
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

func (s *Active) Begin(ctx context.Context, event sc.Event) sc.Event {
	t := ctx.Value("testing.T")
	if t != nil {
		s.T = t.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	s.RegisterReaction((*EvSth)(nil), s.OnSth)
	return nil
}

func (s *Active) End(ctx context.Context, event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Active) GetTransitions() sc.Transitions {
	return activeTrans
}

func (s *Active) InitialChildState() sc.State {
	return (*Stopped)(nil)
}

func (s *Active) OnSth(ctx context.Context, event sc.Event, args ...interface{}) sc.Event {
	logf(s.T, "%T OnSth %#v", s, event)
	return nil
}

type Stopped struct {
	sc.SimpleState
	*testing.T
}

func (s *Stopped) Begin(ctx context.Context, event sc.Event) sc.Event {
	t := ctx.Value("testing.T")
	if t != nil {
		s.T = t.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	return nil
}

func (s *Stopped) End(ctx context.Context, event sc.Event) sc.Event {
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

func (s *Running) Begin(ctx context.Context, event sc.Event) sc.Event {
	t := ctx.Value("testing.T")
	if t != nil {
		s.T = t.(*testing.T)
	}
	logf(s.T, "%T Begin %#v", s, event)
	return nil
}

func (s *Running) End(ctx context.Context, event sc.Event) sc.Event {
	logf(s.T, "%T End %#v", s, event)
	return nil
}

func (s *Running) GetTransitions() sc.Transitions {
	return runningTrans
}

func TestStopWatch(t *testing.T) {
	ctx := context.WithValue(context.Background(), "testing.T", t)
	stopWatch := sc.NewStateMachine((*Active)(nil), ctx)
	require.Nil(t, stopWatch.CurrentState())
	require.NoError(t, stopWatch.Initiate(nil))
	defer func() {
		t.Log("Close")
		stopWatch.Terminate(&EvClose{})
	}()

	active := stopWatch.CurrentState().(*Active)
	require.NotNil(t, active)
	require.IsType(t, (*Active)(nil), active)

	require.NotNil(t, active.CurrentState())
	require.IsType(t, (*Stopped)(nil), active.CurrentState())

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(ctx, &EvStartStop{})
	require.IsType(t, (*Running)(nil), active.CurrentState())
	stopWatch.Run(nil)

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(ctx, &EvStartStop{})
	require.IsType(t, (*Stopped)(nil), active.CurrentState())
	stopWatch.Run(nil)

	t.Logf("EvStartStop")
	stopWatch.ProcessEvent(ctx, &EvStartStop{})
	require.IsType(t, (*Running)(nil), active.CurrentState())
	stopWatch.Run(nil)

	t.Logf("EvReset")
	stopWatch.ProcessEvent(ctx, &EvReset{})
	require.IsType(t, (*Active)(nil), stopWatch.CurrentState())
	require.NotEqual(t, active, stopWatch.CurrentState())
	stopWatch.Run(nil)

	t.Logf("EvSth")
	stopWatch.ProcessEvent(ctx, &EvSth{})
	stopWatch.Run(nil)
}

func BenchmarkTransmit(b *testing.B) {
	ctx := context.WithValue(context.Background(), "testing.B", b)
	stopWatch := sc.NewStateMachine((*Active)(nil), ctx)
	stopWatch.CurrentState()
	_ = stopWatch.Initiate(nil)
	defer func() {
		stopWatch.Terminate(&EvClose{})
	}()

	e := &EvStartStop{}
	for i := 0; i < b.N; i++ {
		stopWatch.ProcessEvent(ctx, e)
	}
}

func BenchmarkProcessEvent(b *testing.B) {
	ctx := context.WithValue(context.Background(), "testing.B", b)
	stopWatch := sc.NewStateMachine((*Active)(nil), nil)
	stopWatch.CurrentState()
	_ = stopWatch.Initiate(nil)
	defer func() {
		stopWatch.Terminate(&EvClose{})
	}()

	e := &EvSth{}
	for i := 0; i < b.N; i++ {
		stopWatch.ProcessEvent(ctx, e)
	}
}
