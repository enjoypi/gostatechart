package gostatechart

type Reaction func(e Event) Event

type State interface {
	Begin(context interface{}, event Event) Event
	End(event Event) Event // implemented in SimpleState
	GetEvent() Event       // implemented in SimpleState
	GetTransitions() Transitions
	InitialChildState() State // implemented in SimpleState

	// implemented in SimpleState
	initiate(machine *StateMachine, state State, context interface{}, event Event) Event
	react(event Event) Event
	terminate(event Event)
}
