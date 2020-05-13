package gostatechart

import (
	"fmt"
	"reflect"
)

type StateMachine struct {
	initialState reflect.Type
	context      interface{}
	events       []Event
	doubleEvents []Event

	currentState       State
	currentTransitions Transitions
}

func NewStateMachine(initialState State, context interface{}) *StateMachine {
	return &StateMachine{
		initialState: reflect.TypeOf(initialState),
		context:      context,
		events:       make([]Event, 0, 16),
		doubleEvents: make([]Event, 0, 16),
	}
}

func (machine *StateMachine) Close(event Event) {
	machine.migrate(nil, event)
}

func (machine *StateMachine) CurrentState() State {
	return machine.currentState
}

func (machine *StateMachine) Initiate(event Event) error {
	if machine.currentState != nil {
		return fmt.Errorf("already running")
	}

	machine.migrate(machine.initialState, event)
	machine.run()
	return nil
}

func (machine *StateMachine) PostEvent(e Event) {
	machine.events = append(machine.events, e)
}

func (machine *StateMachine) ProcessEvent(e Event) {
	ne := machine.currentState.react(e)
	if ne != nil {
		machine.PostEvent(ne)
	}

	next, ok := machine.currentTransitions[reflect.TypeOf(e)]
	if ok {
		machine.migrate(next, e)
	}

	for e := machine.currentState.GetEvent(); e != nil; e = machine.currentState.GetEvent() {
		machine.PostEvent(e)
	}

	machine.run()
}

func (machine *StateMachine) run() {
	if len(machine.events) <= 0 {
		return
	}

	events := machine.events
	machine.events = machine.doubleEvents

	for _, e := range events {
		machine.ProcessEvent(e)
	}
	machine.doubleEvents = events[:0]
}

func (machine *StateMachine) migrate(stateType reflect.Type, event Event) {
	if machine.currentState != nil {
		machine.currentState.close(event)
		if e := machine.currentState.End(event); e != nil {
			machine.PostEvent(e)
		}
	}

	if stateType == nil {
		return
	}

	nextState := reflect.New(stateType.Elem()).Interface().(State)
	machine.currentState = nextState
	machine.currentTransitions = nextState.GetTransitions()

	if e := machine.currentState.initiate(machine, nextState, machine.context, event); e != nil {
		machine.PostEvent(e)
	}

	if e := machine.currentState.Begin(machine.context, event); e != nil {
		machine.PostEvent(e)
	}
}
