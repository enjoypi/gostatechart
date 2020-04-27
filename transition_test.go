package gostatechart

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransitions(t *testing.T) {
	trans := NewTranstions()
	require.NotNil(t, trans)

	trans.RegisterTransition((*EvClose)(nil), (*Greeting)(nil))
	require.Equal(t,
		reflect.TypeOf((*Greeting)(nil)),
		trans[reflect.TypeOf((*EvClose)(nil))],
	)

	trans.RegisterTransition((*EvClose)(nil), (*Greeting)(nil))
	require.Equal(t,
		reflect.TypeOf((*Greeting)(nil)),
		trans[reflect.TypeOf((*EvClose)(nil))],
	)

	defer func() {
		if e := recover(); e != nil {
			require.Equal(t, fmt.Errorf("duplicate event EvClose"), e)
		}
	}()
	trans.RegisterTransition((*EvClose)(nil), (*Active)(nil))
}
