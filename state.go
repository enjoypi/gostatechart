package gostatechart

type Reaction func(e Event) Event

type State interface {
	Begin(context interface{}, event Event) Event
	End(event Event) Event // implemented in SimpleState
	GetEvent() Event       // implemented in SimpleState
	GetTransitions() Transitions
	InitialChildState() State // implemented in SimpleState

	close()
	initiate(machine *StateMachine, context interface{}, event Event) Event
	react(event Event) Event
}
