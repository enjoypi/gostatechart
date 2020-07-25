package gostatechart

import (
	"fmt"
	"reflect"
	"time"
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
	events       chan Event

	currentState       State
	currentTransitions Transitions
	parent             *StateMachine
}

func NewStateMachine(initialState State, context interface{}) *StateMachine {
	if initialState == nil {
		return nil
	}
	return &StateMachine{
		initialState: reflect.TypeOf(initialState),
		context:      context,
		events:       make(chan Event, 32),
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
	return nil
}

func (machine *StateMachine) Parent() *StateMachine {
	if machine.parent == nil {
		return machine
	}
	return machine.parent.Parent()
}

func (machine *StateMachine) PostEvent(e Event) {
	if e == nil {
		return
	}

	if machine.parent == nil {
		machine.events <- e
		return
	}

	machine.parent.PostEvent(e)
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
}

func (machine *StateMachine) Run(exitChan chan int) {
	for machine.currentState != nil {
		currentState := machine.currentState
		if currentState != nil {
			if e := currentState.GetEvent(); e != nil {
				machine.PostEvent(e)
			}
		}

		select {
		case e := <-machine.events:
			machine.ProcessEvent(e)
		case <-time.After(10 * time.Millisecond):
			continue
		case <-exitChan:
			machine.Terminate(nil)
			return
		}
	}
}

func (machine *StateMachine) Terminate(event Event) {
	machine.transit(nil, event)
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
