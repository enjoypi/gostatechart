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
	PostEvent(event Event)
	React(event Event) Event
	RegisterReaction(event Event, reaction Reaction) error
}

type SimpleState struct {
	outermostContext StateMachine
	reactions        map[Event]Reaction
}

func NewSimpleState(outermostContext StateMachine) *SimpleState {
	return &SimpleState{
		outermostContext: outermostContext,
	}
}

func (state *SimpleState) End(event Event) Event {
	return nil
}

func (state *SimpleState) GetEvent() Event {
	return nil
}

func (state *SimpleState) PostEvent(event Event) {
	state.outermostContext.PostEvent(event)
}

func (state *SimpleState) React(event Event) Event {
	reaction, ok := state.reactions[event]
	if !ok {
		return nil
	}

	return reaction(event)
}

func (state *SimpleState) RegisterReaction(event Event, reaction Reaction) error {
	if _, ok := state.reactions[event]; ok {
		panic("event already exists")
	}
	state.reactions[event] = reaction
	return nil
}
