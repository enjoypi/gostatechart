package gostatechart

type Reaction func(e Event) Event

//simple_state<>::post_event()
//simple_state<>::clear_shallow_history<>()
//simple_state<>::clear_deep_history<>()
//simple_state<>::outermost_context()
//simple_state<>::context<>()
//simple_state<>::state_cast<>()
//simple_state<>::state_downcast<>()
//simple_state<>::state_begin()
//simple_state<>::state_end()

type State interface {
	Begin(context interface{}, event Event) Event

	// Default implement in SimpleState
	End(event Event) Event
	GetEvent() Event
	GetTransitions() Transitions
	React(event Event) Event
	RegisterReaction(event Event, reaction Reaction) error
}
