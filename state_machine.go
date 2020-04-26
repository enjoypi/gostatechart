package gostatechart

import (
	"fmt"
	"reflect"
)

type transitions map[reflect.Type]reflect.Type

type StateMachine struct {
	initialState reflect.Type

	context            interface{}
	currentState       State
	currentTransitions transitions
	events             []Event
	allTransitions     map[reflect.Type]transitions
}

func NewStateMachine(initialState State) *StateMachine {
	return &StateMachine{
		initialState:   reflect.TypeOf(initialState),
		events:         make([]Event, 0, 16),
		allTransitions: make(map[reflect.Type]transitions),
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

func (sm *StateMachine) RegisterTransition(state State, event Event, next State) error {
	curType := reflect.TypeOf(state)
	trans, ok := sm.allTransitions[curType]
	if !ok {
		trans = make(transitions)
	}

	eventType := reflect.TypeOf(event)
	nextType := reflect.TypeOf(next)

	n, ok := trans[eventType]
	if ok && n != nextType {
		return fmt.Errorf("the transitions of %s on %s is exists", curType.Elem().Name(), eventType.Elem().Name())
	}

	trans[eventType] = nextType
	sm.allTransitions[curType] = trans
	return nil
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

	// msg0 := reflect.New(typ.Elem()).Interface().(proto.Message)
	//state := sm.Factory.NewState(s)
	state := reflect.New(s.Elem()).Interface().(State)
	sm.currentState = state
	sm.currentTransitions = sm.allTransitions[s]

	if e := state.Begin(sm.context, nil); e != nil {
		sm.PostEvent(e)
	}
}
