package gostatechart

import (
	"fmt"
	"log"
)

type transitions map[Event]State

type StateMachine struct {
	*Factory
	initialState State

	context            interface{}
	currentState       State
	currentTransitions transitions
	events             []Event
	allTransitions     map[string]transitions
}

func NewStateMachine(factory *Factory, initialState State) *StateMachine {
	if factory == nil {
		factory = DefaultFactory
	}
	return &StateMachine{
		Factory:        factory,
		initialState:   initialState,
		events:         make([]Event, 0, 16),
		allTransitions: make(map[string]transitions),
	}
}

func (sm *StateMachine) Initiate(context interface{}) error {
	if sm.currentState != nil {
		return fmt.Errorf("already running")
	}

	sm.context = context
	sm.transit(sm.initialState, nil)

	return nil
}

func (sm *StateMachine) PostEvent(e Event) {
	sm.events = append(sm.events, e)
}

func (sm *StateMachine) ProcessEvent(e Event) {
	next, ok := sm.currentTransitions[e]
	if ok {
		sm.transit(next, e)
	} else {
		log.Printf("%s React %s", typename(sm.currentState), typename(e))
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
	name := typename(state)
	trans, ok := sm.allTransitions[name]
	if !ok {
		trans = make(transitions)
	}

	n, ok := trans[event]
	if ok && n != next {
		return fmt.Errorf("")
	}

	trans[event] = next
	sm.allTransitions[name] = trans
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

func (sm *StateMachine) transit(s State, event Event) {
	if sm.currentState != nil {
		log.Printf("%s End %s", typename(sm.currentState), typename(event))
		if e := sm.currentState.End(event); e != nil {
			sm.PostEvent(e)
		}
	}

	state := sm.Factory.NewState(s)
	sm.currentState = state
	sm.currentTransitions = sm.allTransitions[typename(state)]

	currentName := typename(state)
	for e, s := range sm.currentTransitions {
		log.Printf("%s To %s When %s", currentName, typename(s), typename(e))
	}

	log.Printf("%s Begin %s", typename(state), typename(event))
	if e := state.Begin(sm.context, nil); e != nil {
		sm.PostEvent(e)
	}
}
