package gostatechart

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type EvStartStop struct {
}

type EvReset struct {
}

type Active struct {
	SimpleState
}

func (s *Active) Begin(context interface{}, event Event) Event {
	return nil
}

func (s *Active) GetTransitions() Transitions {
	trans := NewTranstions()
	trans.RegisterTransition((*EvReset)(nil), (*Active)(nil))
	return trans
}

type Stopped struct {
	SimpleState
}

func (s *Stopped) GetTransitions() Transitions {
	trans := NewTranstions()
	trans.RegisterTransition((*EvStartStop)(nil), (*Running)(nil))
	return trans
}

func (s *Stopped) Begin(context interface{}, event Event) Event {
	return nil
}

type Running struct {
	SimpleState
}

func (s *Running) GetTransitions() Transitions {
	trans := NewTranstions()
	trans.RegisterTransition((*EvStartStop)(nil), (*Stopped)(nil))
	return trans
}

func (s *Running) Begin(context interface{}, event Event) Event {
	return nil
}

func TestStopWatch(t *testing.T) {
	stopWatch := NewStateMachine((*Stopped)(nil))
	require.NoError(t, stopWatch.Initiate(nil))
	require.IsType(t, (*Stopped)(nil), stopWatch.currentState)
	stopWatch.ProcessEvent((*EvStartStop)(nil))
	require.IsType(t, (*Running)(nil), stopWatch.currentState)
	stopWatch.ProcessEvent((*EvStartStop)(nil))
	require.IsType(t, (*Stopped)(nil), stopWatch.currentState)
}
