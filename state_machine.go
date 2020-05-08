package gostatechart

import (
	"fmt"
	"reflect"
)

type StateMachine struct {
	initialState reflect.Type
	context      interface{}
	events       []Event

	currentState       State
	currentTransitions Transitions
}

func NewStateMachine(initialState State, context interface{}) *StateMachine {
	return &StateMachine{
		initialState: reflect.TypeOf(initialState),
		context:      context,
		events:       make([]Event, 0, 16),
	}
}

func (machine *StateMachine) Close() {
	machine.transit(nil, nil)
}

func (machine *StateMachine) CurrentState() State {
	return machine.currentState
}

func (machine *StateMachine) Initiate() error {
	if machine.currentState != nil {
		return fmt.Errorf("already running")
	}

	machine.transit(machine.initialState, nil)
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
		machine.transit(next, e)
	}

	for e := machine.currentState.GetEvent(); e != nil; e = machine.currentState.GetEvent() {
		machine.PostEvent(e)
	}
}

func (machine *StateMachine) Run() {
	if len(machine.events) <= 0 {
		return
	}

	events := machine.events
	machine.events = make([]Event, 0, 16)

	for _, e := range events {
		machine.ProcessEvent(e)
	}
}

func (machine *StateMachine) transit(stateType reflect.Type, event Event) {
	if machine.currentState != nil {
		if e := machine.currentState.End(event); e != nil {
			machine.PostEvent(e)
		}
		machine.currentState.close()
	}

	if stateType == nil {
		return
	}

	nextState := reflect.New(stateType.Elem()).Interface().(State)
	machine.currentState = nextState
	machine.currentTransitions = nextState.GetTransitions()

	if e := machine.currentState.initiate(machine, machine.context, event); e != nil {
		machine.PostEvent(e)
	}

	if e := machine.currentState.Begin(machine.context, event); e != nil {
		machine.PostEvent(e)
	}
}
