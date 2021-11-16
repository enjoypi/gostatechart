package gostatechart

import (
	"context"
	"reflect"
)

//namespace boost
//{
//namespace statechart
//{
//  template<
//    class MostDerived,
//    class Context,
//    class InnerInitial = unspecified,
//    history_mode historyMode = has_no_history >
//  class simple_state : implementation-defined
//  {
//    public:
//      // by default, a state has no reactions
//      typedef mpl::list<> reactions;
//
//      // see template parameters
//      template< implementation-defined-unsigned-integer-type
//        innerOrthogonalPosition >
//      struct orthogonal
//      {
//        // implementation-defined
//      };
//
//      typedef typename Context::outermost_context_type
//        outermost_context_type;
//
//      outermost_context_type & outermost_context();
//      const outermost_context_type & outermost_context() const;
//
//      template< class OtherContext >
//      OtherContext & context();
//      template< class OtherContext >
//      const OtherContext & context() const;
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
//      const event_base * triggering_event() const;
//
//      void post_event(
//        const intrusive_ptr< const event_base > & );
//      void post_event( const event_base & );
//
//      result discard_event();
//      result forward_event();
//      result defer_event();
//      template< class DestinationState >
//      result transit();
//      template<
//        class DestinationState,
//        class TransitionContext,
//        class Event >
//      result transit(
//        void ( TransitionContext::* )( const Event & ),
//        const Event & );
//      result terminate();
//
//      template<
//        class HistoryContext,
//        implementation-defined-unsigned-integer-type
//          orthogonalPosition >
//      void clear_shallow_history();
//      template<
//        class HistoryContext,
//        implementation-defined-unsigned-integer-type
//          orthogonalPosition >
//      void clear_deep_history();
//
//      static id_type static_type();
//
//      template< class CustomId >
//      static const CustomId * custom_static_type_ptr();
//
//      template< class CustomId >
//      static void custom_static_type_ptr( const CustomId * );
//
//      // see transit() or terminate() effects
//      void exit() {}
//
//    protected:
//      simple_state();
//      ~simple_state();
//  };
//}
//}

type SimpleState struct {
	machine   *StateMachine
	reactions map[reflect.Type]Reaction
	parent    *StateMachine
}

// End to override
func (state *SimpleState) End(ctx context.Context, event Event) Event {
	return nil
}

// GetEvent to override
func (state *SimpleState) GetEvent() Event {
	if state.machine != nil && state.machine.currentState != nil {
		return state.machine.currentState.GetEvent()
	}
	return nil
}

// GetTransitions to override
func (state *SimpleState) GetTransitions() Transitions {
	return nil
}

// InitialChildState to override
func (state *SimpleState) InitialChildState() State {
	return nil
}

func (state *SimpleState) CurrentState() State {
	return state.machine.CurrentState()
}

func (state *SimpleState) HasReaction(event Event) bool {
	if state.reactions == nil {
		return false
	}
	_, ok := state.reactions[TypeOf(event)]
	return ok
}

func (state *SimpleState) Outermost() *StateMachine {
	if state.parent == nil {
		return nil
	}
	return state.parent.Parent()
}

func (state *SimpleState) RegisterReaction(event Event, reaction Reaction) {
	if state.reactions == nil {
		state.reactions = make(map[reflect.Type]Reaction)
	}

	eventType := TypeOf(event)
	if _, ok := state.reactions[eventType]; ok {
		panic("event already exists")
	}
	state.reactions[eventType] = reaction
}

//nolint
func (state *SimpleState) initiate(ctx context.Context, parent *StateMachine, self State, event Event) Event {
	state.parent = parent
	child := self.InitialChildState()
	if child != nil {
		machine := NewStateMachine(child, ctx)
		machine.parent = parent
		if err := machine.Initiate(event); err != nil {
			return err
		}
		state.machine = machine
	}
	return nil
}

//nolint
func (state *SimpleState) React(ctx context.Context, event Event) (ret Event) {
	if state.reactions != nil {
		reaction, ok := state.reactions[TypeOf(event)]
		if ok {
			ret = reaction(ctx, event)
		}
	}

	machine := state.machine
	if machine != nil {
		machine.ProcessEvent(ctx, event)
	}
	return ret
}

//nolint
func (state *SimpleState) terminate(event Event) {
	machine := state.machine
	if machine != nil {
		machine.Terminate(event)
	}
}
