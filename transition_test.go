package gostatechart_test

import (
	"fmt"
	"reflect"
	"testing"

	sc "github.com/enjoypi/gostatechart"
	"github.com/stretchr/testify/require"
)

type evTest struct {
}

type st1 struct {
	sc.SimpleState
}

func (s *st1) Begin(context interface{}, event sc.Event) sc.Event {
	panic("implement me")
}

func (s *st1) GetTransitions() sc.Transitions {
	panic("implement me")
}

type st2 struct {
	sc.SimpleState
}

func (s *st2) Begin(context interface{}, event sc.Event) sc.Event {
	panic("implement me")
}

func (s *st2) GetTransitions() sc.Transitions {
	panic("implement me")
}

func TestTransitions(t *testing.T) {
	trans := sc.NewTranstions()
	require.NotNil(t, trans)

	trans.RegisterTransition((*evTest)(nil), (*st1)(nil))
	require.Equal(t,
		reflect.TypeOf((*st1)(nil)),
		trans[reflect.TypeOf((*evTest)(nil))],
	)

	trans.RegisterTransition((*evTest)(nil), (*st1)(nil))
	require.Equal(t,
		reflect.TypeOf((*st1)(nil)),
		trans[reflect.TypeOf((*evTest)(nil))],
	)

	defer func() {
		if e := recover(); e != nil {
			require.Equal(t, fmt.Errorf("duplicate event evTest"), e)
		}
	}()
	trans.RegisterTransition((*evTest)(nil), (*st2)(nil))
}
