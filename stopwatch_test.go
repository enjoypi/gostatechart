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

func (s *Active) InitReactions() {
	s.RegisterTransition((*EvReset)(nil), (*Active)(nil))
}

type Stopped struct {
	SimpleState
}

func (s *Stopped) InitReactions() {
	s.RegisterTransition((*EvStartStop)(nil), (*Running)(nil))
}

func (s *Stopped) Begin(context interface{}, event Event) Event {
	return nil
}

type Running struct {
	SimpleState
}

func (s *Running) InitReactions() {
	s.RegisterTransition((*EvStartStop)(nil), (*Stopped)(nil))
}

func (s *Running) Begin(context interface{}, event Event) Event {
	return nil
}

func TestStopWatch(t *testing.T) {
	stopWatch := NewStateMachine((*Stopped)(nil))
	require.NoError(t, stopWatch.Initiate(nil))
}
