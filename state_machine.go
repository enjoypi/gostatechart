package gostatechart

import (
	"fmt"
	"reflect"
)

type StateMachine struct {
	initialState reflect.Type

	context            interface{}
	currentState       State
	currentTransitions Transitions
	events             []Event
}

func NewStateMachine(initialState State) *StateMachine {
	return &StateMachine{
		initialState: reflect.TypeOf(initialState),
		events:       make([]Event, 0, 16),
	}
}

func (sm *StateMachine) Close() {
	sm.currentState.End((*EvClose)(nil))
}

func (sm *StateMachine) Initiate(context interface{}) error {
	if sm.currentState != nil {
		return fmt.Errorf("already running")
	}

	sm.context = context
	sm.transit(sm.initialState, (*EvInit)(nil))

	return nil
}

func (sm *StateMachine) PostEvent(e Event) {
	sm.events = append(sm.events, e)
}

func (sm *StateMachine) ProcessEvent(e Event) {
	next, ok := sm.currentTransitions[reflect.TypeOf(e)]
	if ok {
		sm.transit(next, e)
	} else {
		ne := sm.currentState.React(e)
		if ne != nil {
			sm.PostEvent(ne)
		}
	}

	for e := sm.currentState.GetEvent(); e != nil; e = sm.currentState.GetEvent() {
		sm.PostEvent(e)
	}
}

func (sm *StateMachine) Run() {
	if len(sm.events) <= 0 {
		return
	}

	events := sm.events
	sm.events = make([]Event, 0, 16)

	for _, e := range events {
		sm.ProcessEvent(e)
	}
}

func (sm *StateMachine) transit(s reflect.Type, event Event) {
	if sm.currentState != nil {
		if e := sm.currentState.End(event); e != nil {
			sm.PostEvent(e)
		}
	}

	state := reflect.New(s.Elem()).Interface().(State)
	sm.currentState = state
	sm.currentTransitions = state.GetTransitions()

	if e := state.Begin(sm.context, nil); e != nil {
		sm.PostEvent(e)
	}
}
