package gostatechart

type Reaction func(e Event) Event

type State interface {
	Begin(context interface{}, event Event) Event
	End(event Event) Event       // implemented in SimpleState
	GetEvent() Event             // implemented in SimpleState
	GetTransitions() Transitions // implemented in SimpleState
	InitialChildState() State    // implemented in SimpleState
	React(event Event) Event

	// implemented in SimpleState
	initiate(machine *StateMachine, state State, context interface{}, event Event) Event
	// terminate my machine
	terminate(event Event)
}
