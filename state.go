package gostatechart

import "context"

type Reaction func(ctx context.Context, e Event, args ...interface{}) Event

type State interface {
	Begin(ctx context.Context, event Event) Event
	End(ctx context.Context, event Event) Event // implemented in SimpleState
	GetEvent() Event                            // implemented in SimpleState
	GetTransitions() Transitions                // implemented in SimpleState
	InitialChildState() State                   // implemented in SimpleState
	React(ctx context.Context, event Event, args ...interface{}) Event

	// implemented in SimpleState
	initiate(ctx context.Context, machine *StateMachine, state State, event Event) Event
	// terminate my machine
	terminate(event Event)
}
