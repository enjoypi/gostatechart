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
	context   interface{}
	machine   *StateMachine
	parent    *StateMachine
	reactions map[reflect.Type]Reaction
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

func (state *SimpleState) Context() interface{} {
	return state.context
}

func (state *SimpleState) CurrentState() State {
	return state.machine.CurrentState()
}

func (state *SimpleState) RegisterReaction(event Event, reaction Reaction) error {
	if state.reactions == nil {
		state.reactions = make(map[reflect.Type]Reaction)
	}

	eventType := reflect.TypeOf(event)
	if _, ok := state.reactions[eventType]; ok {
		panic("event already exists")
	}
	state.reactions[eventType] = reaction
	return nil
}

func (state *SimpleState) close(event Event) {
	machine := state.machine
	if machine != nil {
		machine.Close(event)
	}
}

func (state *SimpleState) initiate(parent *StateMachine, self State, context interface{}, event Event) Event {
	child := self.InitialChildState()
	if child != nil {
		machine := NewStateMachine(child, context)
		if err := machine.Initiate(event); err != nil {
			return err
		}
		state.machine = machine
	}
	state.context = context
	state.parent = parent
	return nil
}

func (state *SimpleState) react(event Event) (ret Event) {
	reaction, ok := state.reactions[reflect.TypeOf(event)]
	if ok {
		ret = reaction(event)
	}

	machine := state.machine
	if machine != nil {
		machine.ProcessEvent(event)
	}
	return ret
}
