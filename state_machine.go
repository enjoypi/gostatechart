package gostatechart

import (
	"fmt"
	"reflect"
)

//namespace boost
//{
//namespace statechart
//{
//  template<
//    class MostDerived,
//    class InitialState,
//    class Allocator = std::allocator< void >,
//    class ExceptionTranslator = null_exception_translator >
//  class state_machine : noncopyable
//  {
//    public:
//      typedef MostDerived outermost_context_type;
//
//      void initiate();
//      void terminate();
//      bool terminated() const;
//
//      void process_event( const event_base & );
//
//      template< class Target >
//      Target state_cast() const;
//      template< class Target >
//      Target state_downcast() const;
//
//      // a model of the StateBase concept
//      typedef implementation-defined state_base_type;
//      // a model of the standard Forward Iterator concept
//      typedef implementation-defined state_iterator;
//
//      state_iterator state_begin() const;
//      state_iterator state_end() const;
//
//      void unconsumed_event( const event_base & ) {}
//
//    protected:
//      state_machine();
//      ~state_machine();
//
//      void post_event(
//        const intrusive_ptr< const event_base > & );
//      void post_event( const event_base & );
//
//      const event_base * triggering_event() const;
//  };
//}
//}

type StateMachine struct {
	initialState reflect.Type
	context      interface{}
	events       []Event
	doubleEvents []Event

	currentState       State
	currentTransitions Transitions
	parent             *StateMachine
}

func NewStateMachine(initialState State, context interface{}) *StateMachine {
	return &StateMachine{
		initialState: reflect.TypeOf(initialState),
		context:      context,
		events:       make([]Event, 0, 16),
		doubleEvents: make([]Event, 0, 16),
	}
}

func (machine *StateMachine) CurrentState() State {
	return machine.currentState
}

func (machine *StateMachine) Initiate(event Event) error {
	if machine.currentState != nil {
		return fmt.Errorf("already running")
	}

	machine.transit(machine.initialState, event)
	machine.run()
	return nil
}

func (machine *StateMachine) PostEvent(e Event) {
	m := machine.outermost()
	m.events = append(m.events, e)
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

	currentState := machine.currentState
	if currentState != nil {
		for e := currentState.GetEvent(); e != nil; e = currentState.GetEvent() {
			machine.PostEvent(e)
		}
	}

	machine.run()
}

func (machine *StateMachine) Terminate(event Event) {
	if machine == nil {
		return
	}
	machine.transit(nil, event)
}

func (machine *StateMachine) outermost() *StateMachine {
	if machine == nil {
		return nil
	}

	if machine.parent == nil {
		return machine
	}
	return machine.parent.outermost()
}

func (machine *StateMachine) run() {
	if len(machine.events) <= 0 {
		return
	}

	events := machine.events
	machine.events = machine.doubleEvents

	for _, e := range events {
		machine.ProcessEvent(e)
	}
	machine.doubleEvents = events[:0]
}

func (machine *StateMachine) transit(stateType reflect.Type, event Event) {
	if machine.currentState != nil {
		machine.currentState.terminate(event)
		if e := machine.currentState.End(event); e != nil {
			machine.PostEvent(e)
		}
		machine.currentState = nil
	}

	if stateType == nil {
		return
	}

	nextState := reflect.New(stateType.Elem()).Interface().(State)
	machine.currentState = nextState
	machine.currentTransitions = nextState.GetTransitions()

	if e := machine.currentState.initiate(machine, nextState, machine.context, event); e != nil {
		machine.PostEvent(e)
	}

	if e := machine.currentState.Begin(machine.context, event); e != nil {
		machine.PostEvent(e)
	}
}
