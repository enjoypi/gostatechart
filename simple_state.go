package gostatechart

import "reflect"

//simple_state<>::post_event()
//simple_state<>::clear_shallow_history<>()
//simple_state<>::clear_deep_history<>()
//simple_state<>::outermost_context()
//simple_state<>::context<>()
//simple_state<>::state_cast<>()
//simple_state<>::state_downcast<>()
//simple_state<>::state_begin()
//simple_state<>::state_end()

type SimpleState struct {
	parent *StateMachine
	*StateMachine
	reactions map[Event]Reaction
}

// End to override
func (state *SimpleState) End(event Event) Event {
	return nil
}

// GetEvent to override
func (state *SimpleState) GetEvent() Event {
	return nil
}

// InitialChildState to override
func (state *SimpleState) InitialChildState() State {
	return nil
}

func (state *SimpleState) close() {
	machine := state.StateMachine
	if machine != nil {
		machine.Close()
	}
}

func typename(v interface{}) string {
	if v == nil {
		return "nil"
	}
	return reflect.TypeOf(v).Elem().Name()
}

func (state *SimpleState) initiate(parent *StateMachine, context interface{}, event Event) Event {
	println("initiate ", typename(state), typename(event))
	child := state.InitialChildState()
	if child != nil {
		machine := NewStateMachine(child, context)
		if err := machine.Initiate(); err != nil {
			return err
		}
		state.StateMachine = machine
	}
	state.parent = parent
	return nil
}

func (state *SimpleState) react(event Event) (ret Event) {
	reaction, ok := state.reactions[event]
	if ok {
		ret = reaction(event)
	}

	machine := state.StateMachine
	if machine != nil {
		machine.ProcessEvent(event)
	}
	return ret
}

func (state *SimpleState) RegisterReaction(event Event, reaction Reaction, nextState State) error {
	if _, ok := state.reactions[event]; ok {
		panic("event already exists")
	}
	state.reactions[event] = reaction

	return nil
}
