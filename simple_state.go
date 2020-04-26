package gostatechart

type SimpleState struct {
	StateMachine
	reactions map[Event]Reaction
}

func NewSimpleState(sm StateMachine, initialSubState State) *SimpleState {
	return &SimpleState{
		StateMachine: sm,
	}
}

func (state *SimpleState) End(event Event) Event {
	return nil
}

func (state *SimpleState) GetEvent() Event {
	return nil
}

func (state *SimpleState) PostEvent(event Event) {
	state.StateMachine.PostEvent(event)
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

func (state *SimpleState) RegisterTransition(event Event, next State) error {
	return nil
}
